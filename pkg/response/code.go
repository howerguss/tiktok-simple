package response

// 错误码规范：
// 0        = 成功
// 1xxx     = 参数错误
// 2xxx     = 业务错误（用户名已存在、密码错误等）
// 4001     = 未登录
// 4003     = 无权限
// 5000     = 服务器内部错误
//
// 好处：
// - 前端可以根据错误码做不同处理（比如 4001 自动跳转登录页）
// - 日志监控可以按错误码统计，快速发现问题
// - 多端（iOS/Android/Web）统一处理逻辑

const (
	// 成功
	CodeSuccess = 0

	// 参数错误 1xxx
	CodeInvalidParam = 1001 // 参数格式错误
	CodeMissingParam = 1002 // 缺少必要参数
	CodeParamTooLong = 1003 // 参数超长

	// 业务错误 2xxx
	CodeUserNotFound    = 2001 // 用户不存在
	CodeUserExists      = 2002 // 用户名已存在
	CodeWrongPassword   = 2003 // 密码错误
	CodeVideoNotFound   = 2004 // 视频不存在
	CodeAlreadyFavorite = 2005 // 已经点赞过了
	CodeNotFavorite     = 2006 // 还没有点赞
	CodeAlreadyFollow   = 2007 // 已经关注过了
	CodeNotFollow       = 2008 // 还没有关注
	CodeSelfFollow      = 2009 // 不能关注自己
	CodeCommentNotFound = 2010 // 评论不存在或无权删除

	// 认证错误 4xxx
	CodeUnauthorized = 4001 // 未登录
	CodeForbidden    = 4003 // 无权限

	// 服务器错误 5xxx
	CodeServerError = 5000 // 服务器内部错误
)

// CodeMsg 错误码对应的默认消息
// 统一管理，避免同一个错误在不同地方文案不一致
var CodeMsg = map[int32]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "参数格式错误",
	CodeMissingParam:    "缺少必要参数",
	CodeParamTooLong:    "参数超长",
	CodeUserNotFound:    "用户不存在",
	CodeUserExists:      "用户名已存在",
	CodeWrongPassword:   "用户名或密码错误",
	CodeVideoNotFound:   "视频不存在",
	CodeAlreadyFavorite: "已经点赞过了",
	CodeNotFavorite:     "还没有点赞",
	CodeAlreadyFollow:   "已经关注过了",
	CodeNotFollow:       "还没有关注",
	CodeSelfFollow:      "不能关注自己",
	CodeCommentNotFound: "评论不存在或无权删除",
	CodeUnauthorized:    "请先登录",
	CodeForbidden:       "无权限",
	CodeServerError:     "服务器内部错误",
}

// FailWithCode 用错误码返回失败响应
// 自动填充对应的错误消息，不需要手动传 msg
func FailWithCode(c interface{ JSON(int, interface{}) }, code int32) {
	// 这里为了复用方便，直接在 response.go 里扩展一个方法
	// 实际使用见 response.go 的更新
}
