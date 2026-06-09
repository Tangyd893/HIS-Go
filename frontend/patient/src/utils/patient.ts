/** 演示环境：登录账号与 his_user.patients 表 ID 映射 */
const DEMO_PATIENT_MAP: Record<string, string> = {
  'demo-patient': 'patient_001',
  'patient_001': 'patient_001',
  'patient_002': 'patient_002',
  'patient_003': 'patient_003',
  'patient_004': 'patient_004',
  'patient_005': 'patient_005',
  'patient_006': 'patient_006',
  'patient_007': 'patient_007',
}

/** 获取当前登录用户对应的患者档案 ID */
export function resolvePatientId(userInfo: { userId?: string; username?: string } | null): string {
  const userId = userInfo?.userId || userInfo?.username || ''
  // 先查演示映射表，再回退到 userId 本身（适用于 userId 即为 patientId 的场景）
  return DEMO_PATIENT_MAP[userId] || userId
}

const TIME_SLOT_LABELS: Record<number, string> = {
  1: '上午',
  2: '下午',
  3: '晚上',
}

/** 格式化排班时段 */
export function formatTimeSlot(slot: number | string): string {
  const n = typeof slot === 'string' ? parseInt(slot, 10) : slot
  return TIME_SLOT_LABELS[n] || String(slot)
}

/**  flatten 科室树为下拉选项 */
export function flattenDepartments(depts: { id: string; name: string; children?: typeof depts }[]): { id: string; name: string }[] {
  const result: { id: string; name: string }[] = []
  for (const d of depts || []) {
    if (d.id && d.name) result.push({ id: d.id, name: d.name })
    if (d.children?.length) result.push(...flattenDepartments(d.children))
  }
  return result
}

const EVENT_TYPE_LABELS: Record<string, string> = {
  visit: '就诊',
  prescription: '处方',
  examination: '检查',
  followup: '随访',
}

/** 时间轴事件类型中文 */
export function formatEventType(type: string): string {
  return EVENT_TYPE_LABELS[type] || type
}

/** 规范化号源列表字段（兼容后端 camelCase） */
export function normalizeSchedule(item: Record<string, unknown>) {
  const remain = (item.remainCount ?? item.remainingSlots ?? 0) as number
  const slot = item.timeSlot ?? item.time_slot ?? ''
  return {
    ...item,
    remainCount: remain,
    remainingSlots: remain,
    totalCount: item.totalCount ?? item.total_count,
    timeSlotLabel: formatTimeSlot(slot as number),
  }
}
