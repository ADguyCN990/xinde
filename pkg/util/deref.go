package util

import "time"

// DerefString 安全地解引用一个字符串指针。
// 如果指针为 nil，返回一个空字符串。
func DerefString(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}

// DerefUint 安全地解引用一个 uint 指针。
// 如果指针为 nil，返回 0。
func DerefUint(p *uint) uint {
	if p != nil {
		return *p
	}
	return 0
}

// DerefInt64 安全地解引用一个 int64 指针。
// 如果指针为 nil，返回 0。
func DerefInt64(p *int64) int64 {
	if p != nil {
		return *p
	}
	return 0
}

// FormatTimePtr 安全地格式化一个时间指针。
// 如果指针为 nil，返回一个空字符串。
func FormatTimePtr(t *time.Time, layout string) string {
	if t != nil {
		return t.Format(layout)
	}
	return ""
}
