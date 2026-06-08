package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 是整个项目的配置结构体，对应 config.yaml 的结构
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Storage  StorageConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
	Expire int
}

type StorageConfig struct {
	Path string
}

// Global 是全局配置对象，项目任何地方都可以用 config.Global.xxx 来读取配置
var Global *Config

func Init() {
	// 告诉viper配置文件的名字、格式和位置
	viper.SetConfigName("config")   // 文件名（不含扩展名）
	viper.SetConfigType("yaml")     // 文件格式
	viper.AddConfigPath("./config") // 文件所在目录

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	// 把读取到的配置反序列化到 Global 结构体里
	Global = &Config{}
	if err := viper.Unmarshal(Global); err != nil {
		panic(fmt.Errorf("解析配置失败: %w", err))
	}
}
