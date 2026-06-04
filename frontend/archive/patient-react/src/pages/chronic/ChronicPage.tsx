import { Segment, Header, Grid, Card, Icon, Button, Label, Progress } from 'semantic-ui-react'

export default function ChronicPage() {
  return (
    <div>
      <Header as="h3"><Icon name="shield alternate" /> 慢病管理</Header>
      <Grid columns={2}>
        {[
          { name: '高血压', date: '2025-01', status: '管理中', progress: 75, desc: '坚持服药，定期测血压' },
          { name: '糖尿病', date: '2025-03', status: '管理中', progress: 60, desc: '控制饮食，监测血糖' },
        ].map((r, i) => (
          <Grid.Column key={i}>
            <Card fluid>
              <Card.Content>
                <Card.Header>{r.name} <Label size="mini" color="blue">{r.status}</Label></Card.Header>
                <Card.Meta>确诊日期: {r.date}</Card.Meta>
                <Progress percent={r.progress} color="green" />
                <Card.Description>{r.desc}</Card.Description>
              </Card.Content>
              <Card.Content extra>
                <Button basic fluid color="blue"><Icon name="edit" /> 更新数据</Button>
              </Card.Content>
            </Card>
          </Grid.Column>
        ))}
      </Grid>
    </div>
  )
}
