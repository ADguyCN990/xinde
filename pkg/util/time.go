package util

import "time"

// 定义你需要的日期时间格式为一个常量，这是一个非常好的实践。
// 这可以保证整个项目中使用的格式都是统一的，并且易于修改。
const (
	// StandardDateTimeFormat YYYY-MM-DD HH:MM:SS 格式
	StandardDateTimeFormat = "2006-01-02 15:04:05"

	// StandardDateFormat YYYY-MM-DD 格式 (常用)
	StandardDateFormat = "2006-01-02"

	// StandardTimeFormat HH:MM:SS 格式 (常用)
	StandardTimeFormat = "15:04:05"
)

// FormatTimeToStandardString 将 time.Time 对象格式化为 "YYYY-MM-DD HH:MM:SS" 格式的字符串。
// 它会根据 time.Time 对象自身所在的时区进行格式化。
func FormatTimeToStandardString(t time.Time) string {
	return t.Format(StandardDateTimeFormat)
}

// FormatNullableTimeToStandardString 安全地将一个 *time.Time 指针格式化为字符串。
// 如果指针为 nil，则返回一个空字符串 ""。
func FormatNullableTimeToStandardString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(StandardDateTimeFormat)
}
