package service

import (
	"errors"
	"tiktok-simple/internal/model"
	"tiktok-simple/internal/repository"
)

// SendMessage 发送消息
//
// 参数：
//
//	fromUserID: 发送者ID（从JWT中间件获取）
//	toUserID:   接收者ID
//	content:    消息内容
func SendMessage(fromUserID, toUserID uint, content string) error {
	// 不能给自己发消息
	if fromUserID == toUserID {
		return errors.New("不能给自己发消息")
	}

	// 检查接收者是否存在（消息发给不存在的用户没有意义）
	_, err := repository.GetUserByID(toUserID)
	if err != nil {
		return errors.New("对方用户不存在")
	}

	if content == "" {
		return errors.New("消息内容不能为空")
	}

	// 消息长度限制：500个字符
	if len([]rune(content)) > 500 {
		return errors.New("消息内容不能超过500个字符")
	}

	// 写入 MongoDB
	return repository.SendMessage(fromUserID, toUserID, content)
}

// GetMessageList 获取两人之间的聊天记录
//
// 参数：
//
//	userID:     当前登录用户ID
//	toUserID:   对方用户ID
//	preMsgTime: 上次拉取的最后一条消息时间戳（毫秒），传0获取全部
func GetMessageList(userID, toUserID uint, preMsgTime int64) ([]model.Message, error) {
	if userID == toUserID {
		return nil, errors.New("不能查看和自己的聊天记录")
	}

	return repository.GetMessageList(userID, toUserID, preMsgTime)
}
