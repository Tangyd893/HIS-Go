import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import StatusTag from '../../components/StatusTag'
import { scheduleApi } from '../../api/schedule'
import type { Schedule, PageData } from '../../api/types'

export default function SchedulePage() {
  const [data, setData] = useState<PageData<Schedule>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<Schedule> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await scheduleApi.getSchedules({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<Schedule>[] = [
    { key: 'doctorName', header: '医生', render: (r) => r.doctorName },
    { key: 'deptName', header: '科室', render: (r) => r.deptName },
    { key: 'date', header: '日期', render: (r) => r.date },
    { key: 'timeSlot', header: '时段', render: (r) => r.timeSlot },
    { key: 'totalSlots', header: '总号源', render: (r) => r.totalSlots },
    { key: 'remainingSlots', header: '剩余', render: (r) => r.remainingSlots },
    { key: 'status', header: '状态', render: (r) => <StatusTag status={r.status} labels={{ 0: '未开放', 1: '可预约', 2: '已约满', 3: '已过期' }} /> },
    {
      key: 'actions', header: '操作', render: (r) => (
        <Button.Group size="small">
          <Button icon="edit" onClick={() => { setEditing(r); setModalOpen(true) }} />
          <Button icon="trash" color="red" onClick={async () => {
            if (confirm('确认删除？')) { await scheduleApi.deleteSchedule(r.id); fetch() }
          }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => { setEditing({}); setModalOpen(true) }} createLabel="新增号源" />
      <FormModal
        open={modalOpen} title={editing.id ? '编辑号源' : '新增号源'}
        initialValues={{ doctorName: editing.doctorName || '', deptName: editing.deptName || '', date: editing.date || '', timeSlot: editing.timeSlot || '', totalSlots: String(editing.totalSlots || ''), remainingSlots: String(editing.remainingSlots || '') }}
        fields={[
          { key: 'doctorName', label: '医生', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'deptName', label: '科室', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'date', label: '日期', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="date" /> },
          { key: 'timeSlot', label: '时段', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} placeholder="上午/下午" /> },
          { key: 'totalSlots', label: '总号源', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="number" /> },
          { key: 'remainingSlots', label: '剩余', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="number" /> },
        ]}
        onSubmit={async (vals) => {
          const payload = { ...vals, totalSlots: parseInt(vals.totalSlots) || 0, remainingSlots: parseInt(vals.remainingSlots) || 0 }
          if (editing.id) await scheduleApi.updateSchedule(editing.id, payload)
          else await scheduleApi.createSchedule(payload)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
