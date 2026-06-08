package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tiktok-simple/config"
	"tiktok-simple/internal/handler"
	"tiktok-simple/internal/middleware"
	"tiktok-simple/pkg/database"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// 1. 初始化配置和数据库
	config.Init()
	database.InitMySQL()
	database.InitRedis()

	// 2. 创建 Gin 引擎
	// 注意：这里用 gin.New() 而不是 gin.Default()
	// gin.Default() 会自动加载 Logger 和 Recovery
	// 我们用自定义的，所以用 gin.New() 手动加载
	r := gin.New()

	// 3. 注册全局中间件
	// 顺序很重要：Recovery 要放最前面，保证任何 panic 都能被捕获
	r.Use(middleware.CustomRecovery())                    // 全局 panic 捕获
	r.Use(middleware.RequestLogger())                     // 请求日志
	r.Use(cors.Default())                                 // 跨域
	r.Use(middleware.RateLimiter(rate.Limit(1000), 2000)) // 全局限流：每秒1000请求，最大突发2000

	// 4. 静态文件服务
	r.Static("/static/videos", "./storage/videos")

	// 5. 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// ==================== 不需要登录的路由 ====================
	userGroup := r.Group("/douyin/user")
	{
		userGroup.POST("/register/", handler.Register)
		userGroup.POST("/login/", handler.Login)
	}

	r.GET("/douyin/feed/", handler.Feed)
	r.GET("/douyin/comment/list/", handler.CommentList)

	// ==================== 需要登录的路由 ====================
	authGroup := r.Group("/douyin")
	authGroup.Use(middleware.AuthMiddleware())
	{
		// 用户
		authGroup.GET("/user/", handler.GetUserInfo)

		// 视频
		authGroup.POST("/publish/action/", handler.Publish)
		authGroup.GET("/publish/list/", handler.PublishList)

		// 点赞
		authGroup.POST("/favorite/action/", handler.FavoriteAction)
		authGroup.GET("/favorite/list/", handler.FavoriteList)

		// 评论
		authGroup.POST("/comment/action/", handler.CommentAction)

		// 关注
		authGroup.POST("/relation/action/", handler.RelationAction)
		authGroup.GET("/relation/follow/list/", handler.FollowList)
		authGroup.GET("/relation/follower/list/", handler.FollowerList)
		authGroup.GET("/relation/friend/list/", handler.FriendList)
	}

	// 6. 用 http.Server 包装 Gin，支持优雅关闭
	// 为什么不直接用 r.Run()？
	// r.Run() 内部调用的是 http.ListenAndServe，没有提供关闭的控制权
	// 用 http.Server 可以调用 server.Shutdown() 来优雅关闭
	addr := fmt.Sprintf(":%d", config.Global.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
		// 超时配置，防止慢客户端占用连接资源
		ReadTimeout:  10 * time.Second, // 读取请求的超时时间
		WriteTimeout: 30 * time.Second, // 写入响应的超时时间
	}

	// 7. 在新的 goroutine 里启动服务（不阻塞主 goroutine）
	go func() {
		fmt.Printf("服务启动在 %s\n", addr)
		// ListenAndServe 正常关闭时返回 http.ErrServerClosed，不是真正的错误
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 8. 监听系统退出信号（优雅关闭的核心）
	// make(chan os.Signal, 1) 创建一个带缓冲的 channel，避免信号丢失
	quit := make(chan os.Signal, 1)
	// signal.Notify 把指定的系统信号转发到 channel
	// syscall.SIGINT  = Ctrl+C 产生的信号
	// syscall.SIGTERM = kill 命令产生的信号（Docker 停止容器时发送这个）
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待信号（主 goroutine 在这里等着）
	<-quit
	fmt.Println("\n收到退出信号，正在优雅关闭服务...")

	// 9. 执行优雅关闭
	// context.WithTimeout 设置最长等待时间：5秒内处理完已有请求，否则强制退出
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown 会：
	// - 停止接受新的连接和请求
	// - 等待已有请求处理完毕
	// - 超过 ctx 时间限制后强制关闭
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭失败: %v", err)
	}

	fmt.Println("服务已安全退出")
}
