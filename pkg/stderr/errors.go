package stderr

// 预定义的业务错误code
const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
)

const (
	ErrorInternalServerError = "服务器内部错误"
)

const (
	ErrorDbNil = "dao或数据库连接为空"
)

// pagination
const (
	ErrorOverLargePage = "查询页数过大"
	ErrorOverSmallPage = "查询页数过小"
)

// account
const (
	ErrorUserNotFound     = "用户不存在 "
	ErrorUserAlreadyExist = "用户已经存在"
	ErrorUserUnauthorized = "用户名或密码错误"
	ErrorUserNotPass      = "用户尚未被管理员通过注册申请"
	ErrorUserBanned       = "用户已经被管理员拒绝注册申请"
	ErrorUserPassed       = "用户已经被管理员批准注册申请"
	ErrorUserIDInvalid    = "无效的用户ID格式"
)

// JWT token
const (
	ErrorTokenExpired     = "token已过期"
	ErrorTokenNotValidYet = "token尚未生效"
	ErrorTokenMalFormed   = "token格式错误"
	ErrorTokenInvalid     = "token解析失败"
)

// MsgFlags 预定义的业务错误msg
var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
