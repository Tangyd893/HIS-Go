import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import StatusTag from '../../components/StatusTag'
import { othersApi } from '../../api/others'
import type { FollowupRecord, PageData } from '../../api/types'

export default function FollowupPage() {
  const [data, setData] = useState<PageData<FollowupRecord>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<FollowupRecord> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await othersApi.getFollowups({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<FollowupRecord>[] = [
    { key: 'patientName', header: '患者', render: (r) => r.patientName },
    { key: 'planDate', header: '计划日期', render: (r) => r.planDate },
    { key: 'actualDate', header: '实际日期', render: (r) => r.actualDate || '-' },
    { key: 'content', header: '随访内容', render: (r) => r.content },
    { key: 'status', header: '状态', render: (r) => <StatusTag status={r.status} labels={{ 0: '待随访', 1: '已随访', 2: '已过期' }} /> },
    {
      key: 'actions', header: '操作', render: (r) => (
        <Button.Group size="small">
          <Button icon="edit" onClick={() => { setEditing(r); setModalOpen(true) }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => { setEditing({}); setModalOpen(true) }} createLabel="新增随访" />
      <FormModal
        open={modalOpen} title={editing.id ? '编辑随访' : '新增随访'}
        initialValues={{ patientName: editing.patientName || '', planDate: editing.planDate || '', actualDate: editing.actualDate || '', content: editing.content || '' }}
        fields={[
          { key: 'patientName', label: '患者', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'planDate', label: '计划日期', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="date" /> },
          { key: 'actualDate', label: '实际日期', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="date" /> },
          { key: 'content', label: '内容', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
        ]}
        onSubmit={async (vals) => {
          if (editing.id) await othersApi.updateFollowup(editing.id, vals)
          else await othersApi.createFollowup(vals)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
