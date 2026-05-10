import http from './client'
import type { Schedule, PageQuery } from './types'

export const scheduleApi = {
  /** 生成排班 */
  generate(data: { startDate: string; endDate: string; deptId: string }): Promise<void> {
    return http.post('/schedule/generate', data)
  },

  /** 排班列表 */
  getList(params: { deptId?: string; date?: string; doctorId?: string }): Promise<Schedule[]> {
    return http.get('/schedule/list', { params })
  },

  /** 更新排班 */
  update(data: Partial<Schedule>): Promise<void> {
    return http.put('/schedule/update', data)
  },

  /** 取消排班 */
  cancel(id: string): Promise<void> {
    return http.post(`/schedule/cancel/${id}`)
  },
}
