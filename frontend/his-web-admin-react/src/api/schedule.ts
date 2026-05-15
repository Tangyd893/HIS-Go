import http from './client'
import type { Schedule, ScheduleRule, PageData } from './types'

export const scheduleApi = {
  getSchedules(params?: Record<string, unknown>): Promise<PageData<Schedule>> {
    return http.get('/schedule/list', { params })
  },
  createSchedule(data: Partial<Schedule>): Promise<Schedule> {
    return http.post('/schedule', data)
  },
  updateSchedule(id: string, data: Partial<Schedule>): Promise<Schedule> {
    return http.put(`/schedule/${id}`, data)
  },
  deleteSchedule(id: string): Promise<void> {
    return http.delete(`/schedule/${id}`)
  },
  getRules(params?: Record<string, unknown>): Promise<PageData<ScheduleRule>> {
    return http.get('/schedule/rules', { params })
  },
}
