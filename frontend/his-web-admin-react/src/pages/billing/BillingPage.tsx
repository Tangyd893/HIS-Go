import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import StatusTag from '../../components/StatusTag'
import { billingApi } from '../../api/billing'
import type { Bill, PageData } from '../../api/types'

export default function BillingPage() {
  const [data, setData] = useState<PageData<Bill>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<Bill> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await billingApi.getBills({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<Bill>[] = [
    { key: 'billNo', header: '账单号', render: (r) => r.billNo },
    { key: 'patientName', header: '患者', render: (r) => r.patientName },
    { key: 'totalAmount', header: '金额', render: (r) => `¥${r.totalAmount}` },
    { key: 'paidAmount', header: '已付', render: (r) => `¥${r.paidAmount}` },
    { key: 'createdAt', header: '日期', render: (r) => r.createdAt },
    { key: 'status', header: '状态', render: (r) => <StatusTag status={r.status} labels={{ 0: '未付款', 1: '已付款', 2: '已退费' }} /> },
    {
      key: 'actions', header: '操作', render: (r) => (
        <Button.Group size="small">
          {r.status === 0 && <Button icon="dollar sign" color="green" onClick={async () => { await billingApi.payBill(r.id); fetch() }} />}
          <Button icon="edit" onClick={() => { setEditing(r); setModalOpen(true) }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => { setEditing({}); setModalOpen(true) }} createLabel="新增账单" />
      <FormModal
        open={modalOpen} title={editing.id ? '编辑账单' : '新增账单'}
        initialValues={{ billNo: editing.billNo || '', patientName: editing.patientName || '', totalAmount: String(editing.totalAmount || ''), paidAmount: String(editing.paidAmount || '') }}
        fields={[
          { key: 'billNo', label: '账单号', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'patientName', label: '患者', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'totalAmount', label: '金额', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="number" /> },
          { key: 'paidAmount', label: '已付', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="number" /> },
        ]}
        onSubmit={async (vals) => {
          const payload = { ...vals, totalAmount: parseFloat(vals.totalAmount) || 0, paidAmount: parseFloat(vals.paidAmount) || 0 }
          if (editing.id) await billingApi.updateBill(editing.id, payload)
          else await billingApi.createBill(payload)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
