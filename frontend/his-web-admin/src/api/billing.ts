import http from './client'
import type { Bill, PageData, PageQuery } from './types'

export const billingApi = {
  /** 创建账单 */
  create(data: Partial<Bill>): Promise<Bill> {
    return http.post('/billing/create', data)
  },

  /** 账单详情 */
  getById(id: string): Promise<Bill> {
    return http.get(`/billing/${id}`)
  },

  /** 支付 */
  pay(id: string, payMethod: number): Promise<void> {
    return http.post(`/billing/pay/${id}`, { pay_method: payMethod })
  },

  /** 退款 */
  refund(id: string): Promise<void> {
    return http.post(`/billing/refund/${id}`)
  },

  /** 账单列表 */
  getList(params: PageQuery & { patient_id?: string; status?: number }): Promise<PageData<Bill>> {
    return http.get('/billing/list', { params })
  },
}
