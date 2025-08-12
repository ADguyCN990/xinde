package util

import "time"

// StringToPointer (或者更简洁的 Ptr, String) 返回一个指向字符串值的指针
func StringToPointer(s string) *string {
	return &s
}

// TimeToPointer 返回一个指向 time.Time 值的指针
func TimeToPointer(t time.Time) *time.Time {
	return &t
}

// UintToPointer 返回一个指向 uint 值的指针
func UintToPointer(u uint) *uint {
	return &u
}
