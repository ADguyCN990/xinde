package util

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

// FormatContentDisposition correctly formats the Content-Disposition header
// to support non-ASCII filenames across all modern clients.
func FormatContentDisposition(filename string) string {
	// 1. 准备一个安全的、只包含 ASCII 的备用文件名
	asciiFilename := "attachment" // 默认备用名
	if ext := filepath.Ext(filename); ext != "" {
		asciiFilename += ext
	}

	// 2. 对原始文件名进行符合 RFC 5987/6266 的百分号编码
	var sb strings.Builder
	for _, r := range filename {
		// 对于 ASCII 范围内的安全字符，直接附加
		if r < 128 && (r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '.' || r == '-' || r == '_') {
			sb.WriteRune(r)
		} else {
			// 对于其他所有字符（包括中文、空格、特殊符号），使用 URL 编码
			sb.WriteString(url.QueryEscape(string(r)))
		}
	}
	encodedFilename := sb.String()

	// 3. 构建完整的、兼容性最好的 Content-Disposition 字符串
	return fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, asciiFilename, encodedFilename)
}
