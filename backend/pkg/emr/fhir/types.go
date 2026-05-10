// Package fhir 提供 HL7 FHIR R4 标准数据交换接口
package fhir

import (
	"encoding/json"
	"time"
)

// ResourceType FHIR 资源类型
type ResourceType string

const (
	ResourcePatient            ResourceType = "Patient"
	ResourceObservation        ResourceType = "Observation"
	ResourceCondition          ResourceType = "Condition"
	ResourceMedicationRequest  ResourceType = "MedicationRequest"
	ResourceProcedure          ResourceType = "Procedure"
	ResourceAllergyIntolerance ResourceType = "AllergyIntolerance"
	ResourceDocumentReference  ResourceType = "DocumentReference"
)

// FHIRResource 基础资源接口
type FHIRResource interface {
	GetResourceType() ResourceType
	GetID() string
}

// BaseResource 资源基类
type BaseResource struct {
	ResourceType ResourceType `json:"resourceType"`
	ID           string       `json:"id"`
	Meta         *Meta        `json:"meta,omitempty"`
}

func (r *BaseResource) GetResourceType() ResourceType { return r.ResourceType }
func (r *BaseResource) GetID() string                 { return r.ID }

// Meta 资源元数据
type Meta struct {
	VersionID   string    `json:"versionId,omitempty"`
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
	Source      string    `json:"source,omitempty"`
	Profile     []string  `json:"profile,omitempty"`
}

// Coding 编码
type Coding struct {
	System  string `json:"system,omitempty"`
	Code    string `json:"code,omitempty"`
	Display string `json:"display,omitempty"`
}

// CodeableConcept 可编码概念
type CodeableConcept struct {
	Coding []Coding `json:"coding,omitempty"`
	Text   string   `json:"text,omitempty"`
}

// PatientResource FHIR 患者资源
type PatientResource struct {
	BaseResource
	Identifier    []Identifier     `json:"identifier,omitempty"`
	Name          []HumanName      `json:"name,omitempty"`
	Gender        string           `json:"gender,omitempty"`
	BirthDate     string           `json:"birthDate,omitempty"`
	Telecom       []ContactPoint   `json:"telecom,omitempty"`
	Address       []Address        `json:"address,omitempty"`
	MaritalStatus *CodeableConcept `json:"maritalStatus,omitempty"`
}

// Identifier 标识符
type Identifier struct {
	Use    string `json:"use,omitempty"`
	System string `json:"system,omitempty"`
	Value  string `json:"value,omitempty"`
}

// HumanName 姓名
type HumanName struct {
	Use    string   `json:"use,omitempty"`
	Family string   `json:"family,omitempty"`
	Given  []string `json:"given,omitempty"`
	Text   string   `json:"text,omitempty"`
}

// ContactPoint 联系方式
type ContactPoint struct {
	System string `json:"system,omitempty"`
	Value  string `json:"value,omitempty"`
	Use    string `json:"use,omitempty"`
}

// Address 地址
type Address struct {
	Use        string   `json:"use,omitempty"`
	Line       []string `json:"line,omitempty"`
	City       string   `json:"city,omitempty"`
	District   string   `json:"district,omitempty"`
	PostalCode string   `json:"postalCode,omitempty"`
	Country    string   `json:"country,omitempty"`
}

// ObservationResource FHIR 观察资源（检查检验）
type ObservationResource struct {
	BaseResource
	Status               string            `json:"status"`
	Category             []CodeableConcept `json:"category,omitempty"`
	Code                 CodeableConcept   `json:"code"`
	Subject              *Reference        `json:"subject,omitempty"`
	EffectiveDateTime    string            `json:"effectiveDateTime,omitempty"`
	ValueQuantity        *Quantity         `json:"valueQuantity,omitempty"`
	ValueString          string            `json:"valueString,omitempty"`
	ValueCodeableConcept *CodeableConcept  `json:"valueCodeableConcept,omitempty"`
	Interpretation       []CodeableConcept `json:"interpretation,omitempty"`
}

// Quantity 数量
type Quantity struct {
	Value  float64 `json:"value"`
	Unit   string  `json:"unit,omitempty"`
	System string  `json:"system,omitempty"`
	Code   string  `json:"code,omitempty"`
}

// Reference 引用
type Reference struct {
	Reference string `json:"reference,omitempty"`
	Display   string `json:"display,omitempty"`
}

// ConditionResource FHIR 诊断病症资源
type ConditionResource struct {
	BaseResource
	ClinicalStatus *CodeableConcept `json:"clinicalStatus,omitempty"`
	Code           *CodeableConcept `json:"code,omitempty"`
	Subject        *Reference       `json:"subject,omitempty"`
	OnsetDateTime  string           `json:"onsetDateTime,omitempty"`
	RecordedDate   string           `json:"recordedDate,omitempty"`
}

// MedicationRequestResource FHIR 处方资源
type MedicationRequestResource struct {
	BaseResource
	Status                    string              `json:"status"`
	Intent                    string              `json:"intent"`
	MedicationCodeableConcept *CodeableConcept    `json:"medicationCodeableConcept,omitempty"`
	Subject                   *Reference          `json:"subject"`
	Requester                 *Reference          `json:"requester,omitempty"`
	DosageInstruction         []DosageInstruction `json:"dosageInstruction,omitempty"`
	DispenseRequest           *DispenseRequest    `json:"dispenseRequest,omitempty"`
}

// DosageInstruction 用药指导
type DosageInstruction struct {
	Text        string     `json:"text,omitempty"`
	Timing      *Timing    `json:"timing,omitempty"`
	DoseAndRate []DoseRate `json:"doseAndRate,omitempty"`
}

// Timing 服药时间
type Timing struct {
	Code *CodeableConcept `json:"code,omitempty"`
}

// DoseRate 剂量
type DoseRate struct {
	DoseQuantity *Quantity `json:"doseQuantity,omitempty"`
}

// DispenseRequest 发药请求
type DispenseRequest struct {
	ValidityPeriod         *Period   `json:"validityPeriod,omitempty"`
	Quantity               *Quantity `json:"quantity,omitempty"`
	NumberOfRepeatsAllowed int       `json:"numberOfRepeatsAllowed,omitempty"`
}

// Period 时间段
type Period struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

// Bundle FHIR 资源集合
type Bundle struct {
	ResourceType string        `json:"resourceType"`
	Type         string        `json:"type"`
	Total        int           `json:"total,omitempty"`
	Entry        []BundleEntry `json:"entry,omitempty"`
}

// BundleEntry 集合条目
type BundleEntry struct {
	Resource json.RawMessage `json:"resource"`
}

// Serialize 序列化 FHIR 资源为 JSON
func Serialize(resource FHIRResource) ([]byte, error) {
	return json.Marshal(resource)
}

// Deserialize 从 JSON 反序列化为 FHIR 资源
func Deserialize(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
