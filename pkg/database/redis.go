package database

import (
	"fmt"
	"tiktok-simple/config"

	"github.com/go-redis/redis"
)

// RDB 是全局Redis客户端，项目任何地方都可以用 database.RDB 来操作Redis
var RDB *redis.Client

func InitRedis() {
	cfg := config.Global.Redis

	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB, // Redis有16个数据库(0-15)，选择用哪个
	})

	// Ping测试连接是否正常
	// v6写法：直接 .Ping()，不需要传context（v9才需要传context）
	if err := RDB.Ping().Err(); err != nil {
		panic(fmt.Errorf("Redis 连接失败: %w", err))
	}

	fmt.Println("Redis 连接成功")
}
