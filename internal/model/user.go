package model

import "time"

// User 对应数据库里的 users 表
// GORM 通过 struct tag 来知道每个字段怎么映射到数据库列
type User struct {
	// primaryKey: 主键   autoIncrement: 自增（1,2,3...）
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// uniqueIndex: 唯一索引（不能有两个相同用户名）
	// size:32: 最大长度32字符
	// not null: 不能为空
	Username string `gorm:"uniqueIndex;size:32;not null" json:"username"`

	// json:"-" 非常重要！表示这个字段不会被序列化到JSON里
	// 即使你查出来了User对象，返回给前端时密码字段会自动消失
	Password string `gorm:"size:256;not null" json:"-"`

	Avatar    string `gorm:"size:256" json:"avatar"`
	Signature string `gorm:"size:256" json:"signature"`

	// default:0 数据库默认值为0
	FollowCount   int64 `gorm:"default:0" json:"follow_count"`
	FollowerCount int64 `gorm:"default:0" json:"follower_count"`

	CreatedAt time.Time `json:"created_at"`
}

// TableName 告诉GORM这个struct对应哪张表
// 如果不写，GORM默认会把 User -> users（自动加s变复数）
func (User) TableName() string {
	return "users"
}
