import http from './client'

/** 登录 */
export function login(data: { username: string; password: string }): Promise<{ token: string; userInfo?: any; user?: any }> {
  return http.post('/auth/login', data)
}

/** 患者健康档案 */
export function getHealthSummary(patientId: string): Promise<any> {
  return http.get(`/health-record/summary/${patientId}`)
}

export function getHealthTimeline(patientId: string): Promise<any[]> {
  return http.get(`/health-record/timeline/${patientId}`)
}

/** 科室 */
export function getDepartments(): Promise<{ id: string; name: string }[]> {
  return http.get('/user/departments')
}

/** 挂号 */
export function getSchedules(params: { deptId?: string; date?: string }): Promise<any[]> {
  return http.get('/registration/schedules', { params })
}

export function getRegistrations(params: { patientId?: string; page?: number; pageSize?: number }): Promise<{ list?: any[]; total?: number }> {
  return http.get('/registration/list', { params })
}

export function register(data: { patientId: string; patientName: string; scheduleId: string }): Promise<any> {
  return http.post('/registration/register', data)
}

/** 取消挂号 */
export function cancelRegistration(id: string): Promise<void> {
  return http.post(`/registration/cancel/${id}`)
}

/** 处方 */
export function getPrescriptions(params: any): Promise<any> {
  return http.get('/prescription/list', { params })
}

/** 检查报告 */
export function getReports(params: any): Promise<any> {
  return http.get('/examination/reports', { params })
}

/** 随访 */
export function getFollowupPlans(params: any): Promise<any> {
  return http.get('/followup/plans', { params })
}

/** 慢病管理 */
export function getContract(patientId: string): Promise<any> {
  return http.get('/outpatient/contract', { params: { patientId } })
}

export function createContract(data: any): Promise<any> {
  return http.post('/outpatient/contract', data)
}

export function reportHealthData(data: any): Promise<void> {
  return http.post('/outpatient/health-data', data)
}

export function getHealthData(patientId: string): Promise<any[]> {
  return http.get('/outpatient/health-data', { params: { patientId } })
}

/** 就诊助手 */
export interface TriageResponse {
  symptom: string
  advice: string
  depts: { id: string; name: string }[]
  knowledgeRef: string
  urgency: string
  mode: string
  disclaimer: string
}

export function triageChat(symptom: string): Promise<TriageResponse> {
  return http.post('/outpatient/assistant/chat', { symptom })
}
