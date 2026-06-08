package handler

import (
	"strconv"
	"tiktok-simple/internal/service"
	"tiktok-simple/pkg/response"

	"github.com/gin-gonic/gin"
)

// FavoriteAction 点赞/取消点赞接口
// 路由：POST /douyin/favorite/action/
// 请求参数：
//   - video_id:    视频ID
//   - action_type: 1=点赞 2=取消点赞
func FavoriteAction(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Query 取URL参数，PostForm 取表单参数
	// 这里两种方式都支持，先取表单，取不到再取URL参数
	videoIDStr := c.PostForm("video_id")
	if videoIDStr == "" {
		videoIDStr = c.Query("video_id")
	}
	actionTypeStr := c.PostForm("action_type")
	if actionTypeStr == "" {
		actionTypeStr = c.Query("action_type")
	}

	// 字符串转uint
	videoIDInt, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil || videoIDInt == 0 {
		response.Fail(c, 1, "video_id 参数错误")
		return
	}

	actionType, err := strconv.Atoi(actionTypeStr)
	if err != nil || (actionType != 1 && actionType != 2) {
		response.Fail(c, 1, "action_type 参数错误，1=点赞 2=取消点赞")
		return
	}

	if err := service.FavoriteAction(userID, uint(videoIDInt), actionType); err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	response.Success(c, gin.H{})
}

// FavoriteList 获取点赞列表接口
// 路由：GET /douyin/favorite/list/
func FavoriteList(c *gin.Context) {
	userID := c.GetUint("user_id")

	videos, err := service.GetFavoriteList(userID)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	response.Success(c, gin.H{"video_list": videos})
}
