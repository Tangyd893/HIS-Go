import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import { pharmacyApi } from '../../api/pharmacy'
import type { Drug, PageData } from '../../api/types'
import dayjs from 'dayjs'

export default function PharmacyPage() {
  const [data, setData] = useState<PageData<Drug>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<Drug> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await pharmacyApi.getDrugs({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<Drug>[] = [
    { key: 'name', header: '药品名称', render: (r) => r.name },
    { key: 'specification', header: '规格', render: (r) => r.specification },
    { key: 'manufacturer', header: '厂家', render: (r) => r.manufacturer },
    { key: 'stock', header: '库存', render: (r) => r.stock },
    { key: 'price', header: '单价', render: (r) => `¥${r.price}` },
    { key: 'expiryDate', header: '有效期', render: (r) => r.expiryDate },
    {
      key: 'actions', header: '操作', render: (r) => (
        <Button.Group size="small">
          <Button icon="edit" onClick={() => { setEditing(r); setModalOpen(true) }} />
          <Button icon="trash" color="red" onClick={async () => {
            if (confirm('确认删除？')) { await pharmacyApi.deleteDrug(r.id); fetch() }
          }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => { setEditing({}); setModalOpen(true) }} createLabel="新增药品" />
      <FormModal
        open={modalOpen} title={editing.id ? '编辑药品' : '新增药品'}
        initialValues={{ name: editing.name || '', specification: editing.specification || '', manufacturer: editing.manufacturer || '', stock: String(editing.stock || ''), price: String(editing.price || ''), expiryDate: editing.expiryDate || '' }}
        fields={[
          { key: 'name', label: '药品名称', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'specification', label: '规格', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'manufacturer', label: '厂家', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'stock', label: '库存', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="number" /> },
          { key: 'price', label: '单价', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="number" /> },
          { key: 'expiryDate', label: '有效期', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="date" /> },
        ]}
        onSubmit={async (vals) => {
          const payload = { ...vals, stock: parseInt(vals.stock) || 0, price: parseFloat(vals.price) || 0 }
          if (editing.id) await pharmacyApi.updateDrug(editing.id, payload)
          else await pharmacyApi.createDrug(payload)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
