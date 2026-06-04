import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import { pharmacyApi } from '../../api/pharmacy'
import type { Drug, PageData } from '../../api/types'

export default function PharmacyPage() {
  const [data, setData] = useState<PageData<Drug>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [selectedDrugId, setSelectedDrugId] = useState('')

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
          <Button icon="plus" color="green" onClick={() => { setSelectedDrugId(r.id); setModalOpen(true) }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} />
      <FormModal
        open={modalOpen} title="补充库存"
        initialValues={{ qty: '10' }}
        fields={[
          { key: 'qty', label: '入库数量', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="number" /> },
        ]}
        onSubmit={async (vals) => {
          if (selectedDrugId) {
            await pharmacyApi.addStock(selectedDrugId, parseInt(vals.qty) || 0)
            fetch()
          }
        }}
        onClose={() => { setModalOpen(false); setSelectedDrugId('') }}
      />
    </>
  )
}
