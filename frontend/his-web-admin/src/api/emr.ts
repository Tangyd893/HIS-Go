import http from './client'
import type { EMRRecord, PageData, PageQuery } from './types'

export const emrApi = {
  /** 创建病历 */
  create(data: Partial<EMRRecord>): Promise<EMRRecord> {
    return http.post('/emr/record', data)
  },

  /** 病历详情 */
  getById(id: string): Promise<EMRRecord> {
    return http.get(`/emr/record/${id}`)
  },

  /** 病历列表 */
  getList(params: PageQuery & { patient_id?: string }): Promise<PageData<EMRRecord>> {
    return http.get('/emr/records', { params })
  },

  /** 质控 */
  qualityControl(id: string, data: { reviewer_id: string; level: number; comment?: string }): Promise<void> {
    return http.post(`/emr/quality-control/${id}`, data)
  },

  /** 病历模板 */
  getTemplates(): Promise<any[]> {
    return http.get('/emr/templates')
  },

  /** CDSS 检查 */
  cdssCheck(data: { patient_id: string; drug_id?: string; diagnosis?: string }): Promise<any> {
    return http.post('/emr/cdss-check', data)
  },
}
