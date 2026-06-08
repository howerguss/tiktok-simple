package model

import "time"

// Video 对应数据库里的 videos 表
type Video struct {
	// 主键，自增
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// index: 普通索引，查询某个用户发布的视频时会用到（WHERE user_id = ?）
	// 加索引的原因：如果没有索引，查询时会全表扫描，数据量大时很慢
	UserID uint `gorm:"index;not null" json:"author_id"`

	// 视频标题
	Title string `gorm:"size:128;not null" json:"title"`

	// 视频文件的访问URL（存相对路径，返回时拼完整URL）
	PlayURL string `gorm:"size:256;not null" json:"play_url"`

	// 封面图片的访问URL
	CoverURL string `gorm:"size:256" json:"cover_url"`

	// 点赞数和评论数，查视频列表时直接返回，不用再去count
	FavoriteCount int64 `gorm:"default:0" json:"favorite_count"`
	CommentCount  int64 `gorm:"default:0" json:"comment_count"`

	CreatedAt time.Time `json:"created_at"`

	// 关联作者信息，gorm:"foreignKey:UserID" 告诉GORM用UserID做外键
	// 这样查视频时可以同时把作者信息查出来（Preload）
	// json:"author" 返回给前端时包含作者的完整信息
	Author User `gorm:"foreignKey:UserID" json:"author"`
}

func (Video) TableName() string {
	return "videos"
}
