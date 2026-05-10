// Package payment 第三方支付接口定义与桩实现
// 提供微信支付、支付宝、医保结算的抽象接口和测试桩
package payment

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ─── 公共类型定义 ───

// PayMethod 支付方式
type PayMethod int

const (
	PayCash   PayMethod = 0
	PayWechat PayMethod = 1
	PayAlipay PayMethod = 2
	PayBank   PayMethod = 3
	PayMedIns PayMethod = 4 // 医保
)

func (m PayMethod) String() string {
	switch m {
	case PayCash:
		return "现金"
	case PayWechat:
		return "微信支付"
	case PayAlipay:
		return "支付宝"
	case PayBank:
		return "银行卡"
	case PayMedIns:
		return "医保结算"
	default:
		return fmt.Sprintf("未知(%d)", m)
	}
}

// PayStatus 支付状态
type PayStatus int

const (
	StatusPending PayStatus = 0
	StatusSuccess PayStatus = 1
	StatusFailed  PayStatus = 2
	StatusRefund  PayStatus = 3
	StatusClosed  PayStatus = 4
)

// PayOrder 支付订单
type PayOrder struct {
	OrderNo     string    `json:"orderNo"`
	Amount      int64     `json:"amount"` // 金额（分）
	Method      PayMethod `json:"method"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	PatientID   string    `json:"patientId"`
	PatientName string    `json:"patientName"`
	ReturnURL   string    `json:"returnUrl,omitempty"`
	NotifyURL   string    `json:"notifyUrl,omitempty"`
}

// PayResult 支付结果
type PayResult struct {
	Success       bool      `json:"success"`
	TransactionID string    `json:"transactionId"`
	OrderNo       string    `json:"orderNo"`
	Amount        int64     `json:"amount"`
	PaidAt        time.Time `json:"paidAt"`
	ErrorCode     string    `json:"errorCode,omitempty"`
	ErrorMsg      string    `json:"errorMsg,omitempty"`
}

// RefundOrder 退款订单
type RefundOrder struct {
	OriginalOrderNo string `json:"originalOrderNo"`
	RefundNo        string `json:"refundNo"`
	Amount          int64  `json:"amount"`
	Reason          string `json:"reason"`
}

// RefundResult 退款结果
type RefundResult struct {
	Success    bool      `json:"success"`
	RefundNo   string    `json:"refundNo"`
	Amount     int64     `json:"amount"`
	RefundedAt time.Time `json:"refundedAt"`
	ErrorMsg   string    `json:"errorMsg,omitempty"`
}

// ─── 支付接口定义 ───

// PayChannel 支付通道接口
// 所有第三方支付实现此接口
type PayChannel interface {
	// Name 返回支付通道名称
	Name() string

	// Pay 发起支付
	Pay(ctx context.Context, order *PayOrder) (*PayResult, error)

	// Refund 发起退款
	Refund(ctx context.Context, order *RefundOrder) (*RefundResult, error)

	// Query 查询支付状态
	Query(ctx context.Context, orderNo string) (*PayResult, error)

	// Close 关闭支付订单
	Close(ctx context.Context, orderNo string) error

	// Health 检查支付通道连通性
	Health(ctx context.Context) error
}

// PaymentManager 支付管理器（门面模式）
type PaymentManager struct {
	channels map[PayMethod]PayChannel
}

// NewPaymentManager 创建支付管理器
func NewPaymentManager() *PaymentManager {
	return &PaymentManager{
		channels: make(map[PayMethod]PayChannel),
	}
}

// Register 注册支付通道
func (m *PaymentManager) Register(method PayMethod, channel PayChannel) {
	m.channels[method] = channel
}

// Pay 执行支付
func (m *PaymentManager) Pay(ctx context.Context, order *PayOrder) (*PayResult, error) {
	ch, ok := m.channels[order.Method]
	if !ok {
		return nil, fmt.Errorf("不支持的支付方式: %s", order.Method)
	}
	return ch.Pay(ctx, order)
}

// Refund 执行退款
func (m *PaymentManager) Refund(ctx context.Context, order *RefundOrder, method PayMethod) (*RefundResult, error) {
	ch, ok := m.channels[method]
	if !ok {
		return nil, fmt.Errorf("不支持的支付方式: %s", method)
	}
	return ch.Refund(ctx, order)
}

// GetSupportedMethods 返回所有已注册的支付方式
func (m *PaymentManager) GetSupportedMethods() []PayMethod {
	methods := make([]PayMethod, 0, len(m.channels))
	for m := range m.channels {
		methods = append(methods, m)
	}
	return methods
}

// ─── 桩实现 (测试/开发环境) ───

var ErrNotImplemented = errors.New("该支付通道尚未接入生产环境")

// StubPayChannel 通用支付桩实现
// 模拟第三方支付行为，支持 Pay/Refund/Query/Close
type StubPayChannel struct {
	name     string
	simulate bool          // true 模拟成功，false 模拟失败
	delay    time.Duration // 模拟支付延迟
}

// NewStubPayChannel 创建支付桩
func NewStubPayChannel(name string, simulate bool) *StubPayChannel {
	return &StubPayChannel{
		name:     name,
		simulate: simulate,
		delay:    500 * time.Millisecond,
	}
}

func (s *StubPayChannel) Name() string { return s.name }

func (s *StubPayChannel) Pay(ctx context.Context, order *PayOrder) (*PayResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(s.delay):
	}

	if !s.simulate {
		return &PayResult{
			Success:   false,
			OrderNo:   order.OrderNo,
			Amount:    order.Amount,
			ErrorCode: "SIM_FAIL",
			ErrorMsg:  fmt.Sprintf("[桩] %s 支付模拟失败", s.name),
		}, nil
	}

	return &PayResult{
		Success:       true,
		TransactionID: fmt.Sprintf("STUB_%s_%d", s.name, time.Now().UnixNano()),
		OrderNo:       order.OrderNo,
		Amount:        order.Amount,
		PaidAt:        time.Now(),
	}, nil
}

func (s *StubPayChannel) Refund(ctx context.Context, order *RefundOrder) (*RefundResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(s.delay):
	}

	if !s.simulate {
		return &RefundResult{
			Success:  false,
			RefundNo: order.RefundNo,
			Amount:   order.Amount,
			ErrorMsg: fmt.Sprintf("[桩] %s 退款模拟失败", s.name),
		}, nil
	}

	return &RefundResult{
		Success:    true,
		RefundNo:   order.RefundNo,
		Amount:     order.Amount,
		RefundedAt: time.Now(),
	}, nil
}

func (s *StubPayChannel) Query(ctx context.Context, orderNo string) (*PayResult, error) {
	return &PayResult{
		Success:       true,
		TransactionID: "STUB_QUERY_" + orderNo,
		OrderNo:       orderNo,
		PaidAt:        time.Now(),
	}, nil
}

func (s *StubPayChannel) Close(ctx context.Context, orderNo string) error { return nil }

func (s *StubPayChannel) Health(ctx context.Context) error {
	if !s.simulate {
		return errors.New("[桩] 支付通道不可用")
	}
	return nil
}

// WechatPayStub 微信支付桩
type WechatPayStub struct {
	*StubPayChannel
}

func NewWechatPayStub() *WechatPayStub {
	return &WechatPayStub{NewStubPayChannel("微信支付", true)}
}

// AlipayStub 支付宝桩
type AlipayStub struct {
	*StubPayChannel
}

func NewAlipayStub() *AlipayStub {
	return &AlipayStub{NewStubPayChannel("支付宝", true)}
}

// MedInsStub 医保结算桩
type MedInsStub struct {
	*StubPayChannel
}

func NewMedInsStub() *MedInsStub {
	return &MedInsStub{NewStubPayChannel("医保结算", true)}
}

// ─── 默认管理器工厂 ───

// NewDefaultPaymentManager 创建默认支付管理器（注册所有桩通道）
func NewDefaultPaymentManager() *PaymentManager {
	m := NewPaymentManager()
	m.Register(PayCash, NewStubPayChannel("现金", true))
	m.Register(PayWechat, NewWechatPayStub())
	m.Register(PayAlipay, NewAlipayStub())
	m.Register(PayBank, NewStubPayChannel("银行卡", true))
	m.Register(PayMedIns, NewMedInsStub())
	return m
}
