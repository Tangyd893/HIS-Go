package errors

import (
	"testing"
)

func TestGetMessage(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{CodeSuccess, "操作成功"},
		{CodeUnknownError, "未知错误"},
		{CodeParamInvalid, "参数无效"},
		{CodeUnauthorized, "未授权"},
		{CodeForbidden, "无权限"},
		{CodeNotFound, "未找到"},
		{CodeConflict, "数据冲突"},
		{CodeTimeout, "请求超时"},
		{CodeInternalError, "内部错误"},
		{CodeServiceUnavail, "服务不可用"},
		{CodeRateLimited, "请求过于频繁"},
		{CodeUserNotFound, "用户不存在"},
		{CodeUserPasswordWrong, "密码错误"},
		{CodeScheduleFull, "号源已满"},
		{CodeRegistrationRepeat, "重复挂号"},
		{CodePrescriptionInvalid, "处方不合规"},
		{CodeStockInsufficient, "库存不足"},
		{CodePayFailed, "支付失败"},
		{CodeRefundFailed, "退款失败"},
	}

	for _, tt := range tests {
		got := GetMessage(tt.code)
		if got != tt.want {
			t.Errorf("GetMessage(%d) = %q, 期望 %q", tt.code, got, tt.want)
		}
	}
}

func TestGetMessage_UnknownCode(t *testing.T) {
	got := GetMessage(99999)
	if got != "未知错误" {
		t.Errorf("GetMessage(99999) = %q, 期望 '未知错误'", got)
	}
}

func TestNewAppError(t *testing.T) {
	err := NewAppError(CodeScheduleFull, "号源剩余 0")
	if err.Code != CodeScheduleFull {
		t.Errorf("期望 Code=%d, 实际=%d", CodeScheduleFull, err.Code)
	}
	if err.Message != "号源已满" {
		t.Errorf("期望 Message='号源已满', 实际=%q", err.Message)
	}
	if err.Detail != "号源剩余 0" {
		t.Errorf("期望 Detail='号源剩余 0', 实际=%q", err.Detail)
	}
}

func TestAppError_Error(t *testing.T) {
	err := NewAppError(CodeUserNotFound, "")
	if err.Error() != "用户不存在" {
		t.Errorf("期望 Error()='用户不存在', 实际=%q", err.Error())
	}
}

func TestErrorCodeUniqueness(t *testing.T) {
	codes := map[int]bool{}
	for code := range codeMessages {
		if codes[code] {
			t.Errorf("错误码 %d 重复定义", code)
		}
		codes[code] = true
	}
}
