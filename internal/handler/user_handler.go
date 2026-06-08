package handler

import (
	"tiktok-simple/internal/service"
	"tiktok-simple/pkg/response"

	"github.com/gin-gonic/gin"
)

// Register 注册接口 Handler
// POST /douyin/user/register/
// 参数：username, password（通过表单传参）
func Register(c *gin.Context) {
	// PostForm 从请求Body里取表单参数
	// 对应客户端发送的 Content-Type: application/x-www-form-urlencoded
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Handler层做基础的参数校验（格式、非空等）
	// 业务逻辑校验（用户名是否存在）放在Service层
	if username == "" || password == "" {
		response.Fail(c, 1, "用户名和密码不能为空")
		return
	}
	if len(password) < 6 {
		response.Fail(c, 1, "密码长度不能少于6位")
		return
	}

	// 调用Service层处理业务
	result, err := service.Register(username, password)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	// 返回成功响应
	response.Success(c, gin.H{
		"user_id": result.UserID,
		"token":   result.Token,
	})
}

// Login 登录接口 Handler
// POST /douyin/user/login/
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		response.Fail(c, 1, "用户名和密码不能为空")
		return
	}

	result, err := service.Login(username, password)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}

	response.Success(c, gin.H{
		"user_id": result.UserID,
		"token":   result.Token,
	})
}

// GetUserInfo 获取用户信息接口 Handler
// GET /douyin/user/
// 这个接口需要登录（被 AuthMiddleware 保护）
func GetUserInfo(c *gin.Context) {
	// 从 Gin Context 里取出中间件存入的 user_id
	// 这个值是 AuthMiddleware 解析完token后通过 c.Set("user_id", ...) 存进来的
	// GetUint 取出uint类型的值
	userID := c.GetUint("user_id")

	user, err := service.GetUserInfo(userID)
	if err != nil {
		response.Fail(c, 1, "用户不存在")
		return
	}

	response.Success(c, gin.H{"user": user})
}
