package service

import (
	"errors"
	"fmt"
	"tiktok-simple/internal/model"
	"tiktok-simple/internal/repository"
	"tiktok-simple/pkg/database"

	"gorm.io/gorm"
)

// FavoriteAction 处理点赞/取消点赞
//
// actionType:
//
//	1 = 点赞
//	2 = 取消点赞
//
// 为什么用事务？
//
//	点赞需要同时做两件事：
//	1. 往 favorites 表插入/删除记录
//	2. 更新 videos 表的 favorite_count 字段
//	如果只做了第1步，第2步失败了，点赞数就不准确了
//	事务保证两步要么全成功，要么全回滚，数据始终一致
func FavoriteAction(userID, videoID uint, actionType int) error {
	// 先检查视频是否存在
	_, err := repository.GetVideoByID(videoID)
	if err != nil {
		return errors.New("视频不存在")
	}

	// 用 Redis Set 缓存点赞状态
	// key 格式：like:video:{videoID}
	// value：点赞了这个视频的所有用户ID的集合
	cacheKey := fmt.Sprintf("like:video:%d", videoID)

	if actionType == 1 {
		// ========== 点赞 ==========

		// 先检查是否已经点赞（防止重复点赞）
		// 先查Redis，Redis没有再查数据库
		isMember, err := database.RDB.SIsMember(cacheKey, userID).Result()
		if err == nil && isMember {
			return errors.New("已经点赞过了")
		}
		// Redis没有，查数据库确认
		if err != nil {
			_, dbErr := repository.GetFavoriteByUserAndVideo(userID, videoID)
			if dbErr == nil {
				return errors.New("已经点赞过了")
			}
		}

		// 开启事务
		// database.DB.Transaction 会自动处理提交和回滚：
		// - 回调函数返回 nil → 自动提交（COMMIT）
		// - 回调函数返回 error → 自动回滚（ROLLBACK）
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			// 事务第1步：插入点赞记录
			favorite := &model.Favorite{
				UserID:  userID,
				VideoID: videoID,
			}
			if err := repository.CreateFavorite(tx, favorite); err != nil {
				return err // 返回error触发回滚
			}

			// 事务第2步：视频点赞数+1
			if err := repository.IncrFavoriteCount(tx, videoID); err != nil {
				return err // 返回error触发回滚
			}

			return nil // 返回nil触发提交
		})

		if err != nil {
			return err
		}

		// 事务成功后，更新Redis缓存
		// 为什么在事务成功后才更新Redis，而不是在事务内？
		// 因为Redis操作不参与数据库事务，如果事务回滚了但Redis已更新，就会不一致
		// 先保证数据库成功，再更新缓存，即使Redis更新失败，
		// 下次查询时Redis没有，会去数据库查，数据不会错
		database.RDB.SAdd(cacheKey, userID)

	} else if actionType == 2 {
		// ========== 取消点赞 ==========

		// 开启事务
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			// 事务第1步：删除点赞记录
			if err := repository.DeleteFavorite(tx, userID, videoID); err != nil {
				return err
			}

			// 事务第2步：视频点赞数-1
			if err := repository.DecrFavoriteCount(tx, videoID); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}

		// 事务成功后，从Redis缓存中删除
		database.RDB.SRem(cacheKey, userID)

	} else {
		return errors.New("action_type 参数错误，1=点赞 2=取消点赞")
	}

	return nil
}

// GetFavoriteList 获取用户点赞过的视频列表
func GetFavoriteList(userID uint) ([]model.Video, error) {
	favorites, err := repository.GetFavoriteVideosByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 把 []Favorite 转成 []Video 返回给前端
	// 前端只需要视频信息，不需要知道"点赞记录"这个概念
	videos := make([]model.Video, 0, len(favorites))
	for _, f := range favorites {
		videos = append(videos, f.Video)
	}
	return videos, nil
}
