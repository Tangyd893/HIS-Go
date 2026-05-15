import { Table, Menu, Icon, Button, Label } from 'semantic-ui-react'
import type { PageData } from '../api/types'
import type { ReactNode } from 'react'

export interface Column<T> {
  key: string
  header: string
  render: (record: T) => ReactNode
  width?: number
}

interface DataTableProps<T> {
  columns: Column<T>[]
  data: PageData<T> | null | undefined
  loading: boolean
  onPageChange: (page: number) => void
  onCreate?: () => void
  createLabel?: string
}

export default function DataTable<T extends { id: string }>({
  columns,
  data,
  loading,
  onPageChange,
  onCreate,
  createLabel = '新增',
}: DataTableProps<T>) {
  const list = data?.list ?? []
  const total = data?.total ?? 0
  const page = data?.page ?? 1
  const totalPages = data?.pageSize ? Math.ceil(total / data.pageSize) : 1

  const pageNumbers: number[] = []
  const maxVisible = 5
  let start = Math.max(1, page - Math.floor(maxVisible / 2))
  const end = Math.min(totalPages, start + maxVisible - 1)
  if (end - start + 1 < maxVisible) {
    start = Math.max(1, end - maxVisible + 1)
  }
  for (let i = start; i <= end; i++) pageNumbers.push(i)

  return (
    <div>
      {onCreate && (
        <Button primary icon onClick={onCreate} style={{ marginBottom: 16 }}>
          <Icon name="plus" /> {createLabel}
        </Button>
      )}

      <Table celled striped selectable loading={loading}>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell width={1}>#</Table.HeaderCell>
            {columns.map((col) => (
              <Table.HeaderCell key={col.key} width={col.width}>
                {col.header}
              </Table.HeaderCell>
            ))}
          </Table.Row>
        </Table.Header>

        <Table.Body>
          {list.length === 0 && !loading ? (
            <Table.Row>
              <Table.Cell colSpan={columns.length + 1} textAlign="center">
                暂无数据
              </Table.Cell>
            </Table.Row>
          ) : (
            list.map((record, index) => (
              <Table.Row key={record.id}>
                <Table.Cell>{(page - 1) * (data?.pageSize ?? 20) + index + 1}</Table.Cell>
                {columns.map((col) => (
                  <Table.Cell key={col.key}>{col.render(record)}</Table.Cell>
                ))}
              </Table.Row>
            ))
          )}
        </Table.Body>

        {totalPages > 1 && (
          <Table.Footer>
            <Table.Row>
              <Table.HeaderCell colSpan={columns.length + 1}>
                <Menu floated="right" pagination>
                  <Menu.Item as="a" icon disabled={page <= 1} onClick={() => onPageChange(page - 1)}>
                    <Icon name="chevron left" />
                  </Menu.Item>
                  {pageNumbers.map((p) => (
                    <Menu.Item as="a" key={p} active={p === page} onClick={() => onPageChange(p)}>
                      {p}
                    </Menu.Item>
                  ))}
                  <Menu.Item
                    as="a"
                    icon
                    disabled={page >= totalPages}
                    onClick={() => onPageChange(page + 1)}
                  >
                    <Icon name="chevron right" />
                  </Menu.Item>
                </Menu>
                <Label style={{ float: 'right', marginTop: 8 }}>
                  共 {total} 条
                </Label>
              </Table.HeaderCell>
            </Table.Row>
          </Table.Footer>
        )}
      </Table>
    </div>
  )
}
