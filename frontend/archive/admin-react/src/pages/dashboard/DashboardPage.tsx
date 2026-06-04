import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, Statistic, Grid, Button, Icon, Label, Segment } from 'semantic-ui-react'
import { useAuthStore } from '../../store/authStore'
import http from '../../api/client'
import { scheduleApi } from '../../api/schedule'
import { clinicApi } from '../../api/clinic'
import { billingApi } from '../../api/billing'

const quickActions = [
  { label: '排班管理', icon: 'calendar', path: '/schedule', color: 'teal' },
  { label: '预约挂号', icon: 'edit', path: '/registration', color: 'blue' },
  { label: '门诊接诊', icon: 'stethoscope', path: '/clinic', color: 'green' },
  { label: '开具处方', icon: 'file text', path: '/prescription', color: 'orange' },
  { label: '收费结算', icon: 'dollar sign', path: '/billing', color: 'red' },
  { label: '药房发药', icon: 'flask', path: '/pharmacy', color: 'purple' },
]

/** Profile A 已启用的服务列表 */
const profileAServices = ['Gateway', 'Auth', 'User', 'Registration', 'Clinic', 'Prescription', 'Billing', 'Pharmacy', 'Schedule', 'System']

export default function DashboardPage() {
  const navigate = useNavigate()
  const username = useAuthStore((s) => s.username)
  const role = useAuthStore((s) => s.role)

  const [todaySchedules, setTodaySchedules] = useState<number | null>(null)
  const [todayClinics, setTodayClinics] = useState<number | null>(null)
  const [todayRevenue, setTodayRevenue] = useState<number | null>(null)
  const [gwHealthy, setGwHealthy] = useState(false)

  useEffect(() => {
    // 今日号源数
    scheduleApi.getSchedules({ page: 1, pageSize: 1 }).then(r => setTodaySchedules(r?.total ?? 0)).catch(() => setTodaySchedules(0))
    // 今日门诊记录数
    clinicApi.getRecords({ page: 1, pageSize: 1 }).then(r => setTodayClinics(r?.total ?? 0)).catch(() => setTodayClinics(0))
    // 今日收入（从 billing 列表累计）
    billingApi.getBills({ page: 1, pageSize: 100 }).then(r => {
      const total = (r?.list ?? []).reduce((sum, b) => sum + (b.paidAmount || 0), 0)
      setTodayRevenue(total)
    }).catch(() => setTodayRevenue(0))
    // Gateway 健康检查
    http.get('/health').then(() => setGwHealthy(true)).catch(() => setGwHealthy(false))
  }, [])

  const stats = [
    { label: '排班号源', value: todaySchedules ?? '...', icon: 'calendar', color: '#2185d0' },
    { label: '门诊记录', value: todayClinics ?? '...', icon: 'stethoscope', color: '#21ba45' },
    { label: '今日收入', value: todayRevenue != null ? `¥${todayRevenue.toLocaleString()}` : '...', icon: 'dollar sign', color: '#db2828' },
    { label: '网关状态', value: gwHealthy ? '正常' : '...', icon: 'server', color: gwHealthy ? '#21ba45' : '#fbbd08' },
  ]

  return (
    <div>
      <Statistic.Group widths={4} size="small">
        {stats.map((card) => (
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
            <h3><Icon name="lightning" /> 快捷操作</h3>
            <Button.Group widths={6}>
              {quickActions.map((action) => (
                <Button key={action.label} color={action.color as any} onClick={() => navigate(action.path)}>
                  <Icon name={action.icon as any} /> {action.label}
                </Button>
              ))}
            </Button.Group>
          </Segment>
        </Grid.Column>
        <Grid.Column>
          <Segment>
            <h3><Icon name="info circle" /> 系统信息</h3>
            <Grid columns={2}>
              <Grid.Column>
                <p><strong>当前用户：</strong>{username}</p>
                <p><strong>角色：</strong>{role}</p>
                <p><strong>网关状态：</strong>{gwHealthy ? '🟢 正常' : '🔴 离线'}</p>
              </Grid.Column>
              <Grid.Column>
                <p><strong>演示 Profile：</strong>A — 管理端</p>
                <p><strong>部署模式：</strong>轻量化演示</p>
              </Grid.Column>
            </Grid>
          </Segment>
        </Grid.Column>
      </Grid>

      <Segment style={{ marginTop: 24 }}>
        <h3><Icon name="server" /> 服务模块状态（Profile A）</h3>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
          {profileAServices.map((svc) => (
            <Label key={svc} color="green" basic style={{ marginBottom: 4 }}>
              <Icon name="check circle" /> {svc}
            </Label>
          ))}
          <Label color="grey" basic style={{ marginBottom: 4 }}>
            <Icon name="minus circle" /> 其他服务（演示未启用）
          </Label>
        </div>
      </Segment>
    </div>
  )
}
