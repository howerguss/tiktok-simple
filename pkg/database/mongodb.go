package database

import (
	"context"
	"fmt"
	"tiktok-simple/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB 是全局 MongoDB 客户端
var MongoDB *mongo.Database

func InitMongoDB() {
	cfg := config.Global.MongoDB

	// 设置连接超时为10秒
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建客户端连接
	// options.Client().ApplyURI 解析 MongoDB 连接字符串
	// 连接字符串格式：mongodb://用户名:密码@地址:端口/数据库名
	// 本地无密码：mongodb://localhost:27017
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		panic(fmt.Errorf("MongoDB 连接失败: %w", err))
	}

	// Ping 测试连接是否真正可用
	// Connect 只是创建了客户端对象，不代表连接成功，Ping 才会真正尝试连接
	if err = client.Ping(ctx, nil); err != nil {
		panic(fmt.Errorf("MongoDB Ping 失败: %w", err))
	}

	// 选择数据库
	// MongoDB 里：Client → Database → Collection（类似MySQL的：连接 → 数据库 → 表）
	MongoDB = client.Database(cfg.Database)

	fmt.Println("MongoDB 连接成功")
}

// GetCollection 获取指定集合（MongoDB的"集合"相当于MySQL的"表"）
func GetCollection(name string) *mongo.Collection {
	return MongoDB.Collection(name)
}
