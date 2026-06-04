import http from './client'
import type { ExaminationReport, PageData, PageQuery } from './types'

export const examinationApi = {
  /** 创建检查报告 */
  create(data: Partial<ExaminationReport>): Promise<ExaminationReport> {
    return http.post('/examination/report', data)
  },

  /** 报告详情 */
  getById(id: string): Promise<ExaminationReport> {
    return http.get(`/examination/report/${id}`)
  },

  /** 报告列表 */
  getList(params: PageQuery & { patientId?: string; status?: number }): Promise<PageData<ExaminationReport>> {
    return http.get('/examination/reports', { params })
  },

  /** 审核报告 */
  review(data: { report_id: string; reviewer_id: string; approved: boolean; comment?: string }): Promise<void> {
    return http.post('/examination/review', data)
  },
}
