# tiktok-simple — 极简抖音后端

基于 Go + Gin 实现的短视频社交平台后端，覆盖抖音核心功能。

## 技术栈

| 技术    | 版本    | 用途                                     |
| ------- | ------- | ---------------------------------------- |
| Go      | 1.24.8  | 主语言                                   |
| Gin     | v1.11.0 | HTTP 框架                                |
| GORM    | v1.31.1 | ORM，操作 MySQL                          |
| MySQL   | 8.0+    | 主数据库，存储用户/视频/关系等结构化数据 |
| Redis   | 6.0+    | 缓存点赞/关注状态，加速高频查询          |
| MongoDB | 7.0+    | 存储聊天消息（文档型，高频写入）         |
| MinIO   | 最新版  | 视频与封面的分布式对象存储               |
| JWT     | v5.3.1  | 用户无状态认证                           |
| ffmpeg  | -       | 视频截帧生成封面                         |

## 功能模块

- **用户模块**：注册、登录、获取用户信息
- **视频模块**：视频上传（MinIO存储）、Feed流（游标翻页）、发布列表
- **点赞模块**：点赞/取消点赞、点赞列表
- **评论模块**：发评论/删评论、评论列表
- **关注模块**：关注/取消关注、关注列表、粉丝列表、好友列表
- **消息模块**：发送聊天消息、获取聊天记录（MongoDB存储）

## 项目结构

```
tiktok-simple/
├── cmd/
│   └── main.go                      # 程序入口，初始化、路由、优雅关闭
├── config/
│   ├── config.go                    # 配置加载
│   ├── config.yaml                  # 配置文件
├── internal/
│   ├── handler/                     # HTTP 处理层（解析参数、返回响应）
│   │   ├── user_handler.go
│   │   ├── video_handler.go
│   │   ├── favorite_handler.go
│   │   ├── comment_handler.go
│   │   ├── relation_handler.go
│   │   └── message_handler.go
│   ├── service/                     # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── video_service.go
│   │   ├── favorite_service.go
│   │   ├── comment_service.go
│   │   ├── relation_service.go
│   │   └── message_service.go
│   ├── repository/                  # 数据库操作层
│   │   ├── user_repository.go
│   │   ├── video_repository.go
│   │   ├── favorite_repository.go
│   │   ├── comment_repository.go
│   │   ├── relation_repository.go
│   │   └── message_repository.go
│   ├── model/                       # 数据模型
│   │   ├── user.go
│   │   ├── video.go
│   │   ├── favorite.go
│   │   ├── comment.go
│   │   ├── relation.go
│   │   └── message.go               # MongoDB 文档模型
│   └── middleware/                  # 中间件
│       ├── auth.go                  # JWT 鉴权
│       ├── recovery.go              # panic 捕获
│       ├── logger.go                # 请求日志
│       └── ratelimit.go             # 令牌桶限流
├── pkg/
│   ├── database/
│   │   ├── mysql.go                 # MySQL 连接
│   │   ├── redis.go                 # Redis 连接
│   │   └── mongodb.go               # MongoDB 连接
│   ├── storage/
│   │   └── minio.go                 # MinIO 客户端
│   ├── jwt/
│   │   └── jwt.go                   # JWT 生成和解析
│   ├── response/
│   │   ├── response.go              # 统一响应格式
│   │   └── code.go                  # 错误码定义
│   └── util/
│       └── video.go                 # ffmpeg 截帧工具
├── storage/
│   └── videos/                      # 视频临时目录（上传 MinIO 前的暂存）
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- MongoDB 7.0+
- MinIO（本地部署）
- ffmpeg

## 快速开始

### 第一步：克隆项目

```bash
git clone https://github.com/howerguss/tiktok-simple.git
cd tiktok-simple
```

### 第二步：安装依赖服务

**安装 ffmpeg：**

```bash
sudo apt install ffmpeg -y
```

**安装并启动 MinIO：**

```bash
wget https://dl.min.io/server/minio/release/linux-amd64/minio
chmod +x minio && sudo mv minio /usr/local/bin/
mkdir -p ~/minio-data
export MINIO_ROOT_USER=minioadmin
export MINIO_ROOT_PASSWORD=minioadmin123
minio server ~/minio-data --console-address ":9001"
```

启动后用 mc 工具创建并设置公开桶：

```bash
wget https://dl.min.io/client/mc/release/linux-amd64/mc
chmod +x mc && sudo mv mc /usr/local/bin/
mc alias set local http://localhost:9000 minioadmin minioadmin123
mc mb local/tiktok
mc anonymous set download local/tiktok
```

**安装并启动 MongoDB：**

```bash
sudo apt-get install -y mongodb-org
sudo systemctl start mongod
sudo systemctl enable mongod
```

### 第三步：创建 MySQL 数据库

```sql
CREATE DATABASE tiktok_simple CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 第四步：配置文件

```bash
cp config/config.yaml.example config/config.yaml
```

编辑 `config/config.yaml`，填入你自己的密码和配置：

```yaml
server:
  port: 8080

database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: "你的MySQL密码"
  dbname: tiktok_simple

redis:
  host: 127.0.0.1
  port: 6379
  password: "你的Redis密码"
  db: 0

jwt:
  secret: "your-secret-key"
  expire: 72

storage:
  path: "./storage/videos"

minio:
  endpoint: "localhost:9000"
  access_key_id: "minioadmin"
  secret_access_key: "minioadmin123"
  bucket_name: "tiktok"
  use_ssl: false
  base_url: "http://localhost:9000"

mongodb:
  uri: "mongodb://localhost:27017"
  database: "tiktok_simple"
```

### 第五步：启动服务

```bash
go mod tidy
go run cmd/main.go
```

看到以下输出说明启动成功：

```
MySQL 连接成功
Redis 连接成功
MongoDB 连接成功
MinIO 连接成功
服务启动在 :8080
```

### 第六步：验证

```bash
curl http://localhost:8080/ping
# 返回 {"message":"pong"} 表示成功
```

## 接口文档

### 用户

| 接口                   | 方法 | 是否需要登录 | 说明         |
| ---------------------- | ---- | ------------ | ------------ |
| /douyin/user/register/ | POST | ❌            | 注册         |
| /douyin/user/login/    | POST | ❌            | 登录         |
| /douyin/user/          | GET  | ✅            | 获取用户信息 |

### 视频

| 接口                    | 方法 | 是否需要登录 | 说明         |
| ----------------------- | ---- | ------------ | ------------ |
| /douyin/feed/           | GET  | ❌            | 获取Feed流   |
| /douyin/publish/action/ | POST | ✅            | 上传视频     |
| /douyin/publish/list/   | GET  | ✅            | 我的发布列表 |

### 点赞

| 接口                     | 方法 | 是否需要登录 | 说明          |
| ------------------------ | ---- | ------------ | ------------- |
| /douyin/favorite/action/ | POST | ✅            | 点赞/取消点赞 |
| /douyin/favorite/list/   | GET  | ✅            | 我的点赞列表  |

### 评论

| 接口                    | 方法 | 是否需要登录 | 说明      |
| ----------------------- | ---- | ------------ | --------- |
| /douyin/comment/action/ | POST | ✅            | 发/删评论 |
| /douyin/comment/list/   | GET  | ❌            | 评论列表  |

### 关注

| 接口                            | 方法 | 是否需要登录 | 说明                 |
| ------------------------------- | ---- | ------------ | -------------------- |
| /douyin/relation/action/        | POST | ✅            | 关注/取消关注        |
| /douyin/relation/follow/list/   | GET  | ✅            | 我关注的人           |
| /douyin/relation/follower/list/ | GET  | ✅            | 粉丝列表             |
| /douyin/relation/friend/list/   | GET  | ✅            | 好友列表（互相关注） |

### 消息

| 接口                    | 方法 | 是否需要登录 | 说明         |
| ----------------------- | ---- | ------------ | ------------ |
| /douyin/message/action/ | POST | ✅            | 发送消息     |
| /douyin/message/chat/   | GET  | ✅            | 获取聊天记录 |
