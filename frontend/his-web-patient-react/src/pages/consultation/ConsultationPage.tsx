import { Segment, Header, Button, Card, Grid, Icon, Label, Form } from 'semantic-ui-react'

export default function ConsultationPage() {
  return (
    <div>
      <Segment>
        <Header as="h3"><Icon name="comment" /> 在线问诊</Header>
        <Form>
          <Form.Field><label>症状描述</label><Form.TextArea placeholder="请描述您的症状..." rows={4} /></Form.Field>
          <Form.Group widths={2}>
            <Form.Field><label>科室</label><Form.Input placeholder="选择科室" /></Form.Field>
            <Form.Field><label>医生</label><Form.Input placeholder="选择医生" /></Form.Field>
          </Form.Group>
          <Button primary>发起问诊</Button>
        </Form>
      </Segment>

      <Header as="h4" style={{ marginTop: 24 }}>问诊记录</Header>
      <Grid columns={2}>
        {[
          { doctor: '张医生', dept: '内科', date: '2026-05-14', status: '已完成', msg: '建议多休息，按时服药' },
          { doctor: '王医生', dept: '儿科', date: '2026-05-13', status: '已完成', msg: '建议调整饮食结构' },
        ].map((r, i) => (
          <Grid.Column key={i}>
            <Card fluid>
              <Card.Content>
                <Card.Header>{r.doctor} <Label size="mini" color="green">{r.status}</Label></Card.Header>
                <Card.Meta>{r.dept} | {r.date}</Card.Meta>
                <Card.Description>{r.msg}</Card.Description>
              </Card.Content>
            </Card>
          </Grid.Column>
        ))}
      </Grid>
    </div>
  )
}
