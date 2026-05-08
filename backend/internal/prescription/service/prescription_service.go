package service

import (
	"his-go/internal/prescription/model"
	"his-go/internal/prescription/repository"
	"his-go/pkg/errors"
)

// PrescriptionService 处方业务服务
type PrescriptionService struct {
	repo *repository.PrescriptionRepository
}

// NewPrescriptionService 创建处方业务服务
func NewPrescriptionService(repo *repository.PrescriptionRepository) *PrescriptionService {
	return &PrescriptionService{repo: repo}
}

// Create 创建处方（含明细），初始状态为草稿
func (s *PrescriptionService) Create(p *model.Prescription, details []model.PrescriptionDetail) error {
	p.Status = 0 // 草稿
	return s.repo.Create(p, details)
}

// GetByID 根据ID获取处方
func (s *PrescriptionService) GetByID(id string) (*model.Prescription, error) {
	return s.repo.FindByID(id)
}

// ListByPatient 分页查询患者处方
func (s *PrescriptionService) ListByPatient(patientID string, page, pageSize int) ([]model.Prescription, int64, error) {
	return s.repo.ListByPatient(patientID, page, pageSize)
}

// ListByDoctor 分页查询医生处方
func (s *PrescriptionService) ListByDoctor(doctorID string, page, pageSize int) ([]model.Prescription, int64, error) {
	return s.repo.ListByDoctor(doctorID, page, pageSize)
}

// SubmitForReview 提交审核（草稿 → 待审核）
func (s *PrescriptionService) SubmitForReview(id string) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if p.Status != 0 {
		return errors.NewAppError(errors.CodePrescriptionInvalid, "只有草稿状态可提交审核")
	}
	return s.repo.UpdateStatus(id, 1)
}

// Review 审核处方（待审核 → 已审核 / 退回草稿）
func (s *PrescriptionService) Review(id string, approved bool, comment string) error {
	return s.repo.Review(id, approved, comment)
}

// Cancel 取消处方（将状态改为失效）
func (s *PrescriptionService) Cancel(id string) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if p.Status >= 3 {
		return errors.NewAppError(errors.CodeConflict, "已收费或已发药的处方不可取消")
	}
	return s.repo.UpdateStatus(id, 3)
}
