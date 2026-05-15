import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import { userApi } from '../../api/user'
import type { Department, PageData } from '../../api/types'

export default function DepartmentList() {
  const [data, setData] = useState<PageData<Department>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<Department> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await userApi.getDepartments({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<Department>[] = [
    { key: 'name', header: '科室名称', render: (r) => r.name },
    { key: 'parentId', header: '上级科室', render: (r) => r.parentId || '-' },
    { key: 'description', header: '描述', render: (r) => r.description },
    {
      key: 'actions', header: '操作', render: (r) => (
        <Button.Group size="small">
          <Button icon="edit" onClick={() => { setEditing(r); setModalOpen(true) }} />
          <Button icon="trash" color="red" onClick={async () => {
            if (confirm('确认删除？')) { await userApi.deleteDepartment(r.id); fetch() }
          }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable
        columns={columns} data={data} loading={loading} onPageChange={setPage}
        onCreate={() => { setEditing({}); setModalOpen(true) }} createLabel="新增科室"
      />
      <FormModal
        open={modalOpen}
        title={editing.id ? '编辑科室' : '新增科室'}
        initialValues={{ name: editing.name || '', parentId: editing.parentId || '', description: editing.description || '' }}
        fields={[
          { key: 'name', label: '科室名称', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'parentId', label: '上级科室ID', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'description', label: '描述', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
        ]}
        onSubmit={async (vals) => {
          if (editing.id) await userApi.updateDepartment(editing.id, vals)
          else await userApi.createDepartment(vals)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
