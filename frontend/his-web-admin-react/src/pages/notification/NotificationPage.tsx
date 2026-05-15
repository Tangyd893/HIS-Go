import { useState, useEffect, useCallback } from 'react'
import { Button } from 'semantic-ui-react'
import DataTable, { type Column } from '../../components/DataTable'
import { othersApi } from '../../api/others'
import type { Notification, PageData } from '../../api/types'

export default function NotificationPage() {
  const [data, setData] = useState<PageData<Notification>>()
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)

  const fetch = useCallback(async () => {
    setLoading(true)
    try { setData(await othersApi.getNotifications({ page, pageSize: 10 })) } catch {}
    finally { setLoading(false) }
  }, [page])

  useEffect(() => { fetch() }, [fetch])

  const columns: Column<Notification>[] = [
    { key: 'title', header: '标题', render: (r) => r.title },
    { key: 'content', header: '内容', render: (r) => r.content },
    { key: 'type', header: '类型', render: (r) => r.type },
    { key: 'createdAt', header: '时间', render: (r) => r.createdAt },
    { key: 'read', header: '状态', render: (r) => r.read ? '已读' : '未读' },
    {
      key: 'actions', header: '操作', render: (r) => !r.read ? (
        <Button size="small" basic color="blue" onClick={async () => { await othersApi.markNotificationRead(r.id); fetch() }}>
          标为已读
        </Button>
      ) : null,
    },
  ]

  return (
    <DataTable columns={columns} data={data} loading={loading} onPageChange={setPage} />
  )
}
