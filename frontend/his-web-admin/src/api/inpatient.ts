import http from './client'
import type { InpatientRecord, MedicalOrder, PageData, PageQuery } from './types'

export const inpatientApi = {
  /** 入院登记 */
  admit(data: Partial<InpatientRecord>): Promise<InpatientRecord> {
    return http.post('/inpatient/admit', data)
  },

  /** 出院 */
  discharge(id: string): Promise<void> {
    return http.post(`/inpatient/discharge/${id}`)
  },

  /** 住院详情 */
  getById(id: string): Promise<InpatientRecord> {
    return http.get(`/inpatient/${id}`)
  },

  /** 住院列表 */
  getList(params: PageQuery & { dept_id?: string; status?: number }): Promise<PageData<InpatientRecord>> {
    return http.get('/inpatient/list', { params })
  },

  /** 创建医嘱 */
  createOrder(data: Partial<MedicalOrder>): Promise<MedicalOrder> {
    return http.post('/inpatient/order', data)
  },

  /** 创建护理记录 */
  createNursing(data: any): Promise<void> {
    return http.post('/inpatient/nursing', data)
  },
}
