package database

import (
	"fmt"
	"tiktok-simple/config"
	"tiktok-simple/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 是全局数据库对象，项目任何地方都可以用 database.DB 来操作数据库
var DB *gorm.DB

func InitMySQL() {
	cfg := config.Global.Database

	// DSN (Data Source Name) 是数据库连接字符串，格式固定
	// charset=utf8mb4    支持emoji等特殊字符
	// parseTime=True     自动把数据库的datetime解析成Go的time.Time类型
	// loc=Local          使用本地时区
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 开发时用Info级别，会打印所有SQL语句，方便调试
		// 上线后改为 logger.Warn，只打印慢查询和错误
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Errorf("MySQL 连接失败: %w", err))
	}

	// 获取底层的 *sql.DB 对象来配置连接池
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(100) // 最大连接数，超过这个数的请求会排队等待
	sqlDB.SetMaxIdleConns(10)  // 最大空闲连接数，空闲连接超过这个数会被关闭

	// AutoMigrate 会自动创建或更新表结构，不会删除已有数据
	// 开发阶段很方便，加了新字段自动同步；生产环境建议用专门的migration工具
	// AutoMigrate同时建 users, videos, favorite, comment表
	if err := DB.AutoMigrate(
		&model.User{},
		&model.Video{},
		&model.Favorite{},
		&model.Comment{},
		&model.Relation{}, // 新增这行
	); err != nil {
		panic(fmt.Errorf("自动建表失败: %w", err))
	}

	fmt.Println("MySQL 连接成功")
}
