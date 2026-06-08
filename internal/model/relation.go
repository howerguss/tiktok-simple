package model

import "time"

// Relation 对应数据库里的 relations 表
// 记录"谁关注了谁"，是一张有向关系表
//
// 例子：
//
//	follower_id=1, followee_id=2 → 用户1 关注了 用户2
//	follower_id=2, followee_id=1 → 用户2 关注了 用户1
//	两条记录都存在 → 互相关注（好友）
type Relation struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// follower: 关注者（主动关注别人的人，即"粉丝"视角的"我"）
	// followee: 被关注者（被别人关注的人）
	//
	// 联合唯一索引：同一个用户对同一个人只能关注一次
	// uniqueIndex:uidx_follow 是索引名，两个字段用同一个名字 = 联合索引
	FollowerID uint `gorm:"uniqueIndex:uidx_follow;not null" json:"follower_id"`
	FolloweeID uint `gorm:"uniqueIndex:uidx_follow;not null" json:"followee_id"`

	CreatedAt time.Time `json:"created_at"`

	// 关联被关注者的用户信息
	// 查"我关注的人列表"时，需要返回每个人的用户信息
	Followee User `gorm:"foreignKey:FolloweeID" json:"user"`

	// 关联关注者的用户信息
	// 查"关注我的人列表"时，需要返回每个粉丝的用户信息
	Follower User `gorm:"foreignKey:FollowerID" json:"follower"`
}

func (Relation) TableName() string {
	return "relations"
}
