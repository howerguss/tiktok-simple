package repository

import (
	"tiktok-simple/internal/model"
	"tiktok-simple/pkg/database"

	"gorm.io/gorm"
)

// CreateRelation 插入一条关注记录（在事务中调用）
func CreateRelation(tx *gorm.DB, relation *model.Relation) error {
	return tx.Create(relation).Error
}

// DeleteRelation 删除一条关注记录（在事务中调用）
func DeleteRelation(tx *gorm.DB, followerID, followeeID uint) error {
	result := tx.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&model.Relation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetRelation 查询关注关系是否存在
// 返回 (relation, nil) 表示已关注
// 返回 (nil, gorm.ErrRecordNotFound) 表示未关注
func GetRelation(followerID, followeeID uint) (*model.Relation, error) {
	var relation model.Relation
	err := database.DB.
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		First(&relation).Error
	if err != nil {
		return nil, err
	}
	return &relation, nil
}

// GetFollowList 查询用户关注的人列表（我关注了谁）
func GetFollowList(followerID uint) ([]model.Relation, error) {
	var relations []model.Relation
	err := database.DB.
		// Preload("Followee") 把被关注者的用户信息查出来
		// 这样返回给前端时，每条关注记录里都有被关注者的完整信息
		Preload("Followee").
		Where("follower_id = ?", followerID).
		Order("created_at DESC").
		Find(&relations).Error
	return relations, err
}

// GetFollowerList 查询关注用户的人列表（谁关注了我）
func GetFollowerList(followeeID uint) ([]model.Relation, error) {
	var relations []model.Relation
	err := database.DB.
		// Preload("Follower") 把关注者的用户信息查出来
		Preload("Follower").
		Where("followee_id = ?", followeeID).
		Order("created_at DESC").
		Find(&relations).Error
	return relations, err
}

// GetFriendList 查询好友列表（互相关注的人）
// 好友 = 我关注的人 中 也关注了我的人
func GetFriendList(userID uint) ([]model.User, error) {
	var users []model.User

	// 子查询思路：
	// 先找出"我关注的人"的ID列表
	// 再在这些人里找"也关注了我的人"
	//
	// 等价SQL：
	// SELECT u.* FROM users u
	// INNER JOIN relations r1 ON r1.followee_id = u.id AND r1.follower_id = userID
	// INNER JOIN relations r2 ON r2.follower_id = u.id AND r2.followee_id = userID
	err := database.DB.
		// Joins 做内连接，过滤出同时满足两个条件的用户
		// r1: 我关注了这个用户（follower_id=userID, followee_id=u.id）
		// r2: 这个用户也关注了我（follower_id=u.id, followee_id=userID）
		Joins("INNER JOIN relations r1 ON r1.followee_id = users.id AND r1.follower_id = ?", userID).
		Joins("INNER JOIN relations r2 ON r2.follower_id = users.id AND r2.followee_id = ?", userID).
		Find(&users).Error

	return users, err
}

// IncrFollowCount 用户关注数+1（在事务中调用）
// 谁主动关注别人，谁的 follow_count +1
func IncrFollowCount(tx *gorm.DB, userID uint) error {
	return tx.Model(&model.User{}).
		Where("id = ?", userID).
		Update("follow_count", gorm.Expr("follow_count + 1")).Error
}

// DecrFollowCount 用户关注数-1（在事务中调用）
func DecrFollowCount(tx *gorm.DB, userID uint) error {
	return tx.Model(&model.User{}).
		Where("id = ?", userID).
		Update("follow_count", gorm.Expr("GREATEST(follow_count - 1, 0)")).Error
}

// IncrFollowerCount 用户粉丝数+1（在事务中调用）
// 被别人关注时，被关注者的 follower_count +1
func IncrFollowerCount(tx *gorm.DB, userID uint) error {
	return tx.Model(&model.User{}).
		Where("id = ?", userID).
		Update("follower_count", gorm.Expr("follower_count + 1")).Error
}

// DecrFollowerCount 用户粉丝数-1（在事务中调用）
func DecrFollowerCount(tx *gorm.DB, userID uint) error {
	return tx.Model(&model.User{}).
		Where("id = ?", userID).
		Update("follower_count", gorm.Expr("GREATEST(follower_count - 1, 0)")).Error
}
