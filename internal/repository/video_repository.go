package repository

import (
	"tiktok-simple/internal/model"
	"tiktok-simple/pkg/database"
)

// CreateVideo 创建视频记录
func CreateVideo(video *model.Video) error {
	return database.DB.Create(video).Error
}

// GetFeedVideos 获取Feed流视频列表
// latestTime: 返回这个时间之前发布的视频（用于翻页）
// limit: 每次返回多少条
func GetFeedVideos(latestTime int64, limit int) ([]model.Video, error) {
	var videos []model.Video

	// Preload("Author") 预加载作者信息
	// 执行的SQL大概是：
	// SELECT * FROM videos WHERE created_at < ? ORDER BY created_at DESC LIMIT ?
	// SELECT * FROM users WHERE id IN (查出来的所有user_id)
	// 然后GORM自动把User数据填充到每个Video的Author字段里
	err := database.DB.
		Preload("Author").
		Where("UNIX_TIMESTAMP(created_at) <= ?", latestTime).
		Order("created_at DESC"). // 按发布时间倒序，最新的在前面
		Limit(limit).
		Find(&videos).Error

	return videos, err
}

// GetVideosByUserID 获取某个用户发布的所有视频
func GetVideosByUserID(userID uint) ([]model.Video, error) {
	var videos []model.Video
	err := database.DB.
		Preload("Author").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&videos).Error
	return videos, err
}

// GetVideoByID 根据ID查询视频
func GetVideoByID(id uint) (*model.Video, error) {
	var video model.Video
	err := database.DB.First(&video, id).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}
