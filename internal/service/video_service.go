package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"tiktok-simple/config"
	"tiktok-simple/internal/model"
	"tiktok-simple/internal/repository"
	"tiktok-simple/pkg/util"
	"time"

	"github.com/google/uuid" // 用来生成唯一文件名，需要安装
)

// PublishVideo 处理视频上传业务
// userID: 上传者的ID（从JWT token里拿到）
// title: 视频标题
// fileHeader: 上传的视频文件
func PublishVideo(userID uint, title string, fileHeader *multipart.FileHeader) error {
	storagePath := config.Global.Storage.Path

	// 确保存储目录存在
	// MkdirAll 类似 mkdir -p，目录已存在不会报错
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return fmt.Errorf("创建存储目录失败: %w", err)
	}

	// 生成唯一文件名，防止文件名冲突
	// uuid.New().String() 生成类似 "550e8400-e29b-41d4-a716-446655440000" 的唯一字符串
	// filepath.Ext 获取原文件扩展名，比如 ".mp4"
	videoFileName := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	videoFilePath := filepath.Join(storagePath, videoFileName)

	// TODO: 打开上传的文件
	// 提示：src, err := fileHeader.Open()
	src, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close() // 函数结束时关闭文件，释放资源

	// TODO: 创建目标文件
	// 提示：dst, err := os.Create(videoFilePath)
	dst, err := os.Create(videoFilePath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	// TODO: 把上传的文件内容复制到目标文件
	// 提示：用 io.Copy(dst, src) 来复制
	// 需要在import里加 "io"
	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("保存视频失败: %w", err)
	}

	// 用ffmpeg截取视频封面
	// 截取失败不影响上传，封面用空字符串
	coverFilePath := ""
	coverFileName := ""
	cover, err := util.GenerateCover(videoFilePath)
	if err == nil {
		coverFilePath = cover
		coverFileName = filepath.Base(coverFilePath) // Base 取路径最后一段，即文件名
	}

	// 构建访问URL
	// 这里用相对路径存储，返回给前端时通过接口提供访问
	// 比如视频文件名是 xxx.mp4，访问URL是 /static/videos/xxx.mp4
	playURL := "/static/videos/" + videoFileName
	coverURL := ""
	if coverFileName != "" {
		coverURL = "/static/videos/" + coverFileName
	}
	_ = coverFilePath // 避免未使用变量报错

	// 写入数据库
	video := &model.Video{
		UserID:   userID,
		Title:    title,
		PlayURL:  playURL,
		CoverURL: coverURL,
	}
	return repository.CreateVideo(video)
}

// GetFeed 获取Feed流
// latestTime: 时间戳，返回这个时间之前的视频（实现翻页）
func GetFeed(latestTime int64) ([]model.Video, error) {
	// 如果没传时间戳，默认用当前时间（返回最新的视频）
	if latestTime == 0 {
		latestTime = time.Now().Unix()
	}
	return repository.GetFeedVideos(latestTime, 10) // 每次返回10条
}

// GetPublishList 获取用户发布的视频列表
func GetPublishList(userID uint) ([]model.Video, error) {
	return repository.GetVideosByUserID(userID)
}
