import http from './client'
import type { EMRRecord, PageData } from './types'

export const emrApi = {
  getRecords(params?: Record<string, unknown>): Promise<PageData<EMRRecord>> {
    return http.get('/emr/records', { params })
  },
  createRecord(data: Partial<EMRRecord>): Promise<EMRRecord> {
    return http.post('/emr/records', data)
  },
  updateRecord(id: string, data: Partial<EMRRecord>): Promise<EMRRecord> {
    return http.put(`/emr/records/${id}`, data)
  },
}
