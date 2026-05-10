import http from './client'
import type { Registration, Schedule, PageData, PageQuery } from './types'

export const registrationApi = {
  /** 查询号源 */
  getSchedules(params: { deptId?: string; date?: string }): Promise<Schedule[]> {
    return http.get('/registration/schedules', { params })
  },

  /** 挂号 */
  register(data: { patientId: string; patientName: string; scheduleId: string }): Promise<Registration> {
    return http.post('/registration/register', data)
  },

  /** 取消挂号 */
  cancel(id: string): Promise<void> {
    return http.post(`/registration/cancel/${id}`)
  },

  /** 签到 */
  signin(id: string): Promise<void> {
    return http.post(`/registration/signin/${id}`)
  },

  /** 排队状态 */
  getQueueStatus(registrationId: string): Promise<any> {
    return http.get('/registration/queue-status', { params: { registrationId } })
  },
}
