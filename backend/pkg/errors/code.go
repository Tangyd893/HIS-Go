// Package errors 统一错误码定义
package errors

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

// GetMessage 获取错误码对应消息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
