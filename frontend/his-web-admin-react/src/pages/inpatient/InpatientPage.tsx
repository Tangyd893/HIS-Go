import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import StatusTag from '../../components/StatusTag'
import { inpatientApi } from '../../api/inpatient'
import type { InpatientRecord, PageData } from '../../api/types'

export default function InpatientPage() {
  const [data, setData] = useState<PageData<InpatientRecord>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<InpatientRecord> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await inpatientApi.getRecords({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<InpatientRecord>[] = [
    { key: 'patientName', header: '患者', render: (r) => r.patientName },
    { key: 'deptName', header: '科室', render: (r) => r.deptName },
    { key: 'bedNo', header: '床位', render: (r) => r.bedNo },
    { key: 'admitDate', header: '入院日期', render: (r) => r.admitDate },
    { key: 'dischargeDate', header: '出院日期', render: (r) => r.dischargeDate || '-' },
    { key: 'status', header: '状态', render: (r) => <StatusTag status={r.status} labels={{ 0: '在院', 1: '已出院' }} /> },
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
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => { setEditing({}); setModalOpen(true) }} createLabel="新增住院记录" />
      <FormModal
        open={modalOpen} title={editing.id ? '编辑住院记录' : '新增住院记录'}
        initialValues={{ patientName: editing.patientName || '', deptName: editing.deptName || '', bedNo: editing.bedNo || '', admitDate: editing.admitDate || '', dischargeDate: editing.dischargeDate || '' }}
        fields={[
          { key: 'patientName', label: '患者', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'deptName', label: '科室', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'bedNo', label: '床位', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'admitDate', label: '入院日期', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="date" /> },
          { key: 'dischargeDate', label: '出院日期', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="date" /> },
        ]}
        onSubmit={async (vals) => {
          if (editing.id) await inpatientApi.updateRecord(editing.id, vals)
          else await inpatientApi.createRecord(vals)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
