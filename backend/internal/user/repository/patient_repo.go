package repository

import (
	"gorm.io/gorm"

	"his-go/internal/user/model"
)

// PatientRepository 患者数据仓库
type PatientRepository struct {
	db *gorm.DB
}

// NewPatientRepository 创建患者数据仓库
func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

// FindByID 根据ID查找患者
func (r *PatientRepository) FindByID(id string) (*model.Patient, error) {
	var patient model.Patient
	err := r.db.Where("id = ?", id).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// FindByUserID 根据登录用户 ID 查找患者
func (r *PatientRepository) FindByUserID(userID string) (*model.Patient, error) {
	var patient model.Patient
	err := r.db.Where("user_id = ?", userID).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// FindByIdCard 根据身份证号查找患者
func (r *PatientRepository) FindByIdCard(idCard string) (*model.Patient, error) {
	var patient model.Patient
	err := r.db.Where("id_card = ?", idCard).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// List 分页查询患者列表
func (r *PatientRepository) List(name, phone string, page, pageSize int) ([]model.Patient, int64, error) {
	var patients []model.Patient
	var total int64

	query := r.db.Model(&model.Patient{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&patients).Error; err != nil {
		return nil, 0, err
	}

	return patients, total, nil
}

// Create 创建患者
func (r *PatientRepository) Create(patient *model.Patient) error {
	return r.db.Create(patient).Error
}

// Update 更新患者信息
func (r *PatientRepository) Update(patient *model.Patient) error {
	return r.db.Save(patient).Error
}

// Delete 逻辑删除患者
func (r *PatientRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Patient{}).Error
}
