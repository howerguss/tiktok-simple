package repository

import (
	"context"
	"tiktok-simple/internal/model"
	"tiktok-simple/pkg/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// collectionName 消息集合名（相当于MySQL的表名）
const collectionName = "messages"

// SendMessage 插入一条消息到 MongoDB
func SendMessage(fromUserID, toUserID uint, content string) error {
	// 获取 messages 集合
	col := database.GetCollection(collectionName)

	// 构建消息文档
	msg := model.Message{
		// primitive.NewObjectID() 生成一个新的唯一ID
		// 如果不设置，MongoDB 会自动生成，但显式设置更清晰
		ID:         primitive.NewObjectID(),
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Content:    content,
		// time.Now().UnixMilli() 获取当前时间的毫秒级时间戳
		CreateTime: time.Now().UnixMilli(),
	}

	// 设置操作超时，避免 MongoDB 无响应时一直阻塞
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// InsertOne 插入一条文档
	// 等价于 MySQL 的 INSERT INTO messages (...) VALUES (...)
	_, err := col.InsertOne(ctx, msg)
	return err
}

// GetMessageList 获取两个用户之间的聊天记录
//
// 参数：
//
//	userA, userB: 两个用户的ID（顺序不重要）
//	preMsgTime:   只返回这个时间戳之后的消息（用于增量拉取）
//	              传0表示获取所有消息
//
// 为什么用 preMsgTime 做增量拉取？
// 聊天记录可能有几千条，每次全量拉取浪费流量和性能
// 客户端记录上次拉取的最后一条消息时间戳，下次只拉这之后的新消息
func GetMessageList(userA, userB uint, preMsgTime int64) ([]model.Message, error) {
	col := database.GetCollection(collectionName)

	// 构建查询条件：两人之间的消息 = (A发给B) 或 (B发给A)
	//
	// bson.M 是 MongoDB 的查询文档，类似 map[string]interface{}
	// bson.A 是数组类型
	//
	// 等价的 MongoDB 查询：
	// db.messages.find({
	//   $and: [
	//     { create_time: { $gt: preMsgTime } },
	//     { $or: [
	//       { from_user_id: userA, to_user_id: userB },
	//       { from_user_id: userB, to_user_id: userA }
	//     ]}
	//   ]
	// })
	filter := bson.M{
		"$and": bson.A{
			// 只返回 preMsgTime 之后的消息
			bson.M{"create_time": bson.M{"$gt": preMsgTime}},
			// 两人之间的消息（任意方向）
			bson.M{
				"$or": bson.A{
					bson.M{"from_user_id": userA, "to_user_id": userB},
					bson.M{"from_user_id": userB, "to_user_id": userA},
				},
			},
		},
	}

	// 查询选项：按 create_time 升序排列（最早的消息在前面，符合聊天界面习惯）
	// options.Find() 返回查询选项对象
	// SetSort(bson.D{{"create_time", 1}}) 设置排序：1=升序，-1=降序
	opts := options.Find().SetSort(bson.D{{Key: "create_time", Value: 1}})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find 返回游标（Cursor），需要遍历才能拿到数据
	// 类似 MySQL 的 SELECT ... 返回结果集
	cursor, err := col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	// 函数结束时关闭游标，释放资源
	defer cursor.Close(ctx)

	// cursor.All 把游标里的所有文档解码到 []model.Message
	var messages []model.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	// MongoDB 查不到数据时返回空切片而不是 nil，前端处理更方便
	if messages == nil {
		messages = []model.Message{}
	}

	return messages, nil
}
