package service

import (
	"errors"
	"fmt"
	"tiktok-simple/internal/model"
	"tiktok-simple/internal/repository"
	"tiktok-simple/pkg/database"

	"gorm.io/gorm"
)

// FollowAction 处理关注/取消关注
//
// actionType:
//
//	1 = 关注
//	2 = 取消关注
//
// 业务约束：
//   - 不能关注自己
//   - 不能重复关注
//   - 关注和取关都要用事务同步更新用户计数
func FollowAction(followerID, followeeID uint, actionType int) error {
	// 不能关注自己
	if followerID == followeeID {
		return errors.New("不能关注自己")
	}

	// 检查被关注者是否存在
	_, err := repository.GetUserByID(followeeID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// Redis缓存key：存储"用户A关注的所有人"的集合
	// key格式：follow:{followerID}  value：该用户关注的所有人的ID集合
	cacheKey := fmt.Sprintf("follow:%d", followerID)

	if actionType == 1 {
		// ========== 关注 ==========

		// Redis预检：是否已经关注（O(1)查询，比查数据库快）
		isMember, err := database.RDB.SIsMember(cacheKey, followeeID).Result()
		if err == nil && isMember {
			return errors.New("已经关注过了")
		}
		// Redis没有命中，去数据库确认
		if err != nil {
			_, dbErr := repository.GetRelation(followerID, followeeID)
			if dbErr == nil {
				return errors.New("已经关注过了")
			}
		}

		// 开启事务，同时完成三件事：
		// 1. 插入关注记录
		// 2. 关注者的 follow_count + 1
		// 3. 被关注者的 follower_count + 1
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			// 第1步：插入关注记录
			relation := &model.Relation{
				FollowerID: followerID,
				FolloweeID: followeeID,
			}
			if err := repository.CreateRelation(tx, relation); err != nil {
				return err
			}

			// 第2步：关注者关注数+1（我的关注数增加）
			if err := repository.IncrFollowCount(tx, followerID); err != nil {
				return err
			}

			// 第3步：被关注者粉丝数+1（对方的粉丝数增加）
			if err := repository.IncrFollowerCount(tx, followeeID); err != nil {
				return err
			}

			return nil // 三步都成功，提交事务
		})

		if err != nil {
			return err
		}

		// 事务成功后更新Redis缓存
		database.RDB.SAdd(cacheKey, followeeID)

	} else if actionType == 2 {
		// ========== 取消关注 ==========

		err = database.DB.Transaction(func(tx *gorm.DB) error {
			// 第1步：删除关注记录
			if err := repository.DeleteRelation(tx, followerID, followeeID); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("还没有关注过该用户")
				}
				return err
			}

			// 第2步：关注者关注数-1
			if err := repository.DecrFollowCount(tx, followerID); err != nil {
				return err
			}

			// 第3步：被关注者粉丝数-1
			if err := repository.DecrFollowerCount(tx, followeeID); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}

		// 事务成功后从Redis缓存删除
		database.RDB.SRem(cacheKey, followeeID)

	} else {
		return errors.New("action_type 参数错误，1=关注 2=取消关注")
	}

	return nil
}

// GetFollowList 获取我关注的人列表
func GetFollowList(userID uint) ([]model.User, error) {
	relations, err := repository.GetFollowList(userID)
	if err != nil {
		return nil, err
	}

	// 从关注记录里提取被关注者的用户信息
	users := make([]model.User, 0, len(relations))
	for _, r := range relations {
		users = append(users, r.Followee)
	}
	return users, nil
}

// GetFollowerList 获取关注我的人列表（粉丝）
func GetFollowerList(userID uint) ([]model.User, error) {
	relations, err := repository.GetFollowerList(userID)
	if err != nil {
		return nil, err
	}

	// 从关注记录里提取关注者的用户信息
	users := make([]model.User, 0, len(relations))
	for _, r := range relations {
		users = append(users, r.Follower)
	}
	return users, nil
}

// GetFriendList 获取好友列表（互相关注的人）
func GetFriendList(userID uint) ([]model.User, error) {
	return repository.GetFriendList(userID)
}
