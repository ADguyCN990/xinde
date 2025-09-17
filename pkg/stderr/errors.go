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

// company
const (
	ErrorCompanyNotFound  = "公司不存在"
	ErrorCompanyIDInvalid = "无效的公司ID格式"
)

// attachment
const (
	ErrorAttachmentNotFound       = "附件不存在"
	ErrorAttachmentIDInvalid      = "无效的附件ID格式"
	ErrorAttachmentNotFoundOnDesk = "附件不存在于磁盘上"
)

// group
const (
	ErrorGroupNotFound             = "分组不存在"
	ErrorGroupIDInvalid            = "无效的分组ID格式"
	ErrorCannotMoveGroupIntoItself = "所更改的父级分组不能是其子孙分组"
	ErrorRootGroupCannotBeDeleted  = "root分组不能被删除"
)

// device
const (
	ErrorDeviceNotFound  = "设备类型不存在"
	ErrorDeviceIDInvalid = "无效的设备类型ID格式"
)

// JWT token
const (
	ErrorTokenExpired     = "token已过期"
	ErrorTokenNotValidYet = "token尚未生效"
	ErrorTokenMalFormed   = "token格式错误"
	ErrorTokenInvalid     = "token解析失败"
	ErrorTokenNotAdmin    = "非管理员,权限不足"
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
