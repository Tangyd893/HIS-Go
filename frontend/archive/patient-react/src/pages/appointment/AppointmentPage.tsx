import { Segment, Header, Form, Button, Grid, Card, Icon, Label } from 'semantic-ui-react'

export default function AppointmentPage() {
  return (
    <div>
      <Segment>
        <Header as="h3"><Icon name="calendar alternate" /> 预约挂号</Header>
        <Form>
          <Form.Group widths={3}>
            <Form.Field><label>科室</label><Form.Input placeholder="选择科室" /></Form.Field>
            <Form.Field><label>医生</label><Form.Input placeholder="选择医生" /></Form.Field>
            <Form.Field><label>日期</label><Form.Input type="date" /></Form.Field>
          </Form.Group>
          <Form.Group widths={2}>
            <Form.Field><label>时段</label><Form.Input placeholder="上午/下午" /></Form.Field>
            <Form.Field><label>患者姓名</label><Form.Input placeholder="姓名" /></Form.Field>
          </Form.Group>
          <Button primary>确认预约</Button>
        </Form>
      </Segment>

      <Header as="h4" style={{ marginTop: 24 }}>我的预约</Header>
      <Grid columns={3}>
        {[
          { dept: '内科', doctor: '张医生', date: '2026-05-16', time: '上午', status: '已预约' },
          { dept: '外科', doctor: '李医生', date: '2026-05-18', time: '下午', status: '已完成' },
        ].map((r, i) => (
          <Grid.Column key={i}>
            <Card fluid>
              <Card.Content>
                <Card.Header>{r.dept}</Card.Header>
                <Card.Meta>{r.doctor} | {r.date} {r.time}</Card.Meta>
                <Label color={r.status === '已预约' ? 'blue' : 'green'}>{r.status}</Label>
              </Card.Content>
            </Card>
          </Grid.Column>
        ))}
      </Grid>
    </div>
  )
}
