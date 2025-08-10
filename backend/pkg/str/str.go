package str

import "strings"

// CommaSeparatedToSlice はカンマ区切り文字列を []string に変換します。
func CommaSeparatedToSlice(str string) []string {
	parts := strings.Split(str, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
