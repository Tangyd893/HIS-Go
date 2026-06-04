import http from './client'
import type { Registration, Schedule, PageData } from './types'

export const registrationApi = {
  /** 查询号源列表 */
  getSchedules(params?: Record<string, unknown>): Promise<Schedule[]> {
    return http.get('/registration/schedules', { params })
  },
  /** 挂号 */
  register(data: { patientId: string; patientName?: string; scheduleId: string }): Promise<Registration> {
    return http.post('/registration/register', data)
  },
  /** 取消挂号 */
  cancel(id: string): Promise<void> {
    return http.post(`/registration/cancel/${id}`)
  },
  /** 签到 */
  signIn(id: string): Promise<void> {
    return http.post(`/registration/signin/${id}`)
  },
  /** 查询排队状态 */
  getQueueStatus(registrationId: string): Promise<{ registrationId: string; rank: number }> {
    return http.get('/registration/queue-status', { params: { registrationId } })
  },
}
