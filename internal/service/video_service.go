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
	"tiktok-simple/pkg/storage"
	"tiktok-simple/pkg/util"
	"time"

	"github.com/google/uuid"
)

// PublishVideo 处理视频上传
// 流程：保存临时文件 → 截取封面 → 上传视频到MinIO → 上传封面到MinIO → 写数据库 → 清理临时文件
func PublishVideo(userID uint, title string, fileHeader *multipart.FileHeader) error {
	// 本地临时目录，用于ffmpeg截帧（ffmpeg需要读本地文件）
	tmpPath := config.Global.Storage.Path
	if err := os.MkdirAll(tmpPath, 0755); err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 生成唯一文件名
	videoFileName := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	videoTmpPath := filepath.Join(tmpPath, videoFileName) // 本地临时路径

	// 第一步：把上传的视频保存到本地临时文件
	src, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(videoTmpPath)
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("保存临时文件失败: %w", err)
	}
	dst.Close() // 必须先关闭再让ffmpeg读取

	// 函数结束时清理所有临时文件（无论成功还是失败）
	defer os.Remove(videoTmpPath)

	// 第二步：用ffmpeg截取封面（截本地临时文件）
	coverTmpPath := ""
	coverURL := ""
	cover, err := util.GenerateCover(videoTmpPath)
	if err == nil {
		coverTmpPath = cover
		defer os.Remove(coverTmpPath) // 函数结束时清理封面临时文件
	}

	// 第三步：上传视频到 MinIO
	// objectName 是文件在MinIO桶里的"路径"，用 videos/ 前缀做分类
	videoObjectName := "videos/" + videoFileName
	playURL, err := storage.UploadFile(videoObjectName, videoTmpPath, "video/mp4")
	if err != nil {
		return fmt.Errorf("上传视频到MinIO失败: %w", err)
	}

	// 第四步：上传封面到 MinIO（如果截帧成功了）
	if coverTmpPath != "" {
		coverFileName := filepath.Base(coverTmpPath)
		coverObjectName := "covers/" + coverFileName // 封面用 covers/ 前缀
		coverURL, _ = storage.UploadFile(coverObjectName, coverTmpPath, "image/jpeg")
		// 封面上传失败不影响主流程，coverURL 为空字符串
	}

	// 第五步：写入数据库
	video := &model.Video{
		UserID:   userID,
		Title:    title,
		PlayURL:  playURL,  // MinIO的完整URL，比如 http://localhost:9000/tiktok/videos/xxx.mp4
		CoverURL: coverURL, // MinIO的完整URL
	}
	return repository.CreateVideo(video)
}

// GetFeed 获取Feed流（不变）
func GetFeed(latestTime int64) ([]model.Video, error) {
	if latestTime == 0 {
		latestTime = time.Now().Unix()
	}
	return repository.GetFeedVideos(latestTime, 10)
}

// GetPublishList 获取用户发布列表（不变）
func GetPublishList(userID uint) ([]model.Video, error) {
	return repository.GetVideosByUserID(userID)
}
