package service

import (
	"his-go/internal/inpatient/model"
	"his-go/internal/inpatient/repository"
)

// InpatientService 住院管理业务服务
type InpatientService struct {
	repo *repository.InpatientRepository
}

// NewInpatientService 创建住院管理业务服务
func NewInpatientService(repo *repository.InpatientRepository) *InpatientService {
	return &InpatientService{repo: repo}
}

// AdmitPatient 入院登记
func (s *InpatientService) AdmitPatient(record *model.InpatientRecord) error {
	return s.repo.Create(record)
}

// DischargePatient 出院
func (s *InpatientService) DischargePatient(id string) error {
	return s.repo.Discharge(id)
}

// GetInpatient 获取住院记录
func (s *InpatientService) GetInpatient(id string) (*model.InpatientRecord, error) {
	return s.repo.FindByID(id)
}

// ListInpatients 分页查询住院列表
func (s *InpatientService) ListInpatients(deptID string, status int, page, pageSize int) ([]model.InpatientRecord, int64, error) {
	return s.repo.List(deptID, status, page, pageSize)
}

// CreateMedicalOrder 创建医嘱
func (s *InpatientService) CreateMedicalOrder(order *model.MedicalOrder) error {
	return s.repo.CreateOrder(order)
}

// ListMedicalOrders 查询医嘱列表
func (s *InpatientService) ListMedicalOrders(inpatientID string) ([]model.MedicalOrder, error) {
	return s.repo.ListOrders(inpatientID)
}

// CreateNursingRecord 创建护理记录
func (s *InpatientService) CreateNursingRecord(record *model.NursingRecord) error {
	return s.repo.CreateNursingRecord(record)
}

// ListNursingRecords 查询护理记录列表
func (s *InpatientService) ListNursingRecords(inpatientID string) ([]model.NursingRecord, error) {
	return s.repo.ListNursingRecords(inpatientID)
}
