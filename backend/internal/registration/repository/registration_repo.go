package repository

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"his-go/internal/registration/model"
	"his-go/pkg/errors"
	"his-go/pkg/redis"
)

// RegistrationRepository 挂号数据仓库
type RegistrationRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewRegistrationRepository 创建挂号数据仓库
func NewRegistrationRepository(db *gorm.DB, rdb *redis.Client) *RegistrationRepository {
	return &RegistrationRepository{db: db, rdb: rdb}
}

// FindByID 根据ID查询挂号记录
func (r *RegistrationRepository) FindByID(id string) (*model.Registration, error) {
	var reg model.Registration
	if err := r.db.Where("id = ?", id).First(&reg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "挂号记录不存在")
		}
		return nil, fmt.Errorf("查询挂号记录失败: %w", err)
	}
	return &reg, nil
}

// ListByPatient 分页查询患者挂号记录
func (r *RegistrationRepository) ListByPatient(patientID string, page, pageSize int) ([]model.Registration, int64, error) {
	var list []model.Registration
	var total int64

	query := r.db.Model(&model.Registration{}).Where("patient_id = ?", patientID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计挂号记录失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("查询挂号记录列表失败: %w", err)
	}

	return list, total, nil
}

// ListAll 分页查询全部挂号记录（管理端用）
func (r *RegistrationRepository) ListAll(page, pageSize int, status *int, date string) ([]model.Registration, int64, error) {
	var list []model.Registration
	var total int64

	query := r.db.Model(&model.Registration{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if date != "" {
		query = query.Where("registration_date = ?", date)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计挂号记录失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("查询挂号记录列表失败: %w", err)
	}

	return list, total, nil
}

// Create 创建挂号记录（事务：锁定号源 → 扣减剩余号数 → 插入挂号记录）
func (r *RegistrationRepository) Create(reg *model.Registration) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 悲观锁查询号源
		var schedule model.Schedule
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", reg.ScheduleID).First(&schedule).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.NewAppError(errors.CodeNotFound, "号源不存在")
			}
			return fmt.Errorf("锁定号源失败: %w", err)
		}

		// 2. 检查号源是否已满
		if schedule.RemainCount <= 0 {
			return errors.NewAppError(errors.CodeScheduleFull, "号源已满")
		}

		reg.RegistrationDate = schedule.Date

		// 3. 扣减剩余号源数（乐观锁 version 自动递增）
		schedule.RemainCount--
		result := tx.Select("remain_count", "version").Updates(&schedule)
		if result.Error != nil {
			return fmt.Errorf("扣减号源失败: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return errors.NewAppError(errors.CodeConflict, "号源信息已变更，请重新尝试")
		}

		// 4. 插入挂号记录
		if err := tx.Create(reg).Error; err != nil {
			return fmt.Errorf("创建挂号记录失败: %w", err)
		}

		return nil
	})
}

// Cancel 取消挂号：将状态改为已取消，同时恢复号源数
func (r *RegistrationRepository) Cancel(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var reg model.Registration
		if err := tx.Where("id = ?", id).First(&reg).Error; err != nil {
			return errors.NewAppError(errors.CodeNotFound, "挂号记录不存在")
		}

		if reg.Status == 3 {
			return errors.NewAppError(errors.CodeConflict, "挂号已取消，不可重复操作")
		}

		// 更新挂号状态为已取消
		reg.Status = 3
		if err := tx.Select("status", "updated_at").Updates(&reg).Error; err != nil {
			return fmt.Errorf("更新挂号状态失败: %w", err)
		}

		// 恢复号源数
		if err := tx.Model(&model.Schedule{}).Where("id = ?", reg.ScheduleID).
			UpdateColumn("remain_count", gorm.Expr("remain_count + 1")).Error; err != nil {
			return fmt.Errorf("恢复号源数失败: %w", err)
		}

		return nil
	})
}

// SignIn 签到：将状态改为已签到
func (r *RegistrationRepository) SignIn(id string) error {
	var reg model.Registration
	if err := r.db.Where("id = ?", id).First(&reg).Error; err != nil {
		return errors.NewAppError(errors.CodeNotFound, "挂号记录不存在")
	}

	if reg.Status != 0 {
		return errors.NewAppError(errors.CodeConflict, "当前状态不可签到")
	}

	reg.Status = 1
	if err := r.db.Select("status", "updated_at").Updates(&reg).Error; err != nil {
		return fmt.Errorf("签到失败: %w", err)
	}

	return nil
}

// CountBySchedule 统计某个号源的挂号数量
func (r *RegistrationRepository) CountBySchedule(scheduleID string) int64 {
	var count int64
	r.db.Model(&model.Registration{}).Where("schedule_id = ?", scheduleID).Count(&count)
	return count
}

// ListSchedules 查询某科室某天的号源
func (r *RegistrationRepository) ListSchedules(deptID, date string) ([]model.Schedule, error) {
	var schedules []model.Schedule
	query := r.db.Model(&model.Schedule{})
	if deptID != "" {
		query = query.Where("dept_id = ?", deptID)
	}
	if date != "" {
		query = query.Where("date = ?", date)
	}
	if err := query.Order("time_slot ASC").Find(&schedules).Error; err != nil {
		return nil, fmt.Errorf("查询号源列表失败: %w", err)
	}
	return schedules, nil
}

// PushToQueue 将挂号记录加入 Redis 排队队列（Sorted Set）
func (r *RegistrationRepository) PushToQueue(ctx context.Context, scheduleID, registrationID string, queueNumber int) error {
	key := "queue:" + scheduleID
	return r.rdb.ZAdd(ctx, key, goredis.Z{Score: float64(queueNumber), Member: registrationID})
}

// PopFromQueue 从排队队列中弹出下一个排队号
func (r *RegistrationRepository) PopFromQueue(ctx context.Context, scheduleID string) ([]string, error) {
	key := "queue:" + scheduleID
	zList, err := r.rdb.ZPopMin(ctx, key, 1).Result()
	if err != nil {
		return nil, err
	}
	members := make([]string, 0, len(zList))
	for _, z := range zList {
		members = append(members, z.Member.(string))
	}
	return members, nil
}

// GetQueueRank 查询某挂号记录在排队队列中的位置
func (r *RegistrationRepository) GetQueueRank(ctx context.Context, scheduleID, registrationID string) (int64, error) {
	key := "queue:" + scheduleID
	return r.rdb.ZRank(ctx, key, registrationID).Result()
}
