package repository

import (
	"tiktok-simple/internal/model"
	"tiktok-simple/pkg/database"

	"gorm.io/gorm"
)

// CreateComment 插入一条评论记录（在事务中调用）
func CreateComment(tx *gorm.DB, comment *model.Comment) error {
	return tx.Create(comment).Error
}

// DeleteComment 删除一条评论（在事务中调用）
// 只有评论的作者才能删除自己的评论，所以要同时校验 user_id
func DeleteComment(tx *gorm.DB, commentID, userID uint) error {
	result := tx.Where("id = ? AND user_id = ?", commentID, userID).
		Delete(&model.Comment{})

	if result.Error != nil {
		return result.Error
	}

	// RowsAffected 是实际删除的行数
	// 如果是0，说明这条评论不存在或者不属于当前用户
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetCommentsByVideoID 查询视频的所有评论
func GetCommentsByVideoID(videoID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := database.DB.
		// 同时查出评论者的用户信息
		Preload("User").
		Where("video_id = ?", videoID).
		// 评论按时间正序：最早的评论在最上面（和大多数App一样）
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

// IncrCommentCount 视频评论数 +1（在事务中调用）
func IncrCommentCount(tx *gorm.DB, videoID uint) error {
	return tx.Model(&model.Video{}).
		Where("id = ?", videoID).
		Update("comment_count", gorm.Expr("comment_count + 1")).Error
}

// DecrCommentCount 视频评论数 -1（在事务中调用）
func DecrCommentCount(tx *gorm.DB, videoID uint) error {
	return tx.Model(&model.Video{}).
		Where("id = ?", videoID).
		Update("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)")).Error
}
