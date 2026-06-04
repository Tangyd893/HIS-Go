import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import { userApi } from '../../api/user'
import type { Employee, PageData } from '../../api/types'

export default function EmployeeList() {
  const [data, setData] = useState<PageData<Employee>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<Employee> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await userApi.getEmployees({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<Employee>[] = [
    { key: 'name', header: '姓名', render: (r) => r.name },
    { key: 'gender', header: '性别', render: (r) => r.gender },
    { key: 'phone', header: '电话', render: (r) => r.phone },
    { key: 'deptName', header: '科室', render: (r) => r.deptName },
    { key: 'title', header: '职称', render: (r) => r.title },
    { key: 'role', header: '角色', render: (r) => r.role },
    {
      key: 'actions', header: '操作', render: (r) => (
        <Button.Group size="small">
          <Button icon="edit" onClick={() => { setEditing(r); setModalOpen(true) }} />
          <Button icon="trash" color="red" onClick={async () => {
            if (confirm('确认删除？')) { await userApi.deleteEmployee(r.id); fetch() }
          }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable
        columns={columns} data={data} loading={loading} onPageChange={setPage}
        onCreate={() => { setEditing({}); setModalOpen(true) }} createLabel="新增员工"
      />
      <FormModal
        open={modalOpen}
        title={editing.id ? '编辑员工' : '新增员工'}
        initialValues={{ name: editing.name || '', gender: editing.gender || '', phone: editing.phone || '', deptName: editing.deptName || '', title: editing.title || '', role: editing.role || '' }}
        fields={[
          { key: 'name', label: '姓名', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'gender', label: '性别', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'phone', label: '电话', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'deptName', label: '科室', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'title', label: '职称', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'role', label: '角色', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
        ]}
        onSubmit={async (vals) => {
          if (editing.id) await userApi.updateEmployee(editing.id, vals)
          else await userApi.createEmployee(vals)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
