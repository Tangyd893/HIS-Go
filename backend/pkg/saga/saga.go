// Package saga 分布式事务 SAGA 协调器框架
// 基于 Saga 模式的补偿式分布式事务管理
// 适用于跨微服务的长时间运行事务（如挂号→接诊→处方→收费→发药）
//
// 使用示例:
//
//	saga := saga.NewSAGA("挂号流程")
//	saga.AddStep("扣减号源", registerAction, registerCompensate)
//	saga.AddStep("创建挂号记录", createRegAction, cancelRegCompensate)
//	saga.AddStep("发送通知", notifyAction, notifyCompensate)
//	err := saga.Execute(ctx)
package saga

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Action SAGA 步骤的正向操作函数
// 返回 error 表示该步骤失败，触发补偿
type Action func(ctx context.Context, data *SagaData) error

// Compensate SAGA 步骤的补偿操作函数
// 当后续步骤失败时，按逆序执行补偿
type Compensate func(ctx context.Context, data *SagaData) error

// Step SAGA 事务步骤
type Step struct {
	Name       string     `json:"name"`
	Action     Action     `json:"-"`
	Compensate Compensate `json:"-"`
	Status     StepStatus `json:"status"`
	Error      string     `json:"error,omitempty"`
	StartedAt  time.Time  `json:"startedAt"`
	EndedAt    time.Time  `json:"endedAt"`
}

// StepStatus 步骤状态
type StepStatus int

const (
	StepPending    StepStatus = 0
	StepRunning    StepStatus = 1
	StepSuccess    StepStatus = 2
	StepFailed     StepStatus = 3
	StepCompensate StepStatus = 4
)

// SagaStatus 事务状态
type SagaStatus int

const (
	SagaPending    SagaStatus = 0
	SagaRunning    SagaStatus = 1
	SagaSuccess    SagaStatus = 2
	SagaFailed     SagaStatus = 3
	SagaCompensing SagaStatus = 4
	SagaCompensed  SagaStatus = 5
)

// SagaData 事务上下文数据，可在步骤间传递
type SagaData struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func NewSagaData() *SagaData {
	return &SagaData{data: make(map[string]interface{})}
}

func (d *SagaData) Set(key string, value interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data[key] = value
}

func (d *SagaData) Get(key string) (interface{}, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	v, ok := d.data[key]
	return v, ok
}

func (d *SagaData) GetString(key string) string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if v, ok := d.data[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Saga SAGA 分布式事务协调器
type Saga struct {
	name      string
	steps     []*Step
	status    SagaStatus
	data      *SagaData
	startedAt time.Time
	endedAt   time.Time
}

// NewSAGA 创建 SAGA 事务
func NewSAGA(name string) *Saga {
	return &Saga{
		name:   name,
		steps:  make([]*Step, 0),
		status: SagaPending,
		data:   NewSagaData(),
	}
}

// AddStep 添加事务步骤
func (s *Saga) AddStep(name string, action Action, compensate Compensate) *Saga {
	s.steps = append(s.steps, &Step{
		Name:       name,
		Action:     action,
		Compensate: compensate,
		Status:     StepPending,
	})
	return s
}

// SetData 设置事务上下文数据
func (s *Saga) SetData(key string, value interface{}) *Saga {
	s.data.Set(key, value)
	return s
}

// Execute 执行 SAGA 事务
// 正向执行所有步骤，若任一失败则逆序执行补偿
func (s *Saga) Execute(ctx context.Context) error {
	s.status = SagaRunning
	s.startedAt = time.Now()

	// 正向执行
	for i, step := range s.steps {
		select {
		case <-ctx.Done():
			s.compensate(ctx, i-1)
			return ctx.Err()
		default:
		}

		step.Status = StepRunning
		step.StartedAt = time.Now()

		if err := step.Action(ctx, s.data); err != nil {
			step.Status = StepFailed
			step.Error = err.Error()
			step.EndedAt = time.Now()

			// 从当前步骤的前一步开始逆序补偿
			s.compensate(ctx, i-1)
			s.status = SagaCompensed
			s.endedAt = time.Now()
			return fmt.Errorf("SAGA [%s] 步骤 [%s] 失败: %w", s.name, step.Name, err)
		}

		step.Status = StepSuccess
		step.EndedAt = time.Now()
	}

	s.status = SagaSuccess
	s.endedAt = time.Now()
	return nil
}

// compensate 执行补偿操作（从指定步骤逆序执行到第0步）
func (s *Saga) compensate(ctx context.Context, fromIndex int) {
	s.status = SagaCompensing

	for i := fromIndex; i >= 0; i-- {
		step := s.steps[i]
		if step.Status != StepSuccess {
			continue
		}

		compCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := step.Compensate(compCtx, s.data); err != nil {
			// 补偿失败记录日志，但继续执行剩余补偿
			step.Error = "补偿失败: " + err.Error()
		}

		step.Status = StepCompensate
	}
}

// Steps 返回步骤列表（用于日志/监控）
func (s *Saga) Steps() []*Step { return s.steps }

// Duration 事务耗时
func (s *Saga) Duration() time.Duration {
	if s.endedAt.IsZero() {
		return time.Since(s.startedAt)
	}
	return s.endedAt.Sub(s.startedAt)
}

// ─── 预定义业务 SAGA 模板 ───

// RegistrationSaga 挂号流程 SAGA
// 扣减号源 → 创建挂号记录 → 发送通知
func RegistrationSaga(deductAction, createRegAction, notifyAction Action,
	deductComp, cancelRegComp, notifyComp Compensate) *Saga {
	return NewSAGA("挂号流程").
		AddStep("扣减号源", deductAction, deductComp).
		AddStep("创建挂号记录", createRegAction, cancelRegComp).
		AddStep("发送通知", notifyAction, notifyComp)
}

// PrescriptionSaga 处方流程 SAGA
// 审核处方 → 扣减库存 → 更新处方状态
func PrescriptionSaga(reviewAction, stockAction, updateAction Action,
	reviewComp, stockComp, updateComp Compensate) *Saga {
	return NewSAGA("处方流程").
		AddStep("审核处方", reviewAction, reviewComp).
		AddStep("扣减库存", stockAction, stockComp).
		AddStep("更新处方状态", updateAction, updateComp)
}

// DischargeSaga 出院结算 SAGA
// 计算费用 → 执行支付 → 更新床位 → 更新住院状态
func DischargeSaga(calcAction, payAction, bedAction, statusAction Action,
	calcComp, payComp, bedComp, statusComp Compensate) *Saga {
	return NewSAGA("出院结算").
		AddStep("计算费用", calcAction, calcComp).
		AddStep("执行支付", payAction, payComp).
		AddStep("更新床位", bedAction, bedComp).
		AddStep("更新住院状态", statusAction, statusComp)
}
