package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// GenerateCover 从视频文件截取第一帧作为封面图
// videoPath: 视频文件的完整路径，比如 ./storage/videos/xxx.mp4
// 返回值: 封面图片的完整路径，比如 ./storage/videos/xxx_cover.jpg
func GenerateCover(videoPath string) (string, error) {
	// 封面路径 = 视频路径去掉扩展名 + "_cover.jpg"
	// filepath.Ext 获取扩展名，比如 ".mp4"
	// videoPath[:len(videoPath)-len(ext)] 去掉扩展名
	ext := filepath.Ext(videoPath)
	coverPath := videoPath[:len(videoPath)-len(ext)] + "_cover.jpg"

	// 调用 ffmpeg 命令截取第一帧
	// -i videoPath     输入文件
	// -ss 00:00:01     从第1秒开始截取（第0秒可能是黑屏）
	// -vframes 1       只截取1帧
	// -y               如果文件已存在则覆盖
	// coverPath        输出文件路径
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-y",
		coverPath,
	)

	// CombinedOutput 执行命令并捕获标准输出+错误输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 如果截取失败（比如视频太短），删除可能生成的空文件
		os.Remove(coverPath)
		return "", fmt.Errorf("ffmpeg截取封面失败: %w, output: %s", err, string(output))
	}

	return coverPath, nil
}
