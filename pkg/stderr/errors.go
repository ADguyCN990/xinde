package stderr

// 预定义的业务错误code
const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
)

const (
	ErrorDbNil            = "dao或数据库连接为空"
	ErrorUserAlreadyExist = "用户已经存在"
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
