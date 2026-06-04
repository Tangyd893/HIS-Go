import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import StatusTag from '../../components/StatusTag'
import { registrationApi } from '../../api/registration'
import type { Schedule } from '../../api/types'

export default function RegistrationPage() {
  const [schedules, setSchedules] = useState<Schedule[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [selectedSchedule, setSelectedSchedule] = useState<Schedule | null>(null)

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setSchedules(await registrationApi.getSchedules({})) } catch {}
    finally { setLoading(false) }
  }, [])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<Schedule>[] = [
    { key: 'doctorName', header: '医生', render: (r) => r.doctorName },
    { key: 'deptName', header: '科室', render: (r) => r.deptName },
    { key: 'date', header: '日期', render: (r) => r.date },
    { key: 'timeSlot', header: '时段', render: (r) => r.timeSlot },
    { key: 'remainingSlots', header: '剩余号源', render: (r) => r.remainingSlots },
    { key: 'fee', header: '挂号费', render: (r) => `¥${r.fee}` },
    { key: 'status', header: '状态', render: (r) => <StatusTag status={r.status} labels={{ 0: '停诊', 1: '可预约' }} /> },
    {
      key: 'actions', header: '操作', render: (r) => (
        <Button.Group size="small">
          {r.status === 1 && r.remainingSlots > 0 && (
            <Button icon="user plus" color="green" onClick={() => { setSelectedSchedule(r); setModalOpen(true) }} />
          )}
        </Button.Group>
      ),
    },
  ]

  // 将 schedules 数组转为 DataTable 期望的 PageData 格式
  const pageData = { list: schedules, total: schedules.length, page: 1, pageSize: schedules.length }

  return (
    <>
      <DataTable columns={columns} data={pageData} loading={loading} onPageChange={() => {}} />
      <FormModal
        open={modalOpen} title="挂号"
        initialValues={{ patientName: '王小明', patientId: 'patient_001' }}
        fields={[
          { key: 'patientName', label: '患者姓名', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'patientId', label: '患者ID', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
        ]}
        onSubmit={async (vals) => {
          if (selectedSchedule) {
            await registrationApi.register({ patientId: vals.patientId, patientName: vals.patientName, scheduleId: selectedSchedule.id })
            fetch()
          }
        }}
        onClose={() => { setModalOpen(false); setSelectedSchedule(null) }}
      />
    </>
  )
}
