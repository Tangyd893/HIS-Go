import { useState, useEffect, useCallback } from 'react'
import { Input, Button, TextArea } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import { clinicApi } from '../../api/clinic'
import type { ClinicRecord, PageData } from '../../api/types'

export default function ClinicPage() {
  const [data, setData] = useState<PageData<ClinicRecord>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<ClinicRecord> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await clinicApi.getRecords({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<ClinicRecord>[] = [
    { key: 'patientName', header: '患者', render: (r) => r.patientName },
    { key: 'doctorName', header: '医生', render: (r) => r.doctorName },
    { key: 'chiefComplaint', header: '主诉', render: (r) => r.chiefComplaint },
    { key: 'diagnosis', header: '诊断', render: (r) => r.diagnosis },
    { key: 'createdAt', header: '就诊时间', render: (r) => r.createdAt },
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
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => { setEditing({}); setModalOpen(true) }} createLabel="新增门诊记录" />
      <FormModal
        open={modalOpen} title={editing.id ? '编辑门诊记录' : '新增门诊记录'}
        initialValues={{ patientName: editing.patientName || '', doctorName: editing.doctorName || '', chiefComplaint: editing.chiefComplaint || '', diagnosis: editing.diagnosis || '' }}
        fields={[
          { key: 'patientName', label: '患者', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'doctorName', label: '医生', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'chiefComplaint', label: '主诉', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'diagnosis', label: '诊断', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
        ]}
        onSubmit={async (vals) => {
          if (editing.id) await clinicApi.updateRecord(editing.id, vals)
          else await clinicApi.createRecord(vals)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
