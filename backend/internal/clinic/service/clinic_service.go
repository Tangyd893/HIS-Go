package service

import (
	"his-go/internal/clinic/model"
	"his-go/internal/clinic/repository"
	"his-go/pkg/demo"
)

// ClinicService 门诊诊疗业务服务
type ClinicService struct {
	repo *repository.ClinicRepository
}

// NewClinicService 创建门诊诊疗业务服务
func NewClinicService(repo *repository.ClinicRepository) *ClinicService {
	return &ClinicService{repo: repo}
}

// CreateRecord 创建门诊诊疗记录
func (s *ClinicService) CreateRecord(record *model.ClinicRecord) error {
	return s.repo.CreateRecord(record)
}

// GetByID 根据ID获取诊疗记录
func (s *ClinicService) GetByID(id string) (*model.ClinicRecord, error) {
	return s.repo.FindByID(id)
}

// ListByPatient 分页查询诊疗记录（patientID 为空时返回全部，管理端列表用）
func (s *ClinicService) ListByPatient(patientID string, page, pageSize int) ([]model.ClinicRecord, int64, error) {
	list, total, err := s.repo.ListByPatient(patientID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		if list[i].DoctorName == "" {
			list[i].DoctorName = demo.DoctorName(list[i].DoctorID)
		}
	}
	return list, total, nil
}

// ListByDoctor 分页查询医生诊疗记录
func (s *ClinicService) ListByDoctor(doctorID string, page, pageSize int) ([]model.ClinicRecord, int64, error) {
	return s.repo.ListByDoctor(doctorID, page, pageSize)
}

// UpdateRecord 更新诊疗记录
func (s *ClinicService) UpdateRecord(record *model.ClinicRecord) error {
	return s.repo.UpdateRecord(record)
}

// CreateExamRequest 创建检查申请单
func (s *ClinicService) CreateExamRequest(request *model.ExaminationRequest) error {
	return s.repo.CreateExamRequest(request)
}

// ListExamRequests 查询患者的检查申请单
func (s *ClinicService) ListExamRequests(patientID string) ([]model.ExaminationRequest, error) {
	return s.repo.ListExamRequests(patientID)
}
