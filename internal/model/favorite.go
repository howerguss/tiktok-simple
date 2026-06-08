package model

import "time"

// Favorite 对应数据库里的 favorites 表
// 记录"哪个用户点赞了哪个视频"，是一张纯粹的关系表
type Favorite struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// 联合唯一索引：同一个用户对同一个视频只能点赞一次
	// uniqueIndex:uidx_user_video 是索引名称，两个字段用同一个名字表示是联合索引
	UserID  uint `gorm:"uniqueIndex:uidx_user_video;not null" json:"user_id"`
	VideoID uint `gorm:"uniqueIndex:uidx_user_video;not null" json:"video_id"`

	CreatedAt time.Time `json:"created_at"`

	// 关联视频信息，查点赞列表时可以同时查出视频详情
	Video Video `gorm:"foreignKey:VideoID" json:"video"`
}

func (Favorite) TableName() string {
	return "favorites"
}
