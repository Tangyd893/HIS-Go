import { useState, useEffect, useCallback } from 'react'
import DataTable, { type Column } from '../../components/DataTable'
import { othersApi } from '../../api/others'
import type { HealthRecord, PageData } from '../../api/types'

export default function HealthRecordPage() {
  const [data, setData] = useState<PageData<HealthRecord>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await othersApi.getHealthRecords({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<HealthRecord>[] = [
    { key: 'patientName', header: '患者', render: (r) => r.patientName },
    { key: 'recordType', header: '记录类型', render: (r) => r.recordType },
    { key: 'content', header: '内容', render: (r) => r.content },
    { key: 'createdAt', header: '创建日期', render: (r) => r.createdAt },
  ]

  return (
    <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} />
  )
}
