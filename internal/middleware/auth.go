package middleware

import (
	"strings"
	"tiktok-simple/pkg/jwt"
	"tiktok-simple/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 返回一个 Gin 中间件函数
// 中间件的作用：在真正的handler执行之前，先做一些公共处理
// 这里的作用是：验证请求是否携带了有效的JWT token
func AuthMiddleware() gin.HandlerFunc {
	// 返回一个函数，这个函数就是中间件
	return func(c *gin.Context) {
		var tokenStr string

		// 从两个地方尝试获取token（兼容不同客户端的传参方式）

		// 方式1：URL Query参数，比如 /douyin/user/?token=xxx
		tokenStr = c.Query("token")

		// 方式2：HTTP Header，格式为 Authorization: Bearer xxx
		if tokenStr == "" {
			authHeader := c.GetHeader("Authorization")
			// HasPrefix 判断是否以 "Bearer " 开头
			if strings.HasPrefix(authHeader, "Bearer ") {
				// TrimPrefix 去掉 "Bearer " 前缀，只保留token本身
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// 没有token，拒绝请求
		if tokenStr == "" {
			response.Fail(c, 401, "请先登录")
			// c.Abort() 非常重要！
			// 它会阻止后续的中间件和handler执行
			// 如果不调用Abort，即使我们返回了错误，请求仍会继续往下走
			c.Abort()
			return
		}

		// 解析和验证token
		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			response.Fail(c, 401, "token 无效或已过期")
			c.Abort()
			return
		}

		// 验证通过，把用户信息存入 Gin 的 Context
		// 这样后续的 handler 就可以用 c.GetUint("user_id") 直接拿到当前用户ID
		// 不用每个handler都重新解析一遍token
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// c.Next() 继续执行后续的中间件和handler
		c.Next()
	}
}
