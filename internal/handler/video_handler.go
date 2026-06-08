package handler

import (
	"strconv"
	"tiktok-simple/internal/service"
	"tiktok-simple/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

// Publish 视频上传接口
// POST /douyin/publish/action/
// 需要登录（在路由注册时加了AuthMiddleware）
func Publish(c *gin.Context) {
	// 从中间件里取当前登录用户的ID
	userID := c.GetUint("user_id")
	title := c.PostForm("title")

	if title == "" {
		response.Fail(c, 1, "标题不能为空")
		return
	}

	// FormFile 从请求里取上传的文件
	// "data" 是前端上传时用的字段名
	fileHeader, err := c.FormFile("data")
	if err != nil {
		response.Fail(c, 1, "请选择要上传的视频文件")
		return
	}

	// 简单校验文件大小，限制100MB
	// fileHeader.Size 单位是字节，100MB = 100 * 1024 * 1024
	if fileHeader.Size > 100*1024*1024 {
		response.Fail(c, 1, "视频文件不能超过100MB")
		return
	}

	if err := service.PublishVideo(userID, title, fileHeader); err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	response.Success(c, gin.H{})
}

// Feed 获取视频流接口
// GET /douyin/feed/
// 不需要登录
func Feed(c *gin.Context) {
	// latest_time 是时间戳，用于翻页
	// 第一次请求不传，之后传上次返回的最早一条视频的时间戳
	latestTimeStr := c.Query("latest_time")
	var latestTime int64

	if latestTimeStr != "" {
		// ParseInt 把字符串转成int64
		// 第二个参数10表示十进制，第三个参数64表示int64
		latestTime, _ = strconv.ParseInt(latestTimeStr, 10, 64)
	} else {
		latestTime = time.Now().Unix()
	}

	videos, err := service.GetFeed(latestTime)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	// nextTime 是返回的视频里最早那条的时间戳
	// 前端下次请求时把这个时间戳传过来，就能拿到更早的视频（翻页）
	var nextTime int64
	if len(videos) > 0 {
		nextTime = videos[len(videos)-1].CreatedAt.Unix()
	} else {
		nextTime = time.Now().Unix()
	}

	response.Success(c, gin.H{
		"video_list": videos,
		"next_time":  nextTime,
	})
}

// PublishList 获取用户发布列表接口
// GET /douyin/publish/list/
// 需要登录
func PublishList(c *gin.Context) {
	// 这个接口查的是自己发布的视频，所以直接从token里取user_id
	userID := c.GetUint("user_id")

	videos, err := service.GetPublishList(userID)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	response.Success(c, gin.H{"video_list": videos})
}
