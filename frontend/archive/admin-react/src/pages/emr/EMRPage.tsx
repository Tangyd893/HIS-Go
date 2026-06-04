import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import { emrApi } from '../../api/emr'
import type { EMRRecord, PageData } from '../../api/types'

export default function EMRPage() {
  const [data, setData] = useState<PageData<EMRRecord>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<EMRRecord> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await emrApi.getRecords({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<EMRRecord>[] = [
    { key: 'patientName', header: '患者', render: (r) => r.patientName },
    { key: 'type', header: '病历类型', render: (r) => r.type },
    { key: 'doctorName', header: '医生', render: (r) => r.doctorName },
    { key: 'createdAt', header: '创建日期', render: (r) => r.createdAt },
    { key: 'updatedAt', header: '更新日期', render: (r) => r.updatedAt },
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
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => { setEditing({}); setModalOpen(true) }} createLabel="新增病历" />
      <FormModal
        open={modalOpen} title={editing.id ? '编辑病历' : '新增病历'}
        initialValues={{ patientName: editing.patientName || '', type: editing.type || '', doctorName: editing.doctorName || '', content: editing.content || '' }}
        fields={[
          { key: 'patientName', label: '患者', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'type', label: '病历类型', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} placeholder="门诊/住院/急诊" /> },
          { key: 'doctorName', label: '医生', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'content', label: '病历内容', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
        ]}
        onSubmit={async (vals) => {
          if (editing.id) await emrApi.updateRecord(editing.id, vals)
          else await emrApi.createRecord(vals)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
