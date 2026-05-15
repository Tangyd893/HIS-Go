import { Segment, Header, Grid, Card, Icon, Label } from 'semantic-ui-react'

export default function HealthRecordPage() {
  return (
    <div>
      <Header as="h3"><Icon name="heart" /> 我的健康档案</Header>
      <Grid columns={2}>
        {[
          { type: '体检报告', date: '2026-04-15', status: '正常', desc: '各项指标正常' },
          { type: '血常规', date: '2026-05-01', status: '正常', desc: '白细胞计数正常' },
          { type: '心电图', date: '2026-03-20', status: '异常', desc: '轻微心律不齐，建议复查' },
          { type: 'B超', date: '2026-02-10', status: '正常', desc: '未见异常' },
        ].map((r, i) => (
          <Grid.Column key={i}>
            <Card fluid>
              <Card.Content>
                <Card.Header>{r.type} <Label size="mini" color={r.status === '正常' ? 'green' : 'orange'}>{r.status}</Label></Card.Header>
                <Card.Meta>{r.date}</Card.Meta>
                <Card.Description>{r.desc}</Card.Description>
              </Card.Content>
            </Card>
          </Grid.Column>
        ))}
      </Grid>
    </div>
  )
}
