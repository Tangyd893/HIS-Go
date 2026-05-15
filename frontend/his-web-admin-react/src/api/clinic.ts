import http from './client'
import type { ClinicRecord, PageData } from './types'

export const clinicApi = {
  getRecords(params?: Record<string, unknown>): Promise<PageData<ClinicRecord>> {
    return http.get('/clinic/records', { params })
  },
  createRecord(data: Partial<ClinicRecord>): Promise<ClinicRecord> {
    return http.post('/clinic/records', data)
  },
  updateRecord(id: string, data: Partial<ClinicRecord>): Promise<ClinicRecord> {
    return http.put(`/clinic/records/${id}`, data)
  },
}
