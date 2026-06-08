package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger 请求日志中间件
// 记录：请求方法、路径、状态码、耗时、客户端IP
//
// 为什么需要自定义日志？
// gin.Default() 自带的日志格式比较简单，生产环境通常需要：
// - 记录请求耗时（定位慢接口）
// - 记录用户ID（追踪某个用户的操作）
// - 输出为 JSON 格式（方便日志系统解析）
// 这里我们做一个简化版，打印关键信息即可
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 执行后续中间件和 handler
		c.Next()

		// handler 执行完毕后，计算耗时
		duration := time.Since(startTime)

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 获取当前登录用户ID（如果有的话）
		userID, _ := c.Get("user_id")

		// 打印日志
		// 格式：[时间] 方法 路径 状态码 耗时 IP 用户ID
		fmt.Printf("[GIN] %s | %d | %v | %s | %s %s | userID:%v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			statusCode,
			duration,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
			userID,
		)
	}
}
