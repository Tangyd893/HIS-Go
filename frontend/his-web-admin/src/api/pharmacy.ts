import http from './client'
import type { Drug, PageData, PageQuery } from './types'

export const pharmacyApi = {
  /** 药品列表 */
  getDrugs(params: PageQuery & { name?: string }): Promise<PageData<Drug>> {
    return http.get('/pharmacy/drugs', { params })
  },

  /** 药品详情 */
  getDrugById(id: string): Promise<Drug> {
    return http.get(`/pharmacy/drug/${id}`)
  },

  /** 入库 */
  addStock(id: string, quantity: number): Promise<void> {
    return http.post(`/pharmacy/stock/add/${id}`, { quantity })
  },

  /** 发药 */
  dispense(data: { prescription_id: string; drug_id: string; quantity: number; dispenser_id: string }): Promise<void> {
    return http.post('/pharmacy/dispense', data)
  },

  /** 扫描过期药品 */
  getExpired(): Promise<Drug[]> {
    return http.get('/pharmacy/expired')
  },
}
