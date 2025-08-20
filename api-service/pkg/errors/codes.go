package errors

// 错误码常量定义
const (
	// 通用错误码 (10000-19999)
	CodeSuccess         = 0
	CodeInternalError   = 10001
	CodeInvalidRequest  = 10002
	CodeUnauthorized    = 10003
	CodeForbidden       = 10004
	CodeNotFound        = 10005
	CodeValidationError = 10006

	// 用户相关错误码 (20000-29999)
	CodeUserNotFound         = 20001
	CodeUserAlreadyExists    = 20002
	CodeInvalidCredentials   = 20003
	CodeUserInactive         = 20004
	CodeInvalidPassword      = 20005
	CodePasswordTooWeak      = 20006
	CodeEmailAlreadyExists   = 20007
	CodeInvalidEmail         = 20008
	CodeUsernameReserved     = 20009
	CodeUserQuotaExceeded    = 20010

	// 应用相关错误码 (30000-39999)
	CodeAppNotFound          = 30001
	CodeAppAlreadyExists     = 30002
	CodeAppDeployFailed      = 30003
	CodeAppPortConflict      = 30004
	CodeAppResourcesInsufficient = 30005
)

// 错误码对应的默认消息
var CodeMessages = map[int]string{
	CodeSuccess:         "成功",
	CodeInternalError:   "内部服务器错误",
	CodeInvalidRequest:  "请求参数无效",
	CodeUnauthorized:    "未授权访问",
	CodeForbidden:       "禁止访问",
	CodeNotFound:        "资源不存在",
	CodeValidationError: "数据验证失败",

	// 用户相关
	CodeUserNotFound:         "用户不存在",
	CodeUserAlreadyExists:    "用户名已存在",
	CodeInvalidCredentials:   "用户名或密码错误",
	CodeUserInactive:         "用户账号未激活",
	CodeInvalidPassword:      "密码错误",
	CodePasswordTooWeak:      "密码强度不足",
	CodeEmailAlreadyExists:   "邮箱已存在",
	CodeInvalidEmail:         "邮箱格式无效",
	CodeUsernameReserved:     "用户名为系统保留字",
	CodeUserQuotaExceeded:    "用户配额已超限",

	// 应用相关
	CodeAppNotFound:              "应用不存在",
	CodeAppAlreadyExists:         "应用已存在",
	CodeAppDeployFailed:          "应用部署失败",
	CodeAppPortConflict:          "端口冲突",
	CodeAppResourcesInsufficient: "资源不足",
}
