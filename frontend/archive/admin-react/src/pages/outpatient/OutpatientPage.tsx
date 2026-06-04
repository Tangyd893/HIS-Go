import { Segment, Header, Icon, Grid, Card, Button, Label } from 'semantic-ui-react'

export default function OutpatientPage() {
  return (
    <div>
      <Header as="h3"><Icon name="globe" /> 在线问诊</Header>

      <Grid columns={3}>
        {[
          { name: '张三', dept: '内科', time: '2026-05-15 09:30', status: '待接诊' },
          { name: '李四', dept: '外科', time: '2026-05-15 10:00', status: '问诊中' },
          { name: '王五', dept: '儿科', time: '2026-05-15 10:30', status: '已完成' },
        ].map((item) => (
          <Grid.Column key={item.name}>
            <Card fluid>
              <Card.Content>
                <Card.Header>{item.name}</Card.Header>
                <Card.Meta>{item.dept} | {item.time}</Card.Meta>
                <Label color={item.status === '已完成' ? 'green' : item.status === '问诊中' ? 'blue' : 'grey'}>
                  {item.status}
                </Label>
              </Card.Content>
              <Card.Content extra>
                <Button basic fluid color="blue">进入问诊</Button>
              </Card.Content>
            </Card>
          </Grid.Column>
        ))}
      </Grid>
    </div>
  )
}
