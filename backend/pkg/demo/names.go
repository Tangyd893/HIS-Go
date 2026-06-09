package demo

// PatientName 演示环境患者 ID → 姓名（billing 等跨库服务无法 join his_user 时使用）
func PatientName(id string) string {
	if n, ok := patientNames[id]; ok {
		return n
	}
	return id
}

// DoctorName 演示环境医生 ID → 姓名
func DoctorName(id string) string {
	if n, ok := doctorNames[id]; ok {
		return n
	}
	return id
}

var patientNames = map[string]string{
	"patient_001": "王小明",
	"patient_002": "李小红",
	"patient_003": "张三",
	"patient_004": "赵四",
	"patient_005": "孙七",
	"patient_006": "周八",
	"patient_007": "吴九",
	"patient_008": "郑十",
	"patient_009": "陈十一",
	"patient_010": "刘十二",
}

var deptNames = map[string]string{
	"dept_001": "内科",
	"dept_002": "外科",
	"dept_003": "儿科",
	"dept_004": "妇产科",
	"dept_005": "急诊科",
}

// DeptName 演示环境科室 ID → 名称
func DeptName(id string) string {
	if n, ok := deptNames[id]; ok {
		return n
	}
	return id
}

var doctorNames = map[string]string{
	"demo-doctor":      "张医生",
	"doctor-wang":      "王医生",
	"doctor-li":        "李医生",
	"doctor-zhao":      "赵医生",
	"doctor-chen":      "陈医生",
	"demo-nurse":       "李护士",
	"demo-pharmacist":  "王药师",
}
