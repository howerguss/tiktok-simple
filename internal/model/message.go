package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message 是聊天消息的数据模型，存储在 MongoDB 而不是 MySQL
//
// 和 MySQL 模型的区别：
// - MySQL 用 gorm tag，MongoDB 用 bson tag
// - MySQL 用 uint 作为主键，MongoDB 用 primitive.ObjectID（12字节的唯一ID）
// - MongoDB 不需要 AutoMigrate，Collection 会自动创建
//
// bson tag 说明：
// - bson:"_id" 对应 MongoDB 的主键字段（固定叫 _id）
// - omitempty 表示如果字段是零值就不序列化（避免插入时 _id 为空）
// - bson:"from_user_id" 指定存入 MongoDB 时的字段名
type Message struct {
	// ObjectID 是 MongoDB 自动生成的唯一主键
	// 格式：12字节 = 4字节时间戳 + 5字节机器ID + 3字节计数器
	// 天然有序（按时间），可以用来排序
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 发送者ID（对应 MySQL users 表的 id）
	// MongoDB 不做外键约束，我们自己在业务层保证数据一致性
	FromUserID uint `bson:"from_user_id" json:"from_user_id"`

	// 接收者ID
	ToUserID uint `bson:"to_user_id" json:"to_user_id"`

	// 消息内容
	Content string `bson:"content" json:"content"`

	// 发送时间（Unix时间戳，毫秒级）
	// 为什么用毫秒而不是秒？
	// 同一秒内可能发多条消息，毫秒级时间戳能保证更精确的排序
	CreateTime int64 `bson:"create_time" json:"create_time"`
}
