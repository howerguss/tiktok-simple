package handler

import (
	"strconv"
	"tiktok-simple/internal/service"
	"tiktok-simple/pkg/response"

	"github.com/gin-gonic/gin"
)

// MessageAction 发送消息接口
// 路由：POST /douyin/message/action/
// 请求参数：
//   - to_user_id:  接收者用户ID
//   - action_type: 操作类型，目前只有 1=发送消息
//   - content:     消息内容
func MessageAction(c *gin.Context) {
	// 当前登录用户（发送者）
	fromUserID := c.GetUint("user_id")

	// 获取接收者ID
	toUserIDStr := c.PostForm("to_user_id")
	if toUserIDStr == "" {
		toUserIDStr = c.Query("to_user_id")
	}

	toUserIDInt, err := strconv.ParseUint(toUserIDStr, 10, 64)
	if err != nil || toUserIDInt == 0 {
		response.Fail(c, 1, "to_user_id 参数错误")
		return
	}

	// action_type 目前只支持 1=发送消息
	// 预留这个参数是为了以后扩展（比如 2=撤回消息）
	actionTypeStr := c.PostForm("action_type")
	if actionTypeStr == "" {
		actionTypeStr = c.Query("action_type")
	}
	actionType, _ := strconv.Atoi(actionTypeStr)
	if actionType != 1 {
		response.Fail(c, 1, "action_type 参数错误，目前只支持 1=发送消息")
		return
	}

	content := c.PostForm("content")
	if content == "" {
		content = c.Query("content")
	}

	if err := service.SendMessage(fromUserID, uint(toUserIDInt), content); err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	response.Success(c, gin.H{})
}

// MessageChat 获取聊天记录接口
// 路由：GET /douyin/message/chat/
// 请求参数：
//   - to_user_id:   对方用户ID
//   - pre_msg_time: 上次拉取的最后一条消息时间戳（毫秒），不传或传0则获取全部
func MessageChat(c *gin.Context) {
	userID := c.GetUint("user_id")

	toUserIDStr := c.Query("to_user_id")
	toUserIDInt, err := strconv.ParseUint(toUserIDStr, 10, 64)
	if err != nil || toUserIDInt == 0 {
		response.Fail(c, 1, "to_user_id 参数错误")
		return
	}

	// pre_msg_time 是毫秒级时间戳
	// 不传或传0时，GetMessageList 会返回两人之间的全部消息
	var preMsgTime int64
	preMsgTimeStr := c.Query("pre_msg_time")
	if preMsgTimeStr != "" {
		preMsgTime, _ = strconv.ParseInt(preMsgTimeStr, 10, 64)
	}

	messages, err := service.GetMessageList(userID, uint(toUserIDInt), preMsgTime)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message_list": messages,
	})
}
