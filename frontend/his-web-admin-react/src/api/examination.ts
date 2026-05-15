import http from './client'
import type { ExaminationReport, PageData } from './types'

export const examinationApi = {
  getReports(params?: Record<string, unknown>): Promise<PageData<ExaminationReport>> {
    return http.get('/examination/reports', { params })
  },
  createReport(data: Partial<ExaminationReport>): Promise<ExaminationReport> {
    return http.post('/examination/reports', data)
  },
  updateReport(id: string, data: Partial<ExaminationReport>): Promise<ExaminationReport> {
    return http.put(`/examination/reports/${id}`, data)
  },
  reviewReport(id: string, result: string): Promise<ExaminationReport> {
    return http.put(`/examination/reports/${id}/review`, { result })
  },
}
