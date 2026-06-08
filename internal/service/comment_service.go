package service

import (
	"errors"
	"tiktok-simple/internal/model"
	"tiktok-simple/internal/repository"
	"tiktok-simple/pkg/database"

	"gorm.io/gorm"
)

// CommentAction 处理发评论/删评论
//
// actionType:
//
//	1 = 发评论（需要content参数）
//	2 = 删评论（需要commentID参数）
func CommentAction(userID, videoID uint, actionType int, content string, commentID uint) (*model.Comment, error) {
	// 检查视频是否存在
	_, err := repository.GetVideoByID(videoID)
	if err != nil {
		return nil, errors.New("视频不存在")
	}

	if actionType == 1 {
		// ========== 发评论 ==========

		if content == "" {
			return nil, errors.New("评论内容不能为空")
		}
		if len([]rune(content)) > 200 {
			// len([]rune(content)) 计算的是字符数（支持中文）
			// len(content) 计算的是字节数（一个中文字符占3字节，会不准确）
			return nil, errors.New("评论内容不能超过200个字符")
		}

		var newComment *model.Comment

		// 开启事务：同时完成插入评论和更新评论数
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			// 第1步：插入评论
			comment := &model.Comment{
				UserID:  userID,
				VideoID: videoID,
				Content: content,
			}
			if err := repository.CreateComment(tx, comment); err != nil {
				return err
			}
			newComment = comment // 保存创建好的评论，后面要返回给前端

			// 第2步：视频评论数+1
			return repository.IncrCommentCount(tx, videoID)
		})

		if err != nil {
			return nil, err
		}

		// 查询完整的评论信息（包含用户信息），返回给前端
		// 因为插入后 newComment 里只有基础字段，User 字段是空的
		// 需要重新查一次拿到完整数据
		comments, _ := repository.GetCommentsByVideoID(videoID)
		for _, c := range comments {
			if c.ID == newComment.ID {
				return &c, nil
			}
		}
		return newComment, nil

	} else if actionType == 2 {
		// ========== 删评论 ==========

		if commentID == 0 {
			return nil, errors.New("comment_id 不能为空")
		}

		// 开启事务：同时完成删除评论和更新评论数
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			// 第1步：删除评论（DeleteComment内部会校验是不是自己的评论）
			if err := repository.DeleteComment(tx, commentID, userID); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("评论不存在或无权删除")
				}
				return err
			}

			// 第2步：视频评论数-1
			return repository.DecrCommentCount(tx, videoID)
		})

		return nil, err

	} else {
		return nil, errors.New("action_type 参数错误，1=发评论 2=删评论")
	}
}

// GetCommentList 获取视频的评论列表
func GetCommentList(videoID uint) ([]model.Comment, error) {
	return repository.GetCommentsByVideoID(videoID)
}
