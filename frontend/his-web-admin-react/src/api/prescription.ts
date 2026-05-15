import http from './client'
import type { Prescription, PageData } from './types'

export const prescriptionApi = {
  getPrescriptions(params?: Record<string, unknown>): Promise<PageData<Prescription>> {
    return http.get('/prescription/list', { params })
  },
  createPrescription(data: Partial<Prescription>): Promise<Prescription> {
    return http.post('/prescription', data)
  },
  updatePrescription(id: string, data: Partial<Prescription>): Promise<Prescription> {
    return http.put(`/prescription/${id}`, data)
  },
  deletePrescription(id: string): Promise<void> {
    return http.delete(`/prescription/${id}`)
  },
}
