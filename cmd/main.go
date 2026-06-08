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
	"tiktok-simple/pkg/storage"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	database.InitMySQL()
	database.InitRedis()
	database.InitMongoDB() // 新增：初始化 MongoDB
	storage.InitMinio()    // 新增：初始化 MinIO

	r := gin.New()
	r.Use(middleware.CustomRecovery())
	r.Use(middleware.RequestLogger())
	r.Use(cors.Default())

	r.Static("/static/videos", "./storage/videos")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 不需要登录
	userGroup := r.Group("/douyin/user")
	{
		userGroup.POST("/register/", handler.Register)
		userGroup.POST("/login/", handler.Login)
	}
	r.GET("/douyin/feed/", handler.Feed)
	r.GET("/douyin/comment/list/", handler.CommentList)

	// 需要登录
	authGroup := r.Group("/douyin")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("/user/", handler.GetUserInfo)

		authGroup.POST("/publish/action/", handler.Publish)
		authGroup.GET("/publish/list/", handler.PublishList)

		authGroup.POST("/favorite/action/", handler.FavoriteAction)
		authGroup.GET("/favorite/list/", handler.FavoriteList)

		authGroup.POST("/comment/action/", handler.CommentAction)

		authGroup.POST("/relation/action/", handler.RelationAction)
		authGroup.GET("/relation/follow/list/", handler.FollowList)
		authGroup.GET("/relation/follower/list/", handler.FollowerList)
		authGroup.GET("/relation/friend/list/", handler.FriendList)

		// 新增：消息接口
		authGroup.POST("/message/action/", handler.MessageAction)
		authGroup.GET("/message/chat/", handler.MessageChat)
	}

	addr := fmt.Sprintf(":%d", config.Global.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		fmt.Printf("服务启动在 %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\n收到退出信号，正在优雅关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭失败: %v", err)
	}
	fmt.Println("服务已安全退出")
}
