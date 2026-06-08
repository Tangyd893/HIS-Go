package service

import (
	"strconv"
	"strings"

	"his-go/internal/outpatient/model"
	"his-go/internal/outpatient/repository"
)

// OutpatientService 院外患者业务服务
type OutpatientService struct {
	repo *repository.OutpatientRepository
}

// NewOutpatientService 创建院外患者业务服务
func NewOutpatientService(repo *repository.OutpatientRepository) *OutpatientService {
	return &OutpatientService{repo: repo}
}

// CreateConsultation 创建在线问诊
func (s *OutpatientService) CreateConsultation(c *model.Consultation) error {
	c.Status = 0 // 待接诊
	return s.repo.CreateConsultation(c)
}

// GetConsultation 获取问诊详情
func (s *OutpatientService) GetConsultation(id string) (*model.Consultation, error) {
	return s.repo.FindConsultationByID(id)
}

// ListConsultations 分页查询问诊
func (s *OutpatientService) ListConsultations(patientID, doctorID string, status int, page, pageSize int) ([]model.Consultation, int64, error) {
	return s.repo.ListConsultations(patientID, doctorID, status, page, pageSize)
}

// SendMessage 发送问诊消息
func (s *OutpatientService) SendMessage(msg *model.ConsultationMessage) error {
	return s.repo.SendMessage(msg)
}

// GetMessages 获取问诊消息
func (s *OutpatientService) GetMessages(consultationID string) ([]model.ConsultationMessage, error) {
	return s.repo.GetMessages(consultationID)
}

// CreateChronicContract 创建慢病签约
func (s *OutpatientService) CreateChronicContract(contract *model.ChronicContract) error {
	contract.Status = 1 // 签约中
	return s.repo.CreateContract(contract)
}

// GetContract 查询患者慢病签约
func (s *OutpatientService) GetContract(patientID string) (*model.ChronicContract, error) {
	return s.repo.FindContract(patientID, "")
}

// ReportHealthData 上报健康数据并自动标记异常
func (s *OutpatientService) ReportHealthData(data *model.HealthData) error {
	data.Abnormal = s.CheckHealthDataAbnormal(data)
	return s.repo.ReportHealthData(data)
}

// ListHealthData 查询健康数据
func (s *OutpatientService) ListHealthData(patientID string) ([]model.HealthData, error) {
	return s.repo.ListHealthData(patientID)
}

// CheckHealthDataAbnormal 判断健康数据是否异常
func (s *OutpatientService) CheckHealthDataAbnormal(data *model.HealthData) bool {
	switch data.DataType {
	case "blood_pressure":
		return s.isBloodPressureAbnormal(data.Value)
	case "blood_sugar":
		return s.isBloodSugarAbnormal(data.Value)
	default:
		return false
	}
}

// isBloodPressureAbnormal 判断血压是否异常（>140/90 或 <90/60）
func (s *OutpatientService) isBloodPressureAbnormal(value string) bool {
	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return false
	}
	systolic, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	diastolic, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		return false
	}
	return systolic > 140 || diastolic > 90 || systolic < 90 || diastolic < 60
}

// isBloodSugarAbnormal 判断血糖是否异常（空腹血糖>7.0 或 <3.9）
func (s *OutpatientService) isBloodSugarAbnormal(value string) bool {
	val, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return false
	}
	return val > 7.0 || val < 3.9
}
