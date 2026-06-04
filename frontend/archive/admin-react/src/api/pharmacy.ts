import http from './client'
import type { Drug, PageData } from './types'

export const pharmacyApi = {
  getDrugs(params?: Record<string, unknown>): Promise<PageData<Drug>> {
    return http.get('/pharmacy/drugs', { params })
  },
  getDrug(id: string): Promise<Drug> {
    return http.get(`/pharmacy/drug/${id}`)
  },
  addStock(id: string, qty: number): Promise<void> {
    return http.post(`/pharmacy/stock/add/${id}`, { qty })
  },
  dispense(prescriptionId: string, drugId: string, qty: number, dispenserId: string): Promise<void> {
    return http.post('/pharmacy/dispense', { prescriptionId, drugId, qty, dispenserId })
  },
  getExpired(): Promise<Drug[]> {
    return http.get('/pharmacy/expired')
  },
}
