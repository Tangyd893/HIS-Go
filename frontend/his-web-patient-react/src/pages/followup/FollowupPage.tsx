import { useState, useEffect, useCallback } from 'react'
import DataTable, { type Column } from '../../components/DataTable'
import StatusTag from '../../components/StatusTag'
import http from '../../api/client'
import type { FollowupRecord, PageData } from '../../api/types'

export default function FollowupPage() {
  const [data, setData] = useState<PageData<FollowupRecord>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await http.get('/followup/list', { params: { page, pageSize: 10 } })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<FollowupRecord>[] = [
    { key: 'planDate', header: '计划日期', render: (r) => r.planDate },
    { key: 'actualDate', header: '实际日期', render: (r) => r.actualDate || '-' },
    { key: 'content', header: '随访内容', render: (r) => r.content },
    { key: 'status', header: '状态', render: (r) => <StatusTag status={r.status} labels={{ 0: '待随访', 1: '已随访' }} /> },
  ]

  return <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} />
}
