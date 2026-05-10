// Package cdss 提供临床决策支持系统（CDSS）规则引擎
package cdss

import "strings"

// AlertLevel 告警级别
type AlertLevel int

const (
	LevelInfo    AlertLevel = 0 // 信息提示
	LevelWarning AlertLevel = 1 // 警告
	LevelError   AlertLevel = 2 // 严重警告/禁止
)

func (l AlertLevel) String() string {
	switch l {
	case LevelInfo:
		return "提示"
	case LevelWarning:
		return "警告"
	case LevelError:
		return "禁止"
	default:
		return "未知"
	}
}

// Alert CDSS 告警结果
type Alert struct {
	Level      AlertLevel `json:"level"`
	Message    string     `json:"message"`
	Category   string     `json:"category"` // drug-interaction, allergy, dosage, etc.
	DrugName   string     `json:"drugName,omitempty"`
	Suggestion string     `json:"suggestion,omitempty"`
}

// CheckRequest CDSS 检查请求
type CheckRequest struct {
	PatientID string   `json:"patientId"`
	DrugIDs   []string `json:"drugIds,omitempty"`
	DrugNames []string `json:"drugNames,omitempty"`
	Diagnosis string   `json:"diagnosis,omitempty"`
	Allergies []string `json:"allergies,omitempty"`
	Age       int      `json:"age,omitempty"`
	Weight    float64  `json:"weight,omitempty"`
	Pregnant  bool     `json:"pregnant,omitempty"`
}

// CheckResult CDSS 检查结果
type CheckResult struct {
	Alerts     []Alert `json:"alerts"`
	Passed     bool    `json:"passed"`     // 无严重警告则为通过
	TotalCheck int     `json:"totalCheck"` // 检查总数
}

// Rule 检查规则接口
type Rule interface {
	Name() string
	Check(req *CheckRequest) []Alert
}

// Engine CDSS 规则引擎
type Engine struct {
	rules []Rule
}

// NewEngine 创建 CDSS 引擎
func NewEngine() *Engine {
	e := &Engine{}
	e.RegisterDefaultRules()
	return e
}

// Register 注册规则
func (e *Engine) Register(rule Rule) {
	e.rules = append(e.rules, rule)
}

// RegisterDefaultRules 注册默认规则
func (e *Engine) RegisterDefaultRules() {
	e.Register(&AllergyCheckRule{})
	e.Register(&DrugInteractionRule{})
	e.Register(&DosageCheckRule{})
}

// Check 执行全部规则检查
func (e *Engine) Check(req *CheckRequest) *CheckResult {
	result := &CheckResult{Passed: true, TotalCheck: len(e.rules)}
	for _, rule := range e.rules {
		alerts := rule.Check(req)
		result.Alerts = append(result.Alerts, alerts...)
		for _, a := range alerts {
			if a.Level >= LevelError {
				result.Passed = false
			}
		}
	}
	return result
}

// ==================== 内置规则实现 ====================

// AllergyCheckRule 过敏史检查规则
type AllergyCheckRule struct{}

func (r *AllergyCheckRule) Name() string { return "过敏史检查" }

func (r *AllergyCheckRule) Check(req *CheckRequest) []Alert {
	var alerts []Alert

	// 青霉素交叉过敏
	penicillinDrugs := map[string]bool{"青霉素": true, "阿莫西林": true, "氨苄西林": true, "头孢氨苄": true}
	sulfaAllergies := map[string]bool{"磺胺类": true}

	for _, drug := range req.DrugNames {
		if penicillinDrugs[drug] {
			for _, allergy := range req.Allergies {
				if strings.Contains(allergy, "青霉素") {
					alerts = append(alerts, Alert{
						Level:      LevelError,
						Message:    "患者有青霉素过敏史，禁止使用 " + drug,
						Category:   "allergy",
						DrugName:   drug,
						Suggestion: "建议更换为非青霉素类抗生素，如大环内酯类",
					})
				}
			}
		}

		if sulfaAllergies[drug] {
			for _, allergy := range req.Allergies {
				if strings.Contains(allergy, "磺胺") {
					alerts = append(alerts, Alert{
						Level:    LevelError,
						Message:  "患者有磺胺类过敏史，禁止使用 " + drug,
						Category: "allergy",
						DrugName: drug,
					})
				}
			}
		}
	}
	return alerts
}

// DrugInteractionRule 药物相互作用检查规则
type DrugInteractionRule struct{}

func (r *DrugInteractionRule) Name() string { return "药物相互作用" }

// 已知药物相互作用对
var interactionPairs = []struct {
	a, b    string
	level   AlertLevel
	message string
}{
	{"华法林", "阿司匹林", LevelError, "华法林与阿司匹林联用显著增加出血风险"},
	{"华法林", "布洛芬", LevelWarning, "华法林与布洛芬联用增加出血风险"},
	{"ACEI类", "螺内酯", LevelWarning, "ACEI类药物与螺内酯联用易致高钾血症"},
	{"甲氨蝶呤", "磺胺类", LevelWarning, "甲氨蝶呤与磺胺类联用增加骨髓抑制风险"},
	{"地高辛", "胺碘酮", LevelWarning, "地高辛与胺碘酮联用增加地高辛血药浓度"},
}

func (r *DrugInteractionRule) Check(req *CheckRequest) []Alert {
	var alerts []Alert
	if len(req.DrugNames) < 2 {
		return alerts
	}

	for _, pair := range interactionPairs {
		foundA, foundB := false, false
		for _, drug := range req.DrugNames {
			if strings.Contains(drug, pair.a) {
				foundA = true
			}
			if strings.Contains(drug, pair.b) {
				foundB = true
			}
		}
		if foundA && foundB {
			alerts = append(alerts, Alert{
				Level:      pair.level,
				Message:    pair.message,
				Category:   "drug-interaction",
				Suggestion: "建议评估风险后调整用药方案",
			})
		}
	}
	return alerts
}

// DosageCheckRule 剂量校验规则
type DosageCheckRule struct{}

func (r *DosageCheckRule) Name() string { return "剂量校验" }

// 特殊人群用药禁忌
func (r *DosageCheckRule) Check(req *CheckRequest) []Alert {
	var alerts []Alert

	// 孕妇禁用药物
	pregnancyContraindicated := map[string]bool{
		"利巴韦林": true, "异维A酸": true, "甲氨蝶呤": true,
		"华法林": true, "卡马西平": true, "丙戊酸钠": true,
	}
	if req.Pregnant {
		for _, drug := range req.DrugNames {
			if pregnancyContraindicated[drug] {
				alerts = append(alerts, Alert{
					Level:      LevelError,
					Message:    drug + " 在孕期禁用，可能导致胎儿畸形",
					Category:   "dosage",
					DrugName:   drug,
					Suggestion: "请咨询妇产科医生选择安全的替代药物",
				})
			}
		}
	}

	// 儿童禁用药物
	pediatricContraindicated := map[string]bool{
		"四环素": true, "氯霉素": true, "喹诺酮类": true,
	}
	if req.Age > 0 && req.Age <= 12 {
		for _, drug := range req.DrugNames {
			if pediatricContraindicated[drug] {
				alerts = append(alerts, Alert{
					Level:      LevelWarning,
					Message:    drug + " 在儿童中应谨慎使用",
					Category:   "dosage",
					DrugName:   drug,
					Suggestion: "请评估儿童用药安全性",
				})
			}
		}
	}

	return alerts
}
