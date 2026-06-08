package handler

import (
	"strconv"
	"tiktok-simple/internal/service"
	"tiktok-simple/pkg/response"

	"github.com/gin-gonic/gin"
)

// CommentAction 发评论/删评论接口
// 路由：POST /douyin/comment/action/
// 请求参数：
//   - video_id:    视频ID
//   - action_type: 1=发评论 2=删评论
//   - comment_text: 评论内容（action_type=1时必填）
//   - comment_id:   评论ID（action_type=2时必填）
func CommentAction(c *gin.Context) {
	userID := c.GetUint("user_id")

	videoIDStr := c.PostForm("video_id")
	if videoIDStr == "" {
		videoIDStr = c.Query("video_id")
	}
	actionTypeStr := c.PostForm("action_type")
	if actionTypeStr == "" {
		actionTypeStr = c.Query("action_type")
	}

	videoIDInt, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil || videoIDInt == 0 {
		response.Fail(c, 1, "video_id 参数错误")
		return
	}

	actionType, err := strconv.Atoi(actionTypeStr)
	if err != nil || (actionType != 1 && actionType != 2) {
		response.Fail(c, 1, "action_type 参数错误，1=发评论 2=删评论")
		return
	}

	// 发评论时需要content，删评论时需要comment_id
	content := c.PostForm("comment_text")
	commentIDStr := c.PostForm("comment_id")
	var commentID uint
	if commentIDStr != "" {
		id, _ := strconv.ParseUint(commentIDStr, 10, 64)
		commentID = uint(id)
	}

	comment, err := service.CommentAction(userID, uint(videoIDInt), actionType, content, commentID)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	// 发评论成功时返回评论详情，删评论成功时comment为nil
	if comment != nil {
		response.Success(c, gin.H{"comment": comment})
	} else {
		response.Success(c, gin.H{})
	}
}

// CommentList 获取视频评论列表接口
// 路由：GET /douyin/comment/list/
// 不需要登录，任何人都能看评论
func CommentList(c *gin.Context) {
	videoIDStr := c.Query("video_id")
	videoIDInt, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil || videoIDInt == 0 {
		response.Fail(c, 1, "video_id 参数错误")
		return
	}

	comments, err := service.GetCommentList(uint(videoIDInt))
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	response.Success(c, gin.H{"comment_list": comments})
}
