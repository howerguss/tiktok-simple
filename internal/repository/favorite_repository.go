package repository

import (
	"tiktok-simple/internal/model"
	"tiktok-simple/pkg/database"

	"gorm.io/gorm"
)

// CreateFavorite 插入一条点赞记录（在事务中调用）
// tx 是事务对象，传入事务而不是直接用 database.DB，是为了保证原子性
func CreateFavorite(tx *gorm.DB, favorite *model.Favorite) error {
	return tx.Create(favorite).Error
}

// DeleteFavorite 删除一条点赞记录（在事务中调用）
func DeleteFavorite(tx *gorm.DB, userID, videoID uint) error {
	// Delete 加 Where 条件，精确删除指定记录
	// 必须加 Where 条件，否则 GORM 会拒绝执行（防止误删全表）
	return tx.Where("user_id = ? AND video_id = ?", userID, videoID).
		Delete(&model.Favorite{}).Error
}

// GetFavoriteByUserAndVideo 查询用户是否点赞了某个视频
// 返回 (记录, nil) 表示已点赞
// 返回 (nil, gorm.ErrRecordNotFound) 表示未点赞
func GetFavoriteByUserAndVideo(userID, videoID uint) (*model.Favorite, error) {
	var favorite model.Favorite
	err := database.DB.
		Where("user_id = ? AND video_id = ?", userID, videoID).
		First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

// GetFavoriteVideosByUserID 查询用户点赞过的所有视频
func GetFavoriteVideosByUserID(userID uint) ([]model.Favorite, error) {
	var favorites []model.Favorite
	err := database.DB.
		// Preload Video 把视频信息查出来
		// Preload Video.Author 把视频的作者信息也查出来（嵌套Preload）
		Preload("Video").
		Preload("Video.Author").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&favorites).Error
	return favorites, err
}

// IncrFavoriteCount 视频点赞数 +1（在事务中调用）
func IncrFavoriteCount(tx *gorm.DB, videoID uint) error {
	// gorm.Expr("favorite_count + 1") 生成 SQL 表达式
	// 最终执行的 SQL：UPDATE videos SET favorite_count = favorite_count + 1 WHERE id = ?
	// 为什么用表达式而不是先查出来再+1？
	// 先查后+1（read-modify-write）在并发时有竞态条件：
	//   线程A读到count=10，线程B也读到count=10
	//   线程A写入11，线程B也写入11，但实际应该是12
	// 直接用数据库表达式更新，由数据库保证原子性，不存在竞态问题
	return tx.Model(&model.Video{}).
		Where("id = ?", videoID).
		Update("favorite_count", gorm.Expr("favorite_count + 1")).Error
}

// DecrFavoriteCount 视频点赞数 -1（在事务中调用）
func DecrFavoriteCount(tx *gorm.DB, videoID uint) error {
	// GREATEST(favorite_count - 1, 0) 防止计数变成负数
	// 正常情况不会出现负数，但加个保险更健壮
	return tx.Model(&model.Video{}).
		Where("id = ?", videoID).
		Update("favorite_count", gorm.Expr("GREATEST(favorite_count - 1, 0)")).Error
}
