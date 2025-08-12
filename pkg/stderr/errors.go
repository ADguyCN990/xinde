package stderr

// 预定义的业务错误code
const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
)

const (
	ERROR_DB_NIL             = "dao或数据库连接为空"
	ERROR_USER_ALREADY_EXIST = "用户已经存在"
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
