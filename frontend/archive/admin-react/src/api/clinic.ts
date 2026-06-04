import http from './client'
import type { ClinicRecord, PageData } from './types'

export const clinicApi = {
  getRecords(params?: Record<string, unknown>): Promise<PageData<ClinicRecord>> {
    return http.get('/clinic/records', { params })
  },
  getRecord(id: string): Promise<ClinicRecord> {
    return http.get(`/clinic/record/${id}`)
  },
  createRecord(data: Partial<ClinicRecord>): Promise<ClinicRecord> {
    return http.post('/clinic/record', data)
  },
  createExaminationRequest(data: Record<string, unknown>): Promise<unknown> {
    return http.post('/clinic/examination-request', data)
  },
}
