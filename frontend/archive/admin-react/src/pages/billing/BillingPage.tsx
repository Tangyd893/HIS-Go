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
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => setModalOpen(true)} createLabel="新增账单" />
      <FormModal
        open={modalOpen} title="新增账单"
        initialValues={{ patientName: '', totalAmount: '' }}
        fields={[
          { key: 'patientName', label: '患者', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'totalAmount', label: '金额', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} type="number" /> },
        ]}
        onSubmit={async (vals) => {
          await billingApi.createBill({ patientName: vals.patientName, totalAmount: parseFloat(vals.totalAmount) || 0 })
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
