package handler

import (
	"strconv"
	"tiktok-simple/internal/service"
	"tiktok-simple/pkg/response"

	"github.com/gin-gonic/gin"
)

// RelationAction 关注/取消关注接口
// 路由：POST /douyin/relation/action/
// 请求参数：
//   - to_user_id:  被关注者的用户ID
//   - action_type: 1=关注 2=取消关注
func RelationAction(c *gin.Context) {
	// 当前登录用户（关注者）
	followerID := c.GetUint("user_id")

	// 获取被关注者ID
	toUserIDStr := c.PostForm("to_user_id")
	if toUserIDStr == "" {
		toUserIDStr = c.Query("to_user_id")
	}
	actionTypeStr := c.PostForm("action_type")
	if actionTypeStr == "" {
		actionTypeStr = c.Query("action_type")
	}

	toUserIDInt, err := strconv.ParseUint(toUserIDStr, 10, 64)
	if err != nil || toUserIDInt == 0 {
		response.Fail(c, 1, "to_user_id 参数错误")
		return
	}

	actionType, err := strconv.Atoi(actionTypeStr)
	if err != nil || (actionType != 1 && actionType != 2) {
		response.Fail(c, 1, "action_type 参数错误，1=关注 2=取消关注")
		return
	}

	if err := service.FollowAction(followerID, uint(toUserIDInt), actionType); err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	response.Success(c, gin.H{})
}

// FollowList 获取我关注的人列表
// 路由：GET /douyin/relation/follow/list/
func FollowList(c *gin.Context) {
	userID := c.GetUint("user_id")

	users, err := service.GetFollowList(userID)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	response.Success(c, gin.H{"user_list": users})
}

// FollowerList 获取关注我的人列表（粉丝列表）
// 路由：GET /douyin/relation/follower/list/
func FollowerList(c *gin.Context) {
	userID := c.GetUint("user_id")

	users, err := service.GetFollowerList(userID)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	response.Success(c, gin.H{"user_list": users})
}

// FriendList 获取好友列表（互相关注）
// 路由：GET /douyin/relation/friend/list/
func FriendList(c *gin.Context) {
	userID := c.GetUint("user_id")

	users, err := service.GetFriendList(userID)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	response.Success(c, gin.H{"user_list": users})
}
