package repository

import (
	"strings"

	"his-go/internal/emr/model"
	"his-go/pkg/errors"

	"gorm.io/gorm"
)

// EMRRepository 电子病历数据仓库
type EMRRepository struct {
	db *gorm.DB
}

// NewEMRRepository 创建电子病历数据仓库
func NewEMRRepository(db *gorm.DB) *EMRRepository {
	return &EMRRepository{db: db}
}

// CreateRecord 创建病历
func (r *EMRRepository) CreateRecord(record *model.MedicalRecord) error {
	return r.db.Create(record).Error
}

// FindByID 根据ID查找病历
func (r *EMRRepository) FindByID(id string) (*model.MedicalRecord, error) {
	var record model.MedicalRecord
	err := r.db.Where("id = ?", id).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "病历不存在")
		}
		return nil, err
	}
	return &record, nil
}

// ListByPatient 分页查询患者病历列表
func (r *EMRRepository) ListByPatient(patientID string, page, pageSize int) ([]model.MedicalRecord, int64, error) {
	var records []model.MedicalRecord
	var total int64

	query := r.db.Model(&model.MedicalRecord{}).Where("patient_id = ?", patientID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// UpdateRecord 更新病历
func (r *EMRRepository) UpdateRecord(record *model.MedicalRecord) error {
	result := r.db.Model(&model.MedicalRecord{}).Where("id = ?", record.ID).Updates(record)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "病历不存在")
	}
	return nil
}

// QualityControl 质控：更新质控等级 + 创建质控记录
func (r *EMRRepository) QualityControl(recordID, reviewerID string, level int, comment string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.MedicalRecord{}).
			Where("id = ?", recordID).
			Updates(map[string]interface{}{
				"quality_level": level,
				"status":        2,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.NewAppError(errors.CodeNotFound, "病历不存在")
		}

		qc := &model.QualityControl{
			RecordID:   recordID,
			ReviewerID: reviewerID,
			Level:      level,
			Comment:    comment,
		}
		if err := tx.Create(qc).Error; err != nil {
			return err
		}

		return nil
	})
}

// ListTemplates 查询所有病历模板
func (r *EMRRepository) ListTemplates() ([]model.RecordTemplate, error) {
	var templates []model.RecordTemplate
	err := r.db.Order("created_at DESC").Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// GetTemplateByID 根据ID查找病历模板
func (r *EMRRepository) GetTemplateByID(id string) (*model.RecordTemplate, error) {
	var tmpl model.RecordTemplate
	err := r.db.Where("id = ?", id).First(&tmpl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "模板不存在")
		}
		return nil, err
	}
	return &tmpl, nil
}

// CDSSCheck 临床决策支持检查：查询患者过敏史 + 药物相互作用（简化版）
func (r *EMRRepository) CDSSCheck(patientID, drugID, diagnosis string) ([]string, error) {
	var warnings []string

	var records []model.MedicalRecord
	if err := r.db.Where("patient_id = ?", patientID).Find(&records).Error; err != nil {
		return nil, err
	}

	for _, record := range records {
		if strings.Contains(record.PastHistory, "过敏") ||
			strings.Contains(record.PastHistory, "allergy") ||
			strings.Contains(record.ChiefComplaint, "过敏") {
			warnings = append(warnings, "患者存在过敏史，请谨慎用药")
			break
		}
	}

	for _, record := range records {
		if strings.Contains(record.PastHistory, "高血压") && strings.Contains(diagnosis, "高血压") {
			warnings = append(warnings, "患者有高血压病史，请关注药物相互作用")
			break
		}
		if strings.Contains(record.PastHistory, "糖尿病") && strings.Contains(diagnosis, "糖尿病") {
			warnings = append(warnings, "患者有糖尿病病史，请关注药物相互作用")
			break
		}
		if strings.Contains(record.PastHistory, "肝") || strings.Contains(record.PastHistory, "肾") {
			warnings = append(warnings, "患者肝肾功能异常，请注意药物剂量调整")
			break
		}
	}

	return warnings, nil
}
