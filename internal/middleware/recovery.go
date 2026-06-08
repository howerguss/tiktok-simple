package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"tiktok-simple/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

// CustomRecovery 自定义 panic 恢复中间件
// 作用：捕获任意 handler 里的 panic，防止整个服务崩溃
//
// 什么时候会 panic？
// - 代码里有 bug，比如空指针解引用（nil pointer dereference）
// - 数组越界
// - 类型断言失败
// 没有 Recovery 中间件的话，一个请求的 panic 会导致整个 goroutine 崩溃
// 进而可能影响整个服务进程
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// recover() 捕获 panic，如果没有 panic 则返回 nil
			if err := recover(); err != nil {
				// 打印 panic 的详细堆栈信息，方便排查 bug
				// debug.Stack() 返回当前 goroutine 的完整调用栈
				fmt.Printf("[PANIC] %s\n%s\n%s\n",
					time.Now().Format("2006-01-02 15:04:05"),
					fmt.Sprintf("panic: %v", err),
					string(debug.Stack()),
				)

				// 向客户端返回 500 错误，不暴露内部错误详情
				// 真实错误信息只记录在日志里，不返回给客户端（安全考虑）
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"status_code": response.CodeServerError,
					"status_msg":  "服务器内部错误，请稍后重试",
				})
			}
		}()

		// 继续执行后续中间件和 handler
		c.Next()
	}
}
