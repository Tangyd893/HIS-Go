package service

import (
	"his-go/internal/emr/model"
	"his-go/internal/emr/repository"
	"his-go/pkg/demo"
	"his-go/pkg/errors"
)

// EMRService 电子病历业务服务
type EMRService struct {
	repo *repository.EMRRepository
}

// NewEMRService 创建电子病历业务服务
func NewEMRService(repo *repository.EMRRepository) *EMRService {
	return &EMRService{repo: repo}
}

// CreateRecord 创建病历（SOAP校验）
func (s *EMRService) CreateRecord(record *model.MedicalRecord) error {
	if err := s.validateRecord(record); err != nil {
		return err
	}
	return s.repo.CreateRecord(record)
}

// GetRecord 获取病历
func (s *EMRService) GetRecord(id string) (*model.MedicalRecord, error) {
	return s.repo.FindByID(id)
}

// ListRecords 分页查询患者病历
func (s *EMRService) ListRecords(patientID string, page, pageSize int) ([]model.MedicalRecord, int64, error) {
	records, total, err := s.repo.ListByPatient(patientID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for i := range records {
		records[i].PatientName = demo.PatientName(records[i].PatientID)
		if records[i].DoctorID == "" {
			records[i].DoctorID = "demo-doctor"
		}
	}
	return records, total, nil
}

// UpdateRecord 更新病历
func (s *EMRService) UpdateRecord(record *model.MedicalRecord) error {
	if err := s.validateRecord(record); err != nil {
		return err
	}
	return s.repo.UpdateRecord(record)
}

// QualityControl 病历质控
func (s *EMRService) QualityControl(recordID, reviewerID string, level int, comment string) error {
	return s.repo.QualityControl(recordID, reviewerID, level, comment)
}

// ListTemplates 获取病历模板列表
func (s *EMRService) ListTemplates() ([]model.RecordTemplate, error) {
	return s.repo.ListTemplates()
}

// GetTemplate 获取病历模板详情
func (s *EMRService) GetTemplate(id string) (*model.RecordTemplate, error) {
	return s.repo.GetTemplateByID(id)
}

// CDSSCheck 临床决策支持检查
func (s *EMRService) CDSSCheck(patientID, drugID, diagnosis string) ([]string, error) {
	return s.repo.CDSSCheck(patientID, drugID, diagnosis)
}

// validateRecord SOAP结构体校验，各字段不能为空
func (s *EMRService) validateRecord(record *model.MedicalRecord) error {
	if record.PatientID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, errors.MsgPatientIDRequired)
	}
	if record.ChiefComplaint == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "主诉不能为空")
	}
	if record.PresentIllness == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "现病史不能为空")
	}
	if record.PastHistory == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "既往史不能为空")
	}
	if record.PhysicalExam == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "体格检查不能为空")
	}
	if record.AuxiliaryExam == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "辅助检查不能为空")
	}
	if record.Diagnosis == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "诊断不能为空")
	}
	if record.TreatmentPlan == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "处理计划不能为空")
	}
	return nil
}
