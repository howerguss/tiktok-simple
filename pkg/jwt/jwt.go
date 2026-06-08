package jwt

import (
	"errors"
	"tiktok-simple/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 是JWT的载荷部分（Payload），存放我们想带在token里的信息
// JWT = Header.Payload.Signature 三段，用.分隔，Base64编码
// Payload 是公开的（任何人都能解码看到），所以不要放密码等敏感信息
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	// jwt.RegisteredClaims 是标准字段，包含过期时间(ExpiresAt)、签发时间(IssuedAt)等
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
// 用户登录成功后调用这个函数，把token返回给客户端
// 客户端后续请求都要带上这个token，服务端通过ParseToken验证
func GenerateToken(userID uint, username string) (string, error) {
	secret := config.Global.JWT.Secret
	expire := config.Global.JWT.Expire

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间 = 当前时间 + 配置的小时数
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
		},
	}

	// 用 HS256 算法创建token（HMAC-SHA256，对称加密，用同一个secret签名和验证）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 用密钥签名，生成最终的token字符串
	return token.SignedString([]byte(secret))
}

// ParseToken 解析并验证JWT Token
// 验证内容：1.签名是否正确（防篡改）2.是否过期
func ParseToken(tokenStr string) (*Claims, error) {
	secret := config.Global.JWT.Secret

	// ParseWithClaims 解析token，第三个参数是验证密钥的回调函数
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 返回密钥，用于验证签名
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err // token过期或签名错误都会走这里
	}

	// 类型断言：把 token.Claims 转成我们自定义的 *Claims 类型
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token 无效")
	}

	return claims, nil
}
