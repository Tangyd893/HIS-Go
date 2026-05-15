import http from './client'
import type { FollowupRecord, HealthRecord, Notification, PageData } from './types'

export const othersApi = {
  getFollowups(params?: Record<string, unknown>): Promise<PageData<FollowupRecord>> {
    return http.get('/followup/list', { params })
  },
  createFollowup(data: Partial<FollowupRecord>): Promise<FollowupRecord> {
    return http.post('/followup', data)
  },
  updateFollowup(id: string, data: Partial<FollowupRecord>): Promise<FollowupRecord> {
    return http.put(`/followup/${id}`, data)
  },

  getHealthRecords(params?: Record<string, unknown>): Promise<PageData<HealthRecord>> {
    return http.get('/health-record/list', { params })
  },

  getNotifications(params?: Record<string, unknown>): Promise<PageData<Notification>> {
    return http.get('/notification/list', { params })
  },
  markNotificationRead(id: string): Promise<void> {
    return http.put(`/notification/${id}/read`)
  },
}
