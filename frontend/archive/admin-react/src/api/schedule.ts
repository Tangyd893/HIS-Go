import http from './client'
import type { Schedule, PageData } from './types'

export const scheduleApi = {
  getSchedules(params?: Record<string, unknown>): Promise<PageData<Schedule>> {
    return http.get('/schedule/list', { params })
  },
  /** 自动生成下一周排班 */
  generateWeekly(): Promise<void> {
    return http.post('/schedule/generate')
  },
  /** 更新排班 */
  updateSchedule(data: Partial<Schedule>): Promise<Schedule> {
    return http.put('/schedule/update', data)
  },
  /** 取消排班 */
  cancelSchedule(id: string): Promise<void> {
    return http.post(`/schedule/cancel/${id}`)
  },
}
