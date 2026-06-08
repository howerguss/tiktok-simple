# tiktok-simple — 极简抖音后端

基于 Go + Gin + GORM 实现的短视频社交平台后端，覆盖抖音核心功能。

## 技术栈

| 技术 | 用途 |
|------|------|
| Go 1.21+ | 主语言 |
| Gin | HTTP 框架 |
| GORM | ORM，操作 MySQL |
| MySQL | 主数据库 |
| Redis | 缓存点赞/关注状态 |
| JWT | 用户认证 |
| ffmpeg | 视频截帧生成封面 |

## 功能模块

- **用户模块**：注册、登录、获取用户信息
- **视频模块**：视频上传、Feed流（游标翻页）、发布列表
- **点赞模块**：点赞/取消点赞、点赞列表
- **评论模块**：发评论/删评论、评论列表
- **关注模块**：关注/取消关注、关注列表、粉丝列表、好友列表

## 项目结构

```
tiktok-simple/
├── cmd/main.go              # 程序入口
├── config/                  # 配置文件
├── internal/
│   ├── handler/             # HTTP 处理层
│   ├── service/             # 业务逻辑层
│   ├── repository/          # 数据库操作层
│   ├── model/               # 数据库模型
│   └── middleware/          # 中间件
├── pkg/
│   ├── database/            # 数据库连接
│   ├── jwt/                 # JWT 工具
│   └── response/            # 统一响应
└── storage/videos/          # 视频本地存储
```

## 配置说明

### 重要注意事项
- `config/config.yaml` 文件包含数据库密码、JWT密钥等敏感信息。
- 你需要**手动创建自己的配置文件**，不能直接使用示例文件。

### 快速配置步骤
1. 在 `config/` 目录下，复制示例配置文件：
   ```bash
   cp config/config.yaml.example config/config.yaml
2. 编辑 config/config.yaml，填入你的真实数据库密码、Redis 密码、JWT 密钥。
3. 保存后即可正常启动项目。

## 快速开始

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- ffmpeg

### 启动步骤

1. 克隆项目
   git clone https://github.com/howerguss/tiktok-simple.git
   cd tiktok-simple

2. 创建数据库
   CREATE DATABASE tiktok_simple CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

3. 修改配置
   编辑 config/config.yaml，填入你的 MySQL 和 Redis 密码

4. 启动服务
   go mod tidy
   go run cmd/main.go

5. 测试
   curl http://localhost:8080/ping

## 接口文档

### 用户
| 接口 | 方法 | 说明 |
|------|------|------|
| /douyin/user/register/ | POST | 注册 |
| /douyin/user/login/ | POST | 登录 |
| /douyin/user/ | GET | 获取用户信息（需登录）|

### 视频
| 接口 | 方法 | 说明 |
|------|------|------|
| /douyin/feed/ | GET | 获取Feed流 |
| /douyin/publish/action/ | POST | 上传视频（需登录）|
| /douyin/publish/list/ | GET | 发布列表（需登录）|

### 点赞
| 接口 | 方法 | 说明 |
|------|------|------|
| /douyin/favorite/action/ | POST | 点赞/取消（需登录）|
| /douyin/favorite/list/ | GET | 点赞列表（需登录）|

### 评论
| 接口 | 方法 | 说明 |
|------|------|------|
| /douyin/comment/action/ | POST | 发/删评论（需登录）|
| /douyin/comment/list/ | GET | 评论列表 |

### 关注
| 接口 | 方法 | 说明 |
|------|------|------|
| /douyin/relation/action/ | POST | 关注/取消（需登录）|
| /douyin/relation/follow/list/ | GET | 关注列表（需登录）|
| /douyin/relation/follower/list/ | GET | 粉丝列表（需登录）|
| /douyin/relation/friend/list/ | GET | 好友列表（需登录）|