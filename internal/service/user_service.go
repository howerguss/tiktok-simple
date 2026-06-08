package service

import (
	"errors"
	"tiktok-simple/internal/model"
	"tiktok-simple/internal/repository"
	"tiktok-simple/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterResult 注册成功后返回的数据
type RegisterResult struct {
	UserID uint
	Token  string
}

// Register 处理注册业务逻辑
// 步骤：检查用户名 → 加密密码 → 写数据库 → 生成Token
func Register(username, password string) (*RegisterResult, error) {
	// 第一步：检查用户名是否已经被注册
	_, err := repository.GetUserByUsername(username)
	if err == nil {
		// err==nil 说明查到了这个用户名，意味着已经被注册了
		return nil, errors.New("用户名已存在")
	}
	// errors.Is 判断错误类型
	// 如果不是"记录不存在"这个错误，说明是数据库出了其他问题
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// 走到这里说明：err == gorm.ErrRecordNotFound，即用户名还没被注册，可以继续

	// 第二步：对密码进行 Bcrypt 加密
	// 为什么用Bcrypt而不是MD5？
	// - MD5是确定性哈希：同样的密码每次哈希结果相同，可以用彩虹表反查
	// - Bcrypt内置随机盐：同样的密码每次结果不同，且计算很慢，暴力破解成本极高
	// DefaultCost=10，表示哈希计算的"难度"，越高越安全但越慢
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 第三步：创建用户，写入数据库
	user := &model.User{
		Username: username,
		Password: string(hashedPassword), // 存加密后的密码，绝不存明文
	}
	if err := repository.CreateUser(user); err != nil {
		return nil, err
	}
	// CreateUser 执行后，user.ID 会被自动填充为数据库分配的自增ID

	// 第四步：生成JWT Token，用于后续接口鉴权
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &RegisterResult{UserID: user.ID, Token: token}, nil
}

// LoginResult 登录成功后返回的数据
type LoginResult struct {
	UserID uint
	Token  string
}

// Login 处理登录业务逻辑
// 步骤：查用户 → 验证密码 → 生成Token
func Login(username, password string) (*LoginResult, error) {
	// 第一步：查找用户
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		// ⚠️ 安全细节：不管是用户不存在还是密码错误，都返回同样的提示
		// 如果说"用户不存在"，攻击者可以枚举出哪些用户名是有效的
		return nil, errors.New("用户名或密码错误")
	}

	// 第二步：验证密码
	// CompareHashAndPassword 内部会处理盐的问题，直接比较即可
	// 如果密码不匹配，返回非nil的err
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 第三步：密码正确，生成Token
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginResult{UserID: user.ID, Token: token}, nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(userID uint) (*model.User, error) {
	return repository.GetUserByID(userID)
}
