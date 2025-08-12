package errors

// 预定义的业务错误code
const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
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
