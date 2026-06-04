import http from './client'
import type { Prescription, PageData } from './types'

export const prescriptionApi = {
  getPrescriptions(params?: Record<string, unknown>): Promise<PageData<Prescription>> {
    return http.get('/prescription/list', { params })
  },
  getPrescription(id: string): Promise<Prescription> {
    return http.get(`/prescription/${id}`)
  },
  createPrescription(data: { prescription: Partial<Prescription>; details: Record<string, unknown>[] }): Promise<Prescription> {
    return http.post('/prescription/create', data)
  },
  review(id: string, approved: boolean, comment?: string): Promise<void> {
    return http.post('/prescription/review', { id, approved, comment })
  },
  cancel(id: string): Promise<void> {
    return http.post(`/prescription/cancel/${id}`)
  },
}
