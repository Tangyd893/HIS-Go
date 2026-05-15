import { useNavigate } from 'react-router-dom'
import { Card, Statistic, Grid, Button, Icon, Segment } from 'semantic-ui-react'
import { useAuthStore } from '../../store/authStore'

export default function DashboardPage() {
  const navigate = useNavigate()
  const username = useAuthStore((s) => s.username)

  return (
    <div>
      <Statistic.Group widths={3} size="small">
        <Card fluid><Statistic color="blue"><Statistic.Value><Icon name="calendar check" /> 3</Statistic.Value><Statistic.Label>我的预约</Statistic.Label></Statistic></Card>
        <Card fluid><Statistic color="green"><Statistic.Value><Icon name="heart" /> 2</Statistic.Value><Statistic.Label>慢病管理</Statistic.Label></Statistic></Card>
        <Card fluid><Statistic color="orange"><Statistic.Value><Icon name="phone" /> 1</Statistic.Value><Statistic.Label>待随访</Statistic.Label></Statistic></Card>
      </Statistic.Group>

      <Segment style={{ marginTop: 24 }}>
        <h3><Icon name="lightning" /> 快捷服务</h3>
        <Button.Group widths={4}>
          <Button color="teal" onClick={() => navigate('/appointment')}><Icon name="calendar alternate" /> 预约挂号</Button>
          <Button color="blue" onClick={() => navigate('/consultation')}><Icon name="comment" /> 在线问诊</Button>
          <Button color="green" onClick={() => navigate('/prescription')}><Icon name="file text" /> 我的处方</Button>
          <Button color="purple" onClick={() => navigate('/report')}><Icon name="file alternate" /> 检查报告</Button>
        </Button.Group>
      </Segment>

      <Segment style={{ marginTop: 16 }}>
        <h3><Icon name="info circle" /> 个人信息</h3>
        <Grid columns={2}>
          <Grid.Column><p><strong>当前用户：</strong>{username}</p></Grid.Column>
          <Grid.Column><p><strong>系统版本：</strong>HIS-Go v1.0.0</p></Grid.Column>
        </Grid>
      </Segment>
    </div>
  )
}
