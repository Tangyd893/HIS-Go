import http from './client'
import type { Registration, PageData } from './types'

export const registrationApi = {
  getRegistrations(params?: Record<string, unknown>): Promise<PageData<Registration>> {
    return http.get('/registration/list', { params })
  },
  createRegistration(data: Partial<Registration>): Promise<Registration> {
    return http.post('/registration', data)
  },
  updateRegistration(id: string, data: Partial<Registration>): Promise<Registration> {
    return http.put(`/registration/${id}`, data)
  },
  deleteRegistration(id: string): Promise<void> {
    return http.delete(`/registration/${id}`)
  },
}
