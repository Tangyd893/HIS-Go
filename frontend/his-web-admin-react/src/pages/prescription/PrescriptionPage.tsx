import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import StatusTag from '../../components/StatusTag'
import { prescriptionApi } from '../../api/prescription'
import type { Prescription, PageData } from '../../api/types'

export default function PrescriptionPage() {
  const [data, setData] = useState<PageData<Prescription>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<Prescription> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await prescriptionApi.getPrescriptions({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<Prescription>[] = [
    { key: 'patientName', header: '患者', render: (r) => r.patientName },
    { key: 'doctorName', header: '医生', render: (r) => r.doctorName },
    { key: 'createdAt', header: '开方日期', render: (r) => r.createdAt },
    { key: 'status', header: '状态', render: (r) => <StatusTag status={r.status} labels={{ 0: '待审核', 1: '已审核', 2: '已发药', 3: '已退药' }} /> },
    {
      key: 'actions', header: '操作', render: (r) => (
        <Button.Group size="small">
          <Button icon="edit" onClick={() => { setEditing(r); setModalOpen(true) }} />
          <Button icon="trash" color="red" onClick={async () => {
            if (confirm('确认删除？')) { await prescriptionApi.deletePrescription(r.id); fetch() }
          }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => { setEditing({}); setModalOpen(true) }} createLabel="新增处方" />
      <FormModal
        open={modalOpen} title={editing.id ? '编辑处方' : '新增处方'}
        initialValues={{ patientName: editing.patientName || '', doctorName: editing.doctorName || '' }}
        fields={[
          { key: 'patientName', label: '患者', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'doctorName', label: '医生', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
        ]}
        onSubmit={async (vals) => {
          if (editing.id) await prescriptionApi.updatePrescription(editing.id, vals)
          else await prescriptionApi.createPrescription(vals)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
