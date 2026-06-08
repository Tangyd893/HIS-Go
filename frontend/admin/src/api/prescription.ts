import http from './client'
import type { Prescription, PageData, PageQuery } from './types'

export const prescriptionApi = {
  /** 创建处方（嵌套格式：{prescription: {...}, details: [...]}） */
  create(data: { prescription: Partial<Prescription>; details: any[] }): Promise<Prescription> {
    return http.post('/prescription/create', data)
  },

  /** 处方详情 */
  getById(id: string): Promise<Prescription> {
    return http.get(`/prescription/${id}`)
  },

  /** 处方列表 */
  getList(params: PageQuery & { patientId?: string }): Promise<PageData<Prescription>> {
    return http.get('/prescription/list', { params })
  },

  /** 审核处方 */
  review(data: { id: string; approved: boolean; comment?: string }): Promise<void> {
    return http.post('/prescription/review', data)
  },

  /** 取消处方 */
  cancel(id: string): Promise<void> {
    return http.post(`/prescription/cancel/${id}`)
  },
}
