import { useState, useEffect, useCallback } from 'react'
import { Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import http from '../../api/client'
import type { Prescription, PageData } from '../../api/types'

export default function PrescriptionPage() {
  const [data, setData] = useState<PageData<Prescription>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await http.get('/prescription/list', { params: { page, pageSize: 10 } })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<Prescription>[] = [
    { key: 'doctorName', header: '医生', render: (r) => r.doctorName },
    { key: 'createdAt', header: '开方日期', render: (r) => r.createdAt },
    { key: 'status', header: '状态', render: (r) => r.status === 2 ? '已发药' : r.status === 1 ? '已审核' : '待审核' },
    { key: 'details', header: '药品数', render: (r) => r.details?.length || 0 },
  ]

  return <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} />
}
