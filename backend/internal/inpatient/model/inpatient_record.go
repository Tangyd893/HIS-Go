package model

import "time"

// InpatientRecord 住院记录模型，表名 inpatient_records
type InpatientRecord struct {
	ID            string     `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	PatientID     string     `gorm:"column:patient_id;not null;type:varchar(64);index"`
	PatientName   string     `gorm:"column:patient_name;not null;size:50"`
	AdmissionDate *time.Time `gorm:"column:admission_date;comment:'入院日期'"`
	DischargeDate *time.Time `gorm:"column:discharge_date;comment:'出院日期'"`
	DeptID        string     `gorm:"column:dept_id;not null;type:varchar(64)"`
	RoomNo        string     `gorm:"column:room_no;size:20"`
	BedNo         string     `gorm:"column:bed_no;size:20"`
	Diagnosis     string     `gorm:"column:diagnosis;type:text;comment:'入院诊断'"`
	Deposit       float64    `gorm:"column:deposit;not null;default:0;comment:'押金'"`
	TotalCost     float64    `gorm:"column:total_cost;not null;default:0;comment:'总费用'"`
	Status        int8       `gorm:"column:status;default:0;comment:'0待入院 1住院中 2已出院'"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (InpatientRecord) TableName() string { return "inpatient_records" }
