import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import StatusTag from '../../components/StatusTag'
import { registrationApi } from '../../api/registration'
import type { Registration, PageData } from '../../api/types'

export default function RegistrationPage() {
  const [data, setData] = useState<PageData<Registration>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<Registration> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await registrationApi.getRegistrations({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<Registration>[] = [
    { key: 'patientName', header: '患者', render: (r) => r.patientName },
    { key: 'deptName', header: '科室', render: (r) => r.deptName },
    { key: 'doctorName', header: '医生', render: (r) => r.doctorName },
    { key: 'registrationDate', header: '挂号日期', render: (r) => r.registrationDate },
    { key: 'queueNumber', header: '排队号', render: (r) => r.queueNumber },
    { key: 'status', header: '状态', render: (r) => <StatusTag status={r.status} labels={{ 0: '待就诊', 1: '就诊中', 2: '已完成', 3: '已取消' }} /> },
    {
      key: 'actions', header: '操作', render: (r) => (
        <Button.Group size="small">
          <Button icon="edit" onClick={() => { setEditing(r); setModalOpen(true) }} />
          <Button icon="trash" color="red" onClick={async () => {
            if (confirm('确认删除？')) { await registrationApi.deleteRegistration(r.id); fetch() }
          }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => { setEditing({}); setModalOpen(true) }} createLabel="新增挂号" />
      <FormModal
        open={modalOpen} title={editing.id ? '编辑挂号' : '新增挂号'}
        initialValues={{ patientName: editing.patientName || '', deptName: editing.deptName || '', doctorName: editing.doctorName || '', registrationDate: editing.registrationDate || '', queueNumber: String(editing.queueNumber || '') }}
        fields={[
          { key: 'patientName', label: '患者', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'deptName', label: '科室', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'doctorName', label: '医生', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'registrationDate', label: '挂号日期', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="date" /> },
          { key: 'queueNumber', label: '排队号', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="number" /> },
        ]}
        onSubmit={async (vals) => {
          const payload = { ...vals, queueNumber: parseInt(vals.queueNumber) || 0 }
          if (editing.id) await registrationApi.updateRegistration(editing.id, payload)
          else await registrationApi.createRegistration(payload)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
