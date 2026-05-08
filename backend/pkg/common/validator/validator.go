// Package validator 参数校验工具
package validator

import (
	"regexp"
	"strings"
)

var (
	phoneRegex  = regexp.MustCompile(`^1[3-9]\d{9}$`)
	idCardRegex = regexp.MustCompile(`^\d{17}[\dXx]$`)
	emailRegex  = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// IsPhone 校验手机号
func IsPhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

// IsIdCard 校验身份证号
func IsIdCard(idCard string) bool {
	return idCardRegex.MatchString(idCard)
}

// IsEmail 校验邮箱
func IsEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// SanitizeString 过滤 XSS 敏感字符
func SanitizeString(s string) string {
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// TrimSpace 去除首尾空格
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// IsEmpty 判断字符串是否为空
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
