package service

import (
	"context"
	"fmt"
	"time"

	"his-go/internal/registration/model"
	"his-go/internal/registration/repository"
	"his-go/pkg/common/snowflake"
	"his-go/pkg/errors"
	"his-go/pkg/redis"
)

// RegistrationService 挂号业务服务
type RegistrationService struct {
	repo *repository.RegistrationRepository
	rdb  *redis.Client
	sf   *snowflake.Snowflake
}

// NewRegistrationService 创建挂号业务服务
func NewRegistrationService(repo *repository.RegistrationRepository, rdb *redis.Client, sf *snowflake.Snowflake) *RegistrationService {
	return &RegistrationService{repo: repo, rdb: rdb, sf: sf}
}

// GetByID 根据ID获取挂号记录
func (s *RegistrationService) GetByID(id string) (*model.Registration, error) {
	return s.repo.FindByID(id)
}

// ListByPatient 分页查询患者挂号记录
func (s *RegistrationService) ListByPatient(patientID string, page, pageSize int) ([]model.Registration, int64, error) {
	return s.repo.ListByPatient(patientID, page, pageSize)
}

// ListAll 分页查询全部挂号记录（管理端用）
func (s *RegistrationService) ListAll(page, pageSize int, status *int, date string) ([]model.Registration, int64, error) {
	return s.repo.ListAll(page, pageSize, status, date)
}

// Register 挂号：使用 Redis 分布式锁防止并发超卖
func (s *RegistrationService) Register(patientID, patientName, scheduleID, date string) (*model.Registration, error) {
	ctx := context.Background()
	lockKey := "lock:registration:" + scheduleID

	// 获取分布式锁，TTL 10 秒
	lockValue, lockErr := s.rdb.Lock(ctx, lockKey, 10*time.Second)
	if lockErr != nil {
		return nil, fmt.Errorf("获取分布式锁失败: %w", lockErr)
	}
	if lockValue == "" {
		return nil, errors.NewAppError(errors.CodeConflict, "操作过于频繁，请稍后重试")
	}
	defer func() {
		_ = s.rdb.Unlock(ctx, lockKey, lockValue)
	}()

	// 生成排队号
	queueNumber := int(s.repo.CountBySchedule(scheduleID)) + 1

	// 生成ID
	id, err := s.sf.NextString()
	if err != nil {
		return nil, fmt.Errorf("生成ID失败: %w", err)
	}

	reg := &model.Registration{
		ID:               id,
		PatientID:        patientID,
		PatientName:      patientName,
		ScheduleID:       scheduleID,
		RegistrationDate: date,
		QueueNumber:      queueNumber,
		Status:           0,
	}

	// 执行挂号创建
	if err := s.repo.Create(reg); err != nil {
		return nil, err
	}

	// 挂号成功后加入排队队列
	_ = s.repo.PushToQueue(ctx, scheduleID, reg.ID, queueNumber)

	return reg, nil
}

// Cancel 取消挂号
func (s *RegistrationService) Cancel(id string) error {
	ctx := context.Background()

	reg, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.Cancel(id); err != nil {
		return err
	}

	// 从排队队列中移除
	queueKey := "queue:" + reg.ScheduleID
	_ = s.rdb.ZRem(ctx, queueKey, id)

	return nil
}

// SignIn 签到
func (s *RegistrationService) SignIn(id string) error {
	return s.repo.SignIn(id)
}

// GetQueueStatus 查询排队状态（返回排名，0 表示已就位）
func (s *RegistrationService) GetQueueStatus(registrationID string) (int64, error) {
	ctx := context.Background()

	reg, err := s.repo.FindByID(registrationID)
	if err != nil {
		return 0, err
	}

	rank, err := s.repo.GetQueueRank(ctx, reg.ScheduleID, registrationID)
	if err != nil {
		return 0, fmt.Errorf("查询排队状态失败: %w", err)
	}

	return rank, nil
}

// CallNext 从排队队列中弹出下一个
func (s *RegistrationService) CallNext(scheduleID string) (string, error) {
	ctx := context.Background()

	members, err := s.repo.PopFromQueue(ctx, scheduleID)
	if err != nil {
		return "", fmt.Errorf("叫号失败: %w", err)
	}
	if len(members) == 0 {
		return "", errors.NewAppError(errors.CodeNotFound, "当前无排队患者")
	}

	return members[0], nil
}

// ListSchedules 查询某科室某天的号源
func (s *RegistrationService) ListSchedules(deptID, date string) ([]model.Schedule, error) {
	return s.repo.ListSchedules(deptID, date)
}
