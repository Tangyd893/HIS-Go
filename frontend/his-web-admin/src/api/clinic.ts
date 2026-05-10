import http from './client'
import type { ClinicRecord, PageData, PageQuery } from './types'

export const clinicApi = {
  /** 创建门诊记录 */
  create(data: Partial<ClinicRecord>): Promise<ClinicRecord> {
    return http.post('/clinic/record', data)
  },

  /** 获取门诊详情 */
  getById(id: string): Promise<ClinicRecord> {
    return http.get(`/clinic/record/${id}`)
  },

  /** 门诊记录列表 */
  getList(params: PageQuery & { patientId?: string }): Promise<PageData<ClinicRecord>> {
    return http.get('/clinic/records', { params })
  },

  /** 创建检查申请 */
  createExamination(data: any): Promise<void> {
    return http.post('/clinic/examination-request', data)
  },
}
