package repository

import (
	"fmt"

	"gorm.io/gorm"

	"his-go/internal/outpatient/model"
	"his-go/pkg/errors"
)

// OutpatientRepository 院外患者数据仓库
type OutpatientRepository struct {
	db *gorm.DB
}

// NewOutpatientRepository 创建院外患者数据仓库
func NewOutpatientRepository(db *gorm.DB) *OutpatientRepository {
	return &OutpatientRepository{db: db}
}

// CreateConsultation 创建在线问诊
func (r *OutpatientRepository) CreateConsultation(c *model.Consultation) error {
	if err := r.db.Create(c).Error; err != nil {
		return fmt.Errorf("创建问诊失败: %w", err)
	}
	return nil
}

// FindConsultationByID 根据ID查询问诊
func (r *OutpatientRepository) FindConsultationByID(id string) (*model.Consultation, error) {
	var c model.Consultation
	if err := r.db.Where("id = ?", id).First(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "问诊不存在")
		}
		return nil, fmt.Errorf("查询问诊失败: %w", err)
	}
	return &c, nil
}

// ListConsultations 分页查询问诊列表
func (r *OutpatientRepository) ListConsultations(patientID, doctorID string, status int, page, pageSize int) ([]model.Consultation, int64, error) {
	var list []model.Consultation
	var total int64

	query := r.db.Model(&model.Consultation{})
	if patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}
	if doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计问诊失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("查询问诊列表失败: %w", err)
	}

	return list, total, nil
}

// SendMessage 发送问诊消息
func (r *OutpatientRepository) SendMessage(msg *model.ConsultationMessage) error {
	if err := r.db.Create(msg).Error; err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	return nil
}

// GetMessages 获取问诊消息列表
func (r *OutpatientRepository) GetMessages(consultationID string) ([]model.ConsultationMessage, error) {
	var list []model.ConsultationMessage
	if err := r.db.Where("consultation_id = ?", consultationID).
		Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("查询消息失败: %w", err)
	}
	return list, nil
}

// CreateContract 创建慢病签约
func (r *OutpatientRepository) CreateContract(contract *model.ChronicContract) error {
	if err := r.db.Create(contract).Error; err != nil {
		return fmt.Errorf("创建签约失败: %w", err)
	}
	return nil
}

// FindContract 查询慢病签约
func (r *OutpatientRepository) FindContract(patientID, doctorID string) (*model.ChronicContract, error) {
	var contract model.ChronicContract
	if err := r.db.Where("patient_id = ? AND doctor_id = ? AND status = 1", patientID, doctorID).
		First(&contract).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "签约不存在")
		}
		return nil, fmt.Errorf("查询签约失败: %w", err)
	}
	return &contract, nil
}

// ReportHealthData 上报健康数据
func (r *OutpatientRepository) ReportHealthData(data *model.HealthData) error {
	if err := r.db.Create(data).Error; err != nil {
		return fmt.Errorf("上报健康数据失败: %w", err)
	}
	return nil
}

// ListHealthData 查询健康数据列表
func (r *OutpatientRepository) ListHealthData(patientID string) ([]model.HealthData, error) {
	var list []model.HealthData
	if err := r.db.Where("patient_id = ?", patientID).
		Order("measure_time DESC").Limit(100).Find(&list).Error; err != nil {
		return nil, fmt.Errorf("查询健康数据失败: %w", err)
	}
	return list, nil
}
