package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// Success 返回成功响应，data 里的字段会合并到顶层
func Success(c *gin.Context, data gin.H) {
	resp := gin.H{
		"status_code": CodeSuccess,
		"status_msg":  "success",
	}
	for k, v := range data {
		resp[k] = v
	}
	c.JSON(http.StatusOK, resp)
}

// Fail 返回失败响应（传入自定义msg）
func Fail(c *gin.Context, code int32, msg string) {
	c.JSON(http.StatusOK, Response{
		StatusCode: code,
		StatusMsg:  msg,
	})
}

// FailCode 用错误码返回失败响应（自动查找对应msg）
// 推荐使用这个，保证同一错误文案统一
func FailCode(c *gin.Context, code int32) {
	msg, ok := CodeMsg[code]
	if !ok {
		msg = "未知错误"
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: code,
		StatusMsg:  msg,
	})
}

// ServerError 返回服务器内部错误（不暴露详细错误信息给客户端）
// 详细错误应该打印到日志里
func ServerError(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		StatusCode: CodeServerError,
		StatusMsg:  CodeMsg[CodeServerError],
	})
}
