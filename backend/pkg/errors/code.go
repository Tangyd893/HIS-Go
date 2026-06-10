// Package errors 统一错误码定义
package errors

import "fmt"

// 常用校验错误消息（避免字符串重复 go:S1192）
const (
	MsgPatientIDRequired = "患者ID不能为空"
	MsgDoctorIDRequired  = "医生ID不能为空"
	MsgDeptIDRequired    = "科室ID不能为空"
	MsgNameRequired      = "名称不能为空"
)

// 错误码定义
const (
	CodeSuccess        = 0
	CodeUnknownError   = 10000
	CodeParamInvalid   = 10001
	CodeUnauthorized   = 10002
	CodeForbidden      = 10003
	CodeNotFound       = 10004
	CodeConflict       = 10005
	CodeTimeout        = 10006
	CodeInternalError  = 10007
	CodeServiceUnavail = 10008
	CodeRateLimited    = 10009

	// 业务错误码范围 20000-29999
	CodeUserNotFound        = 20001
	CodeUserPasswordWrong   = 20002
	CodeScheduleFull        = 20003
	CodeRegistrationRepeat  = 20004
	CodePrescriptionInvalid = 20005
	CodeStockInsufficient   = 20006
	CodePayFailed           = 20007
	CodeRefundFailed        = 20008
	CodeAuthExpired         = 20009 // 授权已过期
	CodeAuthDuplicate       = 20010 // 重复授权
	CodeAuthNotFound        = 20011 // 授权不存在
)

var codeMessages = map[int]string{
	CodeSuccess:        "操作成功",
	CodeUnknownError:   "未知错误",
	CodeParamInvalid:   "参数无效",
	CodeUnauthorized:   "未授权",
	CodeForbidden:      "无权限",
	CodeNotFound:       "未找到",
	CodeConflict:       "数据冲突",
	CodeTimeout:        "请求超时",
	CodeInternalError:  "内部错误",
	CodeServiceUnavail: "服务不可用",
	CodeRateLimited:    "请求过于频繁",

	CodeUserNotFound:        "用户不存在",
	CodeUserPasswordWrong:   "密码错误",
	CodeScheduleFull:        "号源已满",
	CodeRegistrationRepeat:  "重复挂号",
	CodePrescriptionInvalid: "处方不合规",
	CodeStockInsufficient:   "库存不足",
	CodePayFailed:           "支付失败",
	CodeRefundFailed:        "退款失败",
	CodeAuthExpired:         "授权已过期",
	CodeAuthDuplicate:       "重复授权",
	CodeAuthNotFound:        "授权不存在",
}

// WrapQueryError 统一查询错误包装（消除 ~28 处字符串重复 go:S1192）
func WrapQueryError(entity string, err error) error {
	return fmt.Errorf("查询%s失败: %w", entity, err)
}

// WrapCreateError 统一创建错误包装
func WrapCreateError(entity string, err error) error {
	return fmt.Errorf("创建%s失败: %w", entity, err)
}

// WrapUpdateError 统一更新错误包装
func WrapUpdateError(entity string, err error) error {
	return fmt.Errorf("更新%s失败: %w", entity, err)
}

// WrapCountError 统一统计查询错误包装
func WrapCountError(entity string, err error) error {
	return fmt.Errorf("统计%s失败: %w", entity, err)
}

// GetMessage 获取错误码对应消息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
