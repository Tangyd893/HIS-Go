package service

import (
	"his-go/internal/examination/model"
	"his-go/internal/examination/repository"
	"his-go/pkg/errors"
)

// ExaminationService 检查检验业务服务
type ExaminationService struct {
	repo *repository.ExaminationRepository
}

// NewExaminationService 创建检查检验业务服务
func NewExaminationService(repo *repository.ExaminationRepository) *ExaminationService {
	return &ExaminationService{repo: repo}
}

// CreateReport 创建检查报告
func (s *ExaminationService) CreateReport(report *model.ExaminationReport) error {
	report.Status = 0 // 待检查
	return s.repo.CreateReport(report)
}

// GetByID 根据ID获取报告
func (s *ExaminationService) GetByID(id string) (*model.ExaminationReport, error) {
	return s.repo.FindByID(id)
}

// ListByPatient 分页查询患者检查报告
func (s *ExaminationService) ListByPatient(patientID string, status int, page, pageSize int) ([]model.ExaminationReport, int64, error) {
	return s.repo.ListByPatient(patientID, status, page, pageSize)
}

// UpdateReport 更新报告（填写检查所见和印象）
func (s *ExaminationService) UpdateReport(report *model.ExaminationReport) error {
	existing, err := s.repo.FindByID(report.ID)
	if err != nil {
		return err
	}
	if existing.Status != 0 {
		return errors.NewAppError(errors.CodeConflict, "只有待检查状态的报告可修改")
	}
	report.Status = 1 // 已检查
	return s.repo.UpdateReport(report)
}

// Review 审核报告
func (s *ExaminationService) Review(reportID, reviewerID string, approved bool, comment string) error {
	return s.repo.Review(reportID, reviewerID, approved, comment)
}
