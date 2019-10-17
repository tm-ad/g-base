package validator

import (
	"strings"
)

// DefaultStringValue 检查字符串是否为空（trim后），不为空则返回原值.trim()，为空则返回默认值
func DefaultStringValue(val, defaultValue string) string {
	val = strings.Trim(val, "")

	if val == "" {
		return defaultValue
	}

	return val
}
