import { useState, useEffect, useCallback } from 'react'
import { Input, Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import StatusTag from '../../components/StatusTag'
import { examinationApi } from '../../api/examination'
import type { ExaminationReport, PageData } from '../../api/types'

export default function ExaminationPage() {
  const [data, setData] = useState<PageData<ExaminationReport>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Partial<ExaminationReport> & { id?: string }>({})

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await examinationApi.getReports({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<ExaminationReport>[] = [
    { key: 'patientName', header: '患者', render: (r) => r.patientName },
    { key: 'examType', header: '检查类型', render: (r) => r.examType },
    { key: 'examItem', header: '检查项目', render: (r) => r.examItem },
    { key: 'result', header: '结果', render: (r) => r.result || '待审核' },
    { key: 'createdAt', header: '日期', render: (r) => r.createdAt },
    { key: 'status', header: '状态', render: (r) => <StatusTag status={r.status} labels={{ 0: '已开具', 1: '已检查', 2: '已审核' }} /> },
    {
      key: 'actions', header: '操作', render: (r) => (
        <Button.Group size="small">
          <Button icon="edit" onClick={() => { setEditing(r); setModalOpen(true) }} />
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => { setEditing({}); setModalOpen(true) }} createLabel="新增检查" />
      <FormModal
        open={modalOpen} title={editing.id ? '编辑检查' : '新增检查'}
        initialValues={{ patientName: editing.patientName || '', examType: editing.examType || '', examItem: editing.examItem || '', result: editing.result || '' }}
        fields={[
          { key: 'patientName', label: '患者', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'examType', label: '检查类型', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'examItem', label: '检查项目', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'result', label: '结果', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
        ]}
        onSubmit={async (vals) => {
          if (editing.id) await examinationApi.updateReport(editing.id, vals)
          else await examinationApi.createReport(vals)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
