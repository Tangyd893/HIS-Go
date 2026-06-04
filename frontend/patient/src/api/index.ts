import http from './client'

/** 登录 */
export function login(data: { username: string; password: string }): Promise<{ token: string; user: any }> {
  return http.post('/auth/login', data)
}

/** 患者健康档案 */
export function getHealthSummary(patientId: string): Promise<any> {
  return http.get(`/health-record/summary/${patientId}`)
}

export function getHealthTimeline(patientId: string): Promise<any[]> {
  return http.get(`/health-record/timeline/${patientId}`)
}

/** 挂号 */
export function getSchedules(params: { deptId?: string; date?: string }): Promise<any[]> {
  return http.get('/registration/schedules', { params })
}

export function register(data: { patientId: string; patientName: string; scheduleId: string }): Promise<any> {
  return http.post('/registration/register', data)
}

/** 在线问诊 */
export function createConsultation(data: any): Promise<any> {
  return http.post('/outpatient/consultation', data)
}

export function getConsultations(params: any): Promise<any> {
  return http.get('/outpatient/consultations', { params })
}

export function sendMessage(data: any): Promise<void> {
  return http.post('/outpatient/message', data)
}

export function getMessages(consultationId: string): Promise<any[]> {
  return http.get('/outpatient/messages', { params: { consultationId } })
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
export function createContract(data: any): Promise<any> {
  return http.post('/outpatient/contract', data)
}

export function reportHealthData(data: any): Promise<void> {
  return http.post('/outpatient/health-data', data)
}

export function getHealthData(patientId: string): Promise<any[]> {
  return http.get('/outpatient/health-data', { params: { patientId } })
}
