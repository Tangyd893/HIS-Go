package service

import (
	"his-go/internal/inpatient/model"
	"his-go/internal/inpatient/repository"
	"his-go/pkg/errors"
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
	if record == nil {
		return errors.NewAppError(errors.CodeParamInvalid, "住院记录不能为空")
	}
	if record.PatientID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "患者ID不能为空")
	}
	if record.PatientName == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "患者姓名不能为空")
	}
	if record.DeptID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "科室ID不能为空")
	}
	if record.Deposit < 0 {
		return errors.NewAppError(errors.CodeParamInvalid, "押金不能为负数")
	}

	records, err := s.repo.FindByPatientID(record.PatientID)
	if err != nil {
		return err
	}
	for _, r := range records {
		if r.Status == 1 {
			return errors.NewAppError(errors.CodeConflict, "该患者已有住院中记录，不允许重复入院")
		}
	}

	return s.repo.Create(record)
}

// DischargePatient 出院
func (s *InpatientService) DischargePatient(id string) error {
	if id == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "住院记录ID不能为空")
	}

	record, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if record.Status != 1 {
		return errors.NewAppError(errors.CodeConflict, "当前状态不允许出院")
	}

	return s.repo.Discharge(id)
}

// GetInpatient 获取住院记录
func (s *InpatientService) GetInpatient(id string) (*model.InpatientRecord, error) {
	return s.repo.FindByID(id)
}

// ListInpatients 分页查询住院列表
func (s *InpatientService) ListInpatients(deptID string, status int, page, pageSize int) ([]model.InpatientRecord, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.List(deptID, status, page, pageSize)
}

// CreateMedicalOrder 创建医嘱
func (s *InpatientService) CreateMedicalOrder(order *model.MedicalOrder) error {
	if order == nil {
		return errors.NewAppError(errors.CodeParamInvalid, "医嘱不能为空")
	}
	if order.InpatientID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "住院记录ID不能为空")
	}
	if order.OrderType != 1 && order.OrderType != 2 {
		return errors.NewAppError(errors.CodeParamInvalid, "医嘱类型无效，必须为1(长期)或2(临时)")
	}
	if order.Content == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "医嘱内容不能为空")
	}

	record, err := s.repo.FindByID(order.InpatientID)
	if err != nil {
		return err
	}
	if record.Status != 1 {
		return errors.NewAppError(errors.CodeConflict, "住院记录状态不允许创建医嘱")
	}

	return s.repo.CreateOrder(order)
}

// ListMedicalOrders 查询医嘱列表
func (s *InpatientService) ListMedicalOrders(inpatientID string) ([]model.MedicalOrder, error) {
	if inpatientID == "" {
		return nil, errors.NewAppError(errors.CodeParamInvalid, "住院记录ID不能为空")
	}
	return s.repo.ListOrders(inpatientID)
}

// CreateNursingRecord 创建护理记录
func (s *InpatientService) CreateNursingRecord(record *model.NursingRecord) error {
	if record == nil {
		return errors.NewAppError(errors.CodeParamInvalid, "护理记录不能为空")
	}
	if record.InpatientID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "住院记录ID不能为空")
	}
	if record.NurseID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "护士ID不能为空")
	}
	if record.Content == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "护理内容不能为空")
	}

	inpatientRecord, err := s.repo.FindByID(record.InpatientID)
	if err != nil {
		return err
	}
	if inpatientRecord.Status != 1 {
		return errors.NewAppError(errors.CodeConflict, "住院记录状态不允许创建护理记录")
	}

	return s.repo.CreateNursingRecord(record)
}

// ListNursingRecords 查询护理记录列表
func (s *InpatientService) ListNursingRecords(inpatientID string) ([]model.NursingRecord, error) {
	if inpatientID == "" {
		return nil, errors.NewAppError(errors.CodeParamInvalid, "住院记录ID不能为空")
	}
	return s.repo.ListNursingRecords(inpatientID)
}
