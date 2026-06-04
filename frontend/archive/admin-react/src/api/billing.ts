import http from './client'
import type { Bill, PageData } from './types'

export const billingApi = {
  getBills(params?: Record<string, unknown>): Promise<PageData<Bill>> {
    return http.get('/billing/list', { params })
  },
  getBill(id: string): Promise<Bill> {
    return http.get(`/billing/${id}`)
  },
  createBill(data: Partial<Bill>): Promise<Bill> {
    return http.post('/billing/create', data)
  },
  payBill(id: string): Promise<Bill> {
    return http.post(`/billing/pay/${id}`)
  },
  refundBill(id: string): Promise<Bill> {
    return http.post(`/billing/refund/${id}`)
  },
}
