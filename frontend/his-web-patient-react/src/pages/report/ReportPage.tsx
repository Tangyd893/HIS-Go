import { useState, useEffect, useCallback } from 'react'
import DataTable, { type Column } from '../../components/DataTable'
import http from '../../api/client'
import type { ExaminationReport, PageData } from '../../api/types'
import StatusTag from '../../components/StatusTag'

export default function ReportPage() {
  const [data, setData] = useState<PageData<ExaminationReport>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await http.get('/examination/reports', { params: { page, pageSize: 10 } })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<ExaminationReport>[] = [
    { key: 'examType', header: '检查类型', render: (r) => r.examType },
    { key: 'examItem', header: '检查项目', render: (r) => r.examItem },
    { key: 'result', header: '结果', render: (r) => r.result || '待审核' },
    { key: 'createdAt', header: '日期', render: (r) => r.createdAt },
    { key: 'status', header: '状态', render: (r) => <StatusTag status={r.status} labels={{ 0: '已开具', 1: '已检查', 2: '已审核' }} /> },
  ]

  return <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} />
}
