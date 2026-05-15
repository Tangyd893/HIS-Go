import http from './client'
import type { Bill, PageData } from './types'

export const billingApi = {
  getBills(params?: Record<string, unknown>): Promise<PageData<Bill>> {
    return http.get('/billing/list', { params })
  },
  createBill(data: Partial<Bill>): Promise<Bill> {
    return http.post('/billing', data)
  },
  updateBill(id: string, data: Partial<Bill>): Promise<Bill> {
    return http.put(`/billing/${id}`, data)
  },
  payBill(id: string): Promise<Bill> {
    return http.post(`/billing/${id}/pay`)
  },
}
