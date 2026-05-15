import http from './client'
import type { InpatientRecord, DoctorOrder, PageData } from './types'

export const inpatientApi = {
  getRecords(params?: Record<string, unknown>): Promise<PageData<InpatientRecord>> {
    return http.get('/inpatient/records', { params })
  },
  createRecord(data: Partial<InpatientRecord>): Promise<InpatientRecord> {
    return http.post('/inpatient/records', data)
  },
  updateRecord(id: string, data: Partial<InpatientRecord>): Promise<InpatientRecord> {
    return http.put(`/inpatient/records/${id}`, data)
  },
  getOrders(inpatientId: string): Promise<PageData<DoctorOrder>> {
    return http.get(`/inpatient/orders`, { params: { inpatientId } })
  },
  createOrder(data: Partial<DoctorOrder>): Promise<DoctorOrder> {
    return http.post('/inpatient/orders', data)
  },
}
