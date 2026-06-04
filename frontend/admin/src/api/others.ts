import http from './client'

export const outpatientApi = {
  /** 创建问诊 */
  createConsultation(data: any): Promise<any> {
    return http.post('/outpatient/consultation', data)
  },

  /** 问诊详情 */
  getConsultation(id: string): Promise<any> {
    return http.get(`/outpatient/consultation/${id}`)
  },

  /** 问诊列表 */
  getConsultations(params: any): Promise<any> {
    return http.get('/outpatient/consultations', { params })
  },

  /** 发送消息 */
  sendMessage(data: any): Promise<void> {
    return http.post('/outpatient/message', data)
  },

  /** 消息列表 */
  getMessages(consultationId: string): Promise<any[]> {
    return http.get('/outpatient/messages', { params: { consultationId } })
  },
}

export const followupApi = {
  /** 随访计划列表 */
  getPlans(params: any): Promise<any> {
    return http.get('/followup/plans', { params })
  },

  /** 执行随访任务 */
  executeTask(data: { task_id: string; result: string }): Promise<void> {
    return http.post('/followup/task/execute', data)
  },

  /** 满意度调查 */
  submitSurvey(data: any): Promise<void> {
    return http.post('/followup/survey', data)
  },
}

export const healthRecordApi = {
  /** 档案摘要 */
  getSummary(patientId: string): Promise<any> {
    return http.get(`/health-record/summary/${patientId}`)
  },

  /** 时间轴 */
  getTimeline(patientId: string): Promise<any[]> {
    return http.get(`/health-record/timeline/${patientId}`)
  },
}

export const notificationApi = {
  /** 发送通知 */
  send(data: any): Promise<void> {
    return http.post('/notification/send', data)
  },

  /** 批量发送 */
  batchSend(data: any): Promise<void> {
    return http.post('/notification/batch-send', data)
  },

  /** 模板列表 */
  getTemplates(): Promise<any[]> {
    return http.get('/notification/templates')
  },
}

export const statisticsApi = {
  /** 运营统计 */
  getOperation(params: { start_date: string; end_date: string }): Promise<any> {
    return http.post('/statistics/operation', params)
  },

  /** 科室工作量 */
  getDeptWorkload(params: { start_date: string; end_date: string }): Promise<any> {
    return http.post('/statistics/dept-workload', params)
  },

  /** 收入趋势 */
  getRevenueTrend(params: { start_date: string; end_date: string }): Promise<any> {
    return http.post('/statistics/revenue-trend', params)
  },
}

export const systemApi = {
  /** 字典类型 */
  getDictTypes(): Promise<any[]> {
    return http.get('/system/dict-types')
  },

  /** 字典项 */
  getDictItems(dictType: string): Promise<any[]> {
    return http.get('/system/dict-items', { params: { dictType } })
  },

  /** 创建字典项 */
  createDictItem(data: any): Promise<void> {
    return http.post('/system/dict-item', data)
  },

  /** 系统参数 */
  getParams(): Promise<any[]> {
    return http.get('/system/params')
  },

  /** 更新参数 */
  updateParam(data: any): Promise<void> {
    return http.put('/system/param', data)
  },

  /** 操作日志 */
  getOperationLogs(params: any): Promise<any> {
    return http.get('/system/operation-logs', { params })
  },
}
