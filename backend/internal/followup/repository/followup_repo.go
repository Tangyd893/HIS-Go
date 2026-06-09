package repository

import (
	"fmt"

	"gorm.io/gorm"

	"his-go/internal/followup/model"
	"his-go/pkg/errors"
)

// FollowupRepository 随访数据仓库
type FollowupRepository struct {
	db *gorm.DB
}

// NewFollowupRepository 创建随访数据仓库
func NewFollowupRepository(db *gorm.DB) *FollowupRepository {
	return &FollowupRepository{db: db}
}

// CreatePlan 创建随访计划并生成任务
func (r *FollowupRepository) CreatePlan(plan *model.FollowupPlan, tasks []model.FollowupTask) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(plan).Error; err != nil {
			return fmt.Errorf("创建随访计划失败: %w", err)
		}

		for i := range tasks {
			tasks[i].PlanID = plan.ID
		}
		if len(tasks) > 0 {
			if err := tx.Create(&tasks).Error; err != nil {
				return fmt.Errorf("创建随访任务失败: %w", err)
			}
		}

		return nil
	})
}

// FindPlanByID 根据ID查询随访计划
func (r *FollowupRepository) FindPlanByID(id string) (*model.FollowupPlan, error) {
	var plan model.FollowupPlan
	if err := r.db.Where("id = ?", id).First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "随访计划不存在")
		}
		return nil, fmt.Errorf("查询随访计划失败: %w", err)
	}
	return &plan, nil
}

// ListPlans 分页查询随访计划
func (r *FollowupRepository) ListPlans(patientID string, status int, page, pageSize int) ([]model.FollowupPlan, int64, error) {
	var list []model.FollowupPlan
	var total int64

	query := r.db.Model(&model.FollowupPlan{})
	if patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计随访计划失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("查询随访计划列表失败: %w", err)
	}

	return list, total, nil
}

// ExecuteTask 执行随访任务
func (r *FollowupRepository) ExecuteTask(taskID string, result string) error {
	updates := map[string]interface{}{
		"status":  int8(1), // 已完成
		"content": result,
	}
	resultDB := r.db.Model(&model.FollowupTask{}).Where("id = ?", taskID).Updates(updates)
	if resultDB.Error != nil {
		return fmt.Errorf("执行随访任务失败: %w", resultDB.Error)
	}
	if resultDB.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "随访任务不存在")
	}
	return nil
}

// SubmitSurvey 提交满意度调查
func (r *FollowupRepository) SubmitSurvey(survey *model.SatisfactionSurvey) error {
	if err := r.db.Create(survey).Error; err != nil {
		return fmt.Errorf("提交满意度调查失败: %w", err)
	}
	return nil
}

// GetTasksByPlanID 获取随访计划的所有任务
func (r *FollowupRepository) GetTasksByPlanID(planID string) ([]model.FollowupTask, error) {
	var list []model.FollowupTask
	if err := r.db.Where("plan_id = ?", planID).Order("execute_date ASC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("查询随访任务失败: %w", err)
	}
	return list, nil
}
