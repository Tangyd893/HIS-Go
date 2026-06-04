import { useState, useEffect, useCallback } from 'react'
import { Input, Button, Icon } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import StatusTag from '../../components/StatusTag'
import { userApi } from '../../api/user'
import type { Patient, PageData } from '../../api/types'

export default function PatientList() {
  const [data, setData] = useState<PageData<Patient>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<Patient> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try {
      const res = await userApi.getPatients({ page, pageSize: 10 })
      setData(res)
    } catch { /* error shown in client interceptor */ }
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<Patient>[] = [
    { key: 'name', header: '姓名', render: (r) => r.name },
    { key: 'gender', header: '性别', render: (r) => r.gender },
    { key: 'age', header: '年龄', render: (r) => r.age },
    { key: 'phone', header: '电话', render: (r) => r.phone },
    { key: 'idCard', header: '身份证', render: (r) => r.idCard },
    { key: 'address', header: '地址', render: (r) => r.address },
    {
      key: 'actions', header: '操作', render: (r) => (
        <Button.Group size="small">
          <Button icon="edit" onClick={() => { setEditing(r); setModalOpen(true) }} />
          <Button icon="trash" color="red" onClick={async () => {
            if (confirm('确认删除？')) { await userApi.deletePatient(r.id); fetch() }
          }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable
        columns={columns}
        data={data}
        loading={loading}
        onPageChange={setPage}
        onCreate={() => { setEditing({}); setModalOpen(true) }}
        createLabel="新增患者"
      />
      <FormModal
        open={modalOpen}
        title={editing.id ? '编辑患者' : '新增患者'}
        initialValues={{ name: editing.name || '', gender: editing.gender || '', age: String(editing.age || ''), phone: editing.phone || '', idCard: editing.idCard || '', address: editing.address || '' }}
        fields={[
          { key: 'name', label: '姓名', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} placeholder="姓名" /> },
          { key: 'gender', label: '性别', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} placeholder="男/女" /> },
          { key: 'age', label: '年龄', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="number" placeholder="年龄" /> },
          { key: 'phone', label: '电话', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} placeholder="电话" /> },
          { key: 'idCard', label: '身份证', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} placeholder="身份证" /> },
          { key: 'address', label: '地址', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} placeholder="地址" /> },
        ]}
        onSubmit={async (vals) => {
          const payload = { ...vals, age: parseInt(vals.age) || 0 }
          if (editing.id) await userApi.updatePatient(editing.id, payload)
          else await userApi.createPatient(payload)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
