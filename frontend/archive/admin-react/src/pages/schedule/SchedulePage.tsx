import { useState, useEffect, useCallback } from 'react'
import { Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import StatusTag from '../../components/StatusTag'
import { scheduleApi } from '../../api/schedule'
import type { Schedule, PageData } from '../../api/types'

export default function SchedulePage() {
  const [data, setData] = useState<PageData<Schedule>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)

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
          <Button icon="trash" color="red" onClick={async () => {
            if (confirm('确认取消？')) { await scheduleApi.cancelSchedule(r.id); fetch() }
          }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <DataTable
      columns={columns} data={data} loading={loading} onPageChange={setPage}
      onCreate={async () => { await scheduleApi.generateWeekly(); fetch() }}
      createLabel="生成下一周排班"
    />
  )
}
