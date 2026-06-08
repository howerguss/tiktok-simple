package repository

import (
	"tiktok-simple/internal/model"
	"tiktok-simple/pkg/database"
)

// GetUserByUsername 根据用户名查询用户
// 返回值：找到返回(*User, nil)，没找到返回(nil, gorm.ErrRecordNotFound)
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	// Where 设置查询条件，? 是占位符，防止SQL注入
	// First 查找第一条匹配的记录，找不到会返回 gorm.ErrRecordNotFound 错误
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID查询用户
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	// DB.First(&user, id) 等价于 SELECT * FROM users WHERE id = id LIMIT 1
	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建新用户，插入一条记录到数据库
func CreateUser(user *model.User) error {
	// DB.Create 会执行 INSERT 语句，并且会自动填充 user.ID（自增主键）
	return database.DB.Create(user).Error
}
