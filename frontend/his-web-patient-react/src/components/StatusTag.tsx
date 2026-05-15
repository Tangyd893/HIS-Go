import { Label, Icon } from 'semantic-ui-react'

const statusMap: Record<number, { text: string; color: any }> = {
  0: { text: '待处理', color: 'grey' },
  1: { text: '进行中', color: 'blue' },
  2: { text: '已完成', color: 'green' },
  3: { text: '已取消', color: 'red' },
  4: { text: '待审核', color: 'orange' },
  5: { text: '已退费', color: 'red' },
}

interface StatusTagProps {
  status: number
  labels?: Record<number, string>
}

export default function StatusTag({ status, labels }: StatusTagProps) {
  const config = labels ? { text: labels[status] || '未知', color: 'grey' } : statusMap[status] || { text: '未知', color: 'grey' }

  return (
    <Label color={config.color as any} basic>
      <Icon name="circle" color={config.color as any} />
      {config.text}
    </Label>
  )
}
