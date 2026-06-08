/** 后端统一响应包装 */
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

/** 分页响应数据 */
export interface PageData<T = any> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

/** 分页查询参数 */
export interface PageQuery {
  page?: number
  pageSize?: number
}

/** 用户信息 */
export interface UserInfo {
  id: string
  username: string
  realName: string
  role: string
  deptId: string
  deptName: string
  perms: string[]
}

/** 登录请求 */
export interface LoginRequest {
  username: string
  password: string
}

/** 登录响应 */
export interface LoginResponse {
  token: string
  user: UserInfo
}

/** 患者 */
export interface Patient {
  id: string
  name: string
  gender: string
  age: number
  phone: string
  idCard: string
  address: string
  createdAt: string
}

/** 员工 */
export interface Employee {
  id: string
  name: string
  gender: string
  phone: string
  deptId: string
  deptName: string
  title: string
  role: string
}

/** 科室 */
export interface Department {
  id: string
  name: string
  parentId: string
  description: string
}

/** 挂号记录 */
export interface Registration {
  id: string
  patientId: string
  patientName: string
  scheduleId: string
  deptName: string
  doctorName: string
  registrationDate: string
  status: number
  queueNumber: number
}

/** 号源 (registration 服务) */
export interface Schedule {
  id: string
  doctorId: string
  doctorName: string
  deptId: string
  deptName: string
  date: string
  timeSlot: string
  totalCount: number
  remainCount: number
  status: number
}

/** 排班记录 (schedule 服务) */
export interface ScheduleRecord {
  id: string
  doctorId: string
  doctorName: string
  deptId: string
  deptName: string
  date: string
  timeSlot: string
  maxPatients: number
  currentPatients: number
  status: number
}

/** 门诊记录 */
export interface ClinicRecord {
  id: string
  patientId: string
  patientName: string
  doctorId: string
  doctorName: string
  chiefComplaint: string
  diagnosis: string
  createdAt: string
}

/** 处方 */
export interface Prescription {
  id: string
  patientId: string
  patientName: string
  doctorId: string
  doctorName: string
  status: number
  createdAt: string
  details: PrescriptionDetail[]
}

/** 处方明细 */
export interface PrescriptionDetail {
  id: string
  drugId: string
  drugName: string
  quantity: number
  usage: string
  dosage: string
}

/** 账单 */
export interface Bill {
  id: string
  patientId: string
  patientName: string
  registrationId: string
  billNo: string
  totalAmount: number
  paidAmount: number
  status: number
  createdAt: string
  details: BillDetail[]
}

/** 账单明细 */
export interface BillDetail {
  id: string
  itemName: string
  itemType: string
  quantity: number
  unitPrice: number
  amount: number
}

/** 药品 */
export interface Drug {
  id: string
  name: string
  specification: string
  manufacturer: string
  unit: string
  stock: number
  price: number
  expiryDate: string
}

/** 检查报告 */
export interface ExaminationReport {
  id: string
  patientId: string
  patientName: string
  examType: string
  examItem: string
  result: string
  status: number
  reviewerId: string
  createdAt: string
}

/** 住院记录 */
export interface InpatientRecord {
  id: string
  patientId: string
  patientName: string
  deptId: string
  deptName: string
  bedNo: string
  admitDate: string
  dischargeDate: string
  status: number
}

/** 医嘱 */
export interface MedicalOrder {
  id: string
  inpatientId: string
  orderType: string
  content: string
  doctorId: string
  status: number
  createdAt: string
}

/** 病历 */
export interface EMRRecord {
  id: string
  patientId: string
  patientName: string
  doctorId: string
  templateId: string
  subjective: string
  objective: string
  assessment: string
  plan: string
  status: number
  createdAt: string
}

/** 字典 */
export interface DictItem {
  id: string
  dictType: string
  dictLabel: string
  dictValue: string
  sort: number
  status: number
}

/** 统计 */
export interface OperationStats {
  totalRegistrations: number
  totalClinic: number
  totalPrescriptions: number
  totalBilling: number
  totalInpatient: number
}
