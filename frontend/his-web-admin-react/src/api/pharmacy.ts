import http from './client'
import type { Drug, PageData } from './types'

export const pharmacyApi = {
  getDrugs(params?: Record<string, unknown>): Promise<PageData<Drug>> {
    return http.get('/pharmacy/drugs', { params })
  },
  createDrug(data: Partial<Drug>): Promise<Drug> {
    return http.post('/pharmacy/drugs', data)
  },
  updateDrug(id: string, data: Partial<Drug>): Promise<Drug> {
    return http.put(`/pharmacy/drugs/${id}`, data)
  },
  deleteDrug(id: string): Promise<void> {
    return http.delete(`/pharmacy/drugs/${id}`)
  },
}
