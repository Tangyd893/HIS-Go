import { useState, useEffect, useCallback } from 'react'
import { Input } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import FormModal from '../../components/FormModal'
import { clinicApi } from '../../api/clinic'
import type { ClinicRecord, PageData } from '../../api/types'

export default function ClinicPage() {
  const [data, setData] = useState<PageData<ClinicRecord>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await clinicApi.getRecords({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<ClinicRecord>[] = [
    { key: 'patientName', header: '患者', render: (r) => r.patientName },
    { key: 'doctorName', header: '医生', render: (r) => r.doctorName },
    { key: 'chiefComplaint', header: '主诉', render: (r) => r.chiefComplaint },
    { key: 'diagnosis', header: '诊断', render: (r) => r.diagnosis },
    { key: 'createdAt', header: '就诊时间', render: (r) => r.createdAt },
  ]

  return (
    <>
      <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} onCreate={() => setModalOpen(true)} createLabel="新增门诊记录" />
      <FormModal
        open={modalOpen} title="新增门诊记录"
        initialValues={{ patientName: '', doctorName: '', chiefComplaint: '', diagnosis: '' }}
        fields={[
          { key: 'patientName', label: '患者', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'doctorName', label: '医生', required: true, render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'chiefComplaint', label: '主诉', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
          { key: 'diagnosis', label: '诊断', render: (v, onChange) => <Input value={v} onChange={(_, d) => onChange(d.value)} /> },
        ]}
        onSubmit={async (vals) => {
          await clinicApi.createRecord(vals)
          fetch()
        }}
        onClose={() => setModalOpen(false)}
      />
    </>
  )
}
