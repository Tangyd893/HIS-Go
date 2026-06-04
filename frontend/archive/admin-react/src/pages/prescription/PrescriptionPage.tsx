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
          <Button icon="check" color="green" onClick={async () => { await prescriptionApi.review(r.id, true); fetch() }} />
          <Button icon="trash" color="red" onClick={async () => {
            if (confirm('确认取消？')) { await prescriptionApi.cancel(r.id); fetch() }
          }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => setModalOpen(true)} createLabel="新增处方" />
      <FormModal
        open={modalOpen} title="新增处方"
        initialValues={{ patientName: '', doctorName: '', drugName: '', quantity: '1' }}
        fields={[
          { key: 'patientName', label: '患者', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'doctorName', label: '医生', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'drugName', label: '药品', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'quantity', label: '数量', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="number" /> },
        ]}
        onSubmit={async (vals) => {
          await prescriptionApi.createPrescription({
            prescription: { patientName: vals.patientName, doctorName: vals.doctorName, prescriptionType: 1 },
            details: [{ drugName: vals.drugName, quantity: parseInt(vals.quantity) || 1 }],
          })
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
