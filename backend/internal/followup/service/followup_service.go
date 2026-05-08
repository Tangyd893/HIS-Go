package service

import (
	"time"

	"his-go/internal/followup/model"
	"his-go/internal/followup/repository"
	"his-go/pkg/errors"
)

// FollowupService 随访业务服务
type FollowupService struct {
	repo *repository.FollowupRepository
}

// NewFollowupService 创建随访业务服务
func NewFollowupService(repo *repository.FollowupRepository) *FollowupService {
	return &FollowupService{repo: repo}
}

// CreatePlan 创建随访计划并根据频率生成任务
func (s *FollowupService) CreatePlan(plan *model.FollowupPlan) error {
	tasks, err := s.GenerateTasks(plan)
	if err != nil {
		return err
	}
	plan.Status = 1 // 进行中
	return s.repo.CreatePlan(plan, tasks)
}

// GetPlan 获取随访计划
func (s *FollowupService) GetPlan(id string) (*model.FollowupPlan, error) {
	return s.repo.FindPlanByID(id)
}

// ListPlans 分页查询随访计划
func (s *FollowupService) ListPlans(patientID string, status int, page, pageSize int) ([]model.FollowupPlan, int64, error) {
	return s.repo.ListPlans(patientID, status, page, pageSize)
}

// ExecuteTask 执行随访任务
func (s *FollowupService) ExecuteTask(taskID string, result string) error {
	return s.repo.ExecuteTask(taskID, result)
}

// SubmitSurvey 提交满意度调查
func (s *FollowupService) SubmitSurvey(survey *model.SatisfactionSurvey) error {
	if survey.Score < 1 || survey.Score > 5 {
		return errors.NewAppError(errors.CodeParamInvalid, "评分范围为1-5")
	}
	return s.repo.SubmitSurvey(survey)
}

// GenerateTasks 根据频率从StartDate到EndDate生成任务
func (s *FollowupService) GenerateTasks(plan *model.FollowupPlan) ([]model.FollowupTask, error) {
	start, err := time.Parse("2006-01-02", plan.StartDate)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeParamInvalid, "开始日期格式错误")
	}
	end, err := time.Parse("2006-01-02", plan.EndDate)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeParamInvalid, "结束日期格式错误")
	}

	if start.After(end) {
		return nil, errors.NewAppError(errors.CodeParamInvalid, "开始日期不能晚于结束日期")
	}

	interval := plan.Frequency
	if interval < 1 || interval > 3 {
		interval = 1
	}

	var tasks []model.FollowupTask
	for d := start; !d.After(end); {
		tasks = append(tasks, model.FollowupTask{
			ExecuteDate: d.Format("2006-01-02"),
			Type:        1, // 默认电话随访
			Status:      0, // 待执行
			Content:     plan.PlanName,
		})

		switch interval {
		case 1:
			d = d.AddDate(0, 0, 7)
		case 2:
			d = d.AddDate(0, 0, 14)
		case 3:
			d = d.AddDate(0, 1, 0)
		default:
			d = d.AddDate(0, 0, 7)
		}
	}

	return tasks, nil
}
