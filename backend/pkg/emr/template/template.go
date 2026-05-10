// Package template 提供结构化病历模板解析引擎
package template

import (
	"regexp"
	"strings"
)

// Section 病历段落定义
type Section struct {
	Key      string `json:"key"`
	Title    string `json:"title"`
	Required bool   `json:"required"`
	Order    int    `json:"order"`
}

// Template 病历模板定义
type Template struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"` // 科室/专科分类
	Description string    `json:"description"`
	Sections    []Section `json:"sections"`
}

// Record 一份具体的病历记录（模板实例）
type Record struct {
	TemplateID string            `json:"templateId"`
	Sections   map[string]string `json:"sections"` // key → 内容
}

// 预定义的 SOAP 段落
var SOAPSections = []Section{
	{Key: "subjective", Title: "主观资料 (Subjective)", Required: true, Order: 1},
	{Key: "objective", Title: "客观资料 (Objective)", Required: true, Order: 2},
	{Key: "assessment", Title: "评估 (Assessment)", Required: true, Order: 3},
	{Key: "plan", Title: "计划 (Plan)", Required: true, Order: 4},
}

// 预定义专科模板
var (
	TemplateGeneral = Template{
		ID:       "general",
		Name:     "通用门诊病历",
		Category: "general",
		Sections: SOAPSections,
	}

	TemplatePediatrics = Template{
		ID:       "pediatrics",
		Name:     "儿科病历",
		Category: "pediatrics",
		Sections: append(SOAPSections, Section{Key: "growth", Title: "生长发育", Required: false, Order: 5}),
	}

	TemplateSurgery = Template{
		ID:       "surgery",
		Name:     "外科病历",
		Category: "surgery",
		Sections: append(SOAPSections, Section{Key: "surgical_history", Title: "手术史", Required: false, Order: 5}),
	}

	TemplateObstetrics = Template{
		ID:       "obstetrics",
		Name:     "产科病历",
		Category: "obstetrics",
		Sections: append(SOAPSections,
			Section{Key: "menstrual_history", Title: "月经史", Required: false, Order: 5},
			Section{Key: "obstetric_history", Title: "孕产史", Required: false, Order: 6},
		),
	}
)

// GetTemplates 获取全部模板
func GetTemplates() []Template {
	return []Template{TemplateGeneral, TemplatePediatrics, TemplateSurgery, TemplateObstetrics}
}

// GetTemplateByID 按 ID 获取模板
func GetTemplateByID(id string) *Template {
	for _, t := range GetTemplates() {
		if t.ID == id {
			return &t
		}
	}
	return nil
}

// Render 将病历记录渲染为结构化文本
func Render(record *Record, tmpl *Template) string {
	var sb strings.Builder

	sb.WriteString("═══════════════════════════════════════\n")
	sb.WriteString("  ")
	sb.WriteString(tmpl.Name)
	sb.WriteString("\n")
	sb.WriteString("═══════════════════════════════════════\n\n")

	for _, section := range tmpl.Sections {
		sb.WriteString("【")
		sb.WriteString(section.Title)
		sb.WriteString("】\n")

		if content, ok := record.Sections[section.Key]; ok && content != "" {
			sb.WriteString(content)
		} else {
			sb.WriteString("（未记录）")
		}
		sb.WriteString("\n\n")
	}

	return sb.String()
}

// Validate 校验病历记录完整性（必填字段检查）
func Validate(record *Record, tmpl *Template) []string {
	var missing []string
	for _, section := range tmpl.Sections {
		if !section.Required {
			continue
		}
		content, ok := record.Sections[section.Key]
		if !ok || strings.TrimSpace(content) == "" {
			missing = append(missing, section.Title)
		}
	}
	return missing
}

// placeholderPattern 匹配 {{...}} 占位符
var placeholderPattern = regexp.MustCompile(`\{\{(\w+)\}\}`)

// RenderWithPlaceholders 渲染含占位符的文本
// 例如："患者{{name}}，年龄{{age}}岁" → "患者张三，年龄35岁"
func RenderWithPlaceholders(text string, params map[string]string) string {
	return placeholderPattern.ReplaceAllStringFunc(text, func(match string) string {
		key := match[2 : len(match)-2] // 去掉 {{ 和 }}
		if val, ok := params[key]; ok {
			return val
		}
		return match
	})
}

// GetPlaceholderKeys 提取文本中所有占位符 key
func GetPlaceholderKeys(text string) []string {
	matches := placeholderPattern.FindAllStringSubmatch(text, -1)
	keys := make([]string, 0, len(matches))
	seen := make(map[string]bool)
	for _, m := range matches {
		if !seen[m[1]] {
			keys = append(keys, m[1])
			seen[m[1]] = true
		}
	}
	return keys
}
