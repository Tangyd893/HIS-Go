import { useNavigate } from 'react-router-dom'
import { Card, Statistic, Grid, Button, Icon, Label, Segment } from 'semantic-ui-react'
import { useAuthStore } from '../../store/authStore'

const statCards = [
  { label: '今日挂号', value: 128, icon: 'edit', color: '#2185d0' },
  { label: '今日门诊', value: 96, icon: 'stethoscope', color: '#21ba45' },
  { label: '住院患者', value: 45, icon: 'hospital', color: '#fbbd08' },
  { label: '今日收入', value: '¥12,580', icon: 'dollar sign', color: '#db2828' },
]

const quickActions = [
  { label: '预约挂号', icon: 'edit', path: '/registration', color: 'teal' },
  { label: '门诊接诊', icon: 'stethoscope', path: '/clinic', color: 'blue' },
  { label: '开具处方', icon: 'file text', path: '/prescription', color: 'green' },
  { label: '收费结算', icon: 'dollar sign', path: '/billing', color: 'red' },
  { label: '药房发药', icon: 'flask', path: '/pharmacy', color: 'purple' },
  { label: '住院登记', icon: 'hospital', path: '/inpatient', color: 'orange' },
]

const services = [
  'Gateway', 'Auth', 'User', 'Registration', 'Clinic', 'Prescription',
  'Billing', 'Pharmacy', 'Examination', 'Inpatient', 'Schedule',
  'Outpatient', 'Followup', 'HealthRecord', 'Notification', 'Statistics', 'System', 'EMR',
]

export default function DashboardPage() {
  const navigate = useNavigate()
  const username = useAuthStore((s) => s.username)
  const role = useAuthStore((s) => s.role)

  return (
    <div>
      <Statistic.Group widths={4} size="small">
        {statCards.map((card) => (
          <Card key={card.label} fluid>
            <Statistic color={card.color as any}>
              <Statistic.Value>
                <Icon name={card.icon as any} size="small" style={{ marginRight: 8 }} />
                {card.value}
              </Statistic.Value>
              <Statistic.Label>{card.label}</Statistic.Label>
            </Statistic>
          </Card>
        ))}
      </Statistic.Group>

      <Grid columns={2} style={{ marginTop: 24 }}>
        <Grid.Column>
          <Segment>
            <h3>
              <Icon name="lightning" /> 快捷操作
            </h3>
            <Button.Group widths={6}>
              {quickActions.map((action) => (
                <Button
                  key={action.label}
                  color={action.color as any}
                  onClick={() => navigate(action.path)}
                >
                  <Icon name={action.icon as any} /> {action.label}
                </Button>
              ))}
            </Button.Group>
          </Segment>
        </Grid.Column>
        <Grid.Column>
          <Segment>
            <h3>
              <Icon name="info circle" /> 系统信息
            </h3>
            <Grid columns={2}>
              <Grid.Column>
                <p><strong>系统版本：</strong>HIS-Go v1.0.0</p>
                <p><strong>当前用户：</strong>{username}</p>
                <p><strong>角色：</strong>{role}</p>
              </Grid.Column>
              <Grid.Column>
                <p><strong>后端服务：</strong>Go + Gin + gRPC</p>
                <p><strong>前端框架：</strong>React + Semantic UI</p>
                <p><strong>数据库：</strong>PostgreSQL 17</p>
              </Grid.Column>
            </Grid>
          </Segment>
        </Grid.Column>
      </Grid>

      <Segment style={{ marginTop: 24 }}>
        <h3>
          <Icon name="server" /> 服务模块状态
        </h3>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
          {services.map((svc) => (
            <Label key={svc} color="green" basic style={{ marginBottom: 4 }}>
              <Icon name="check circle" /> {svc}
            </Label>
          ))}
        </div>
      </Segment>
    </div>
  )
}
