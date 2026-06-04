import { useMemo } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Sidebar, Menu, Icon, Dropdown, Breadcrumb, Segment, Container } from 'semantic-ui-react'
import { useAuthStore } from '../store/authStore'
import { useAppStore } from '../store/appStore'

const menuItems = [
  { key: 'dashboard', icon: 'home', label: '首页' },
  { key: 'appointment', icon: 'calendar alternate', label: '预约挂号' },
  { key: 'consultation', icon: 'comment', label: '在线问诊' },
  { key: 'prescription', icon: 'file text', label: '我的处方' },
  { key: 'report', icon: 'file alternate', label: '检查报告' },
  { key: 'health-record', icon: 'heart', label: '健康档案' },
  { key: 'chronic', icon: 'shield alternate', label: '慢病管理' },
  { key: 'followup', icon: 'phone', label: '我的随访' },
]

const pathLabels: Record<string, string> = {
  dashboard: '首页', appointment: '预约挂号', consultation: '在线问诊',
  prescription: '我的处方', report: '检查报告', 'health-record': '健康档案',
  chronic: '慢病管理', followup: '我的随访',
}

export default function PatientLayout() {
  const collapsed = useAppStore((s) => s.collapsed)
  const toggleCollapsed = useAppStore((s) => s.toggleCollapsed)
  const username = useAuthStore((s) => s.username)
  const logout = useAuthStore((s) => s.logout)
  const navigate = useNavigate()
  const location = useLocation()

  const activeKey = useMemo(() => location.pathname.replace('/', '') || 'dashboard', [location.pathname])
  const breadcrumbs = useMemo(() => [pathLabels[activeKey] || activeKey], [activeKey])

  const sidebarWidth = collapsed ? 60 : 200

  return (
    <Sidebar.Pushable as={Segment} style={{ minHeight: '100vh', margin: 0, border: 'none', borderRadius: 0 }}>
      <Sidebar
        as={Menu} animation="overlay" icon={collapsed ? 'labeled' : undefined}
        inverted vertical visible width="thin"
        style={{ width: sidebarWidth, overflowY: 'auto', transition: 'width 0.2s' }}
      >
        <Menu.Item header style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 54 }}>
          <Icon name="heart" color="red" size="large" />
          {!collapsed && <span style={{ marginLeft: 10 }}>患者中心</span>}
        </Menu.Item>
        {menuItems.map((item) => (
          <Menu.Item key={item.key} active={activeKey === item.key} onClick={() => navigate(`/${item.key}`)}>
            <Icon name={item.icon as any} />
            {!collapsed && item.label}
          </Menu.Item>
        ))}
      </Sidebar>

      <Sidebar.Pusher style={{ marginLeft: 0, minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
        <Menu attached="top" borderless style={{ margin: 0, borderRadius: 0, padding: '0 16px' }}>
          <Menu.Item onClick={toggleCollapsed}>
            <Icon name={collapsed ? 'content' : 'sidebar'} />
          </Menu.Item>
          <Menu.Menu position="right">
            <Dropdown item text={username || '用户'} icon="user circle">
              <Dropdown.Menu>
                <Dropdown.Item icon="sign-out" text="退出登录" onClick={logout} />
              </Dropdown.Menu>
            </Dropdown>
          </Menu.Menu>
        </Menu>

        <Breadcrumb style={{ padding: '12px 20px', background: '#f8f8f8', borderBottom: '1px solid #e0e0e0' }}>
          <Breadcrumb.Section><Icon name="home" /></Breadcrumb.Section>
          <Breadcrumb.Divider icon="right angle" />
          <Breadcrumb.Section active>{breadcrumbs[0]}</Breadcrumb.Section>
        </Breadcrumb>

        <Container fluid style={{ flex: 1, padding: 24, background: '#f5f5f5' }}>
          <Outlet />
        </Container>
      </Sidebar.Pusher>
    </Sidebar.Pushable>
  )
}
