package model

import "time"

// Comment 对应数据库里的 comments 表
type Comment struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// index: 给 video_id 建索引，查某个视频的评论时用到（WHERE video_id = ?）
	VideoID uint `gorm:"index;not null" json:"video_id"`
	UserID  uint `gorm:"not null" json:"user_id"`

	// 评论内容，最长512个字符
	Content string `gorm:"size:512;not null" json:"content"`

	CreatedAt time.Time `json:"created_at"`

	// 关联评论者信息，查评论列表时同时查出用户信息
	// 这样前端显示评论时能知道是谁发的
	User User `gorm:"foreignKey:UserID" json:"user"`
}

func (Comment) TableName() string {
	return "comments"
}
