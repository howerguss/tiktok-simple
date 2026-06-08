package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

// RateLimiter 全局令牌桶限流中间件
//
// 令牌桶算法原理：
// - 有一个容量为 burst 的桶
// - 系统以固定速率（r个/秒）往桶里放令牌
// - 每个请求进来时从桶里取一个令牌
// - 取到令牌 → 放行；取不到 → 拒绝（返回429）
// - 桶满了就不再放令牌（允许短时间的突发流量，最多 burst 个）
//
// 举例：rate.NewLimiter(100, 200)
// - 每秒产生100个令牌（即每秒最多处理100个请求）
// - 桶容量200（允许瞬间突发200个请求）
//
// 对比漏桶算法：
// - 漏桶：严格匀速处理，不允许突发
// - 令牌桶：允许短时突发（消耗桶里积累的令牌），更适合Web场景
func RateLimiter(r rate.Limit, burst int) gin.HandlerFunc {
	// 创建一个全局共享的限流器
	// 注意：这是整个服务共享的，不是每个用户单独的
	// 如果要做"每个用户每秒最多10次"，需要用 map 存每个用户的限流器
	limiter := rate.NewLimiter(r, burst)

	return func(c *gin.Context) {
		// Allow() 尝试取一个令牌
		// 如果桶里有令牌，返回 true；没有令牌，返回 false（不会阻塞等待）
		if !limiter.Allow() {
			// 返回 429 Too Many Requests
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status_code": 4029,
				"status_msg":  "请求太频繁，请稍后再试",
			})
			return
		}
		c.Next()
	}
}
