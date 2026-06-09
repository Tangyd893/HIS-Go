package model

import "time"

// InpatientRecord 住院记录模型，表名 inpatient_records
type InpatientRecord struct {
	ID            string     `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	PatientID     string     `gorm:"column:patient_id;not null;type:varchar(64);index" json:"patientId"`
	PatientName   string     `gorm:"column:patient_name;not null;size:50" json:"patientName"`
	AdmissionDate *time.Time `gorm:"column:admission_date" json:"admitDate,omitempty"`
	DischargeDate *time.Time `gorm:"column:discharge_date" json:"dischargeDate,omitempty"`
	DeptID        string     `gorm:"column:dept_id;not null;type:varchar(64)" json:"deptId"`
	DeptName      string     `gorm:"-" json:"deptName,omitempty"`
	RoomNo        string     `gorm:"column:room_no;size:20" json:"roomNo,omitempty"`
	BedNo         string     `gorm:"column:bed_no;size:20" json:"bedNo"`
	Diagnosis     string     `gorm:"column:diagnosis;type:text" json:"diagnosis,omitempty"`
	Deposit       float64    `gorm:"column:deposit;not null;default:0" json:"deposit,omitempty"`
	TotalCost     float64    `gorm:"column:total_cost;not null;default:0" json:"totalCost,omitempty"`
	Status        int8       `gorm:"column:status;default:0" json:"status"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
}

func (InpatientRecord) TableName() string { return "inpatient_records" }
