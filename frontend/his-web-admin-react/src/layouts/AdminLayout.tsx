import { useState, useMemo, Fragment } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import {
  Sidebar, Menu, Icon, Dropdown, Breadcrumb, Segment, Container,
} from 'semantic-ui-react'
import { useAuthStore } from '../store/authStore'
import { useAppStore } from '../store/appStore'

const menuItems = [
  { key: 'dashboard', icon: 'dashboard', label: '工作台' },
  {
    key: 'user', icon: 'users', label: '用户管理', children: [
      { key: 'user/patients', label: '患者管理' },
      { key: 'user/employees', label: '员工管理' },
      { key: 'user/departments', label: '科室管理' },
    ],
  },
  { key: 'registration', icon: 'edit', label: '挂号管理' },
  { key: 'schedule', icon: 'calendar alternate', label: '排班管理' },
  { key: 'clinic', icon: 'stethoscope', label: '门诊诊疗' },
  { key: 'prescription', icon: 'file text', label: '处方管理' },
  { key: 'pharmacy', icon: 'flask', label: '药房管理' },
  { key: 'billing', icon: 'dollar sign', label: '收费结算' },
  {
    key: 'inpatient-group', icon: 'hospital', label: '住院管理', children: [
      { key: 'inpatient', label: '住院记录' },
      { key: 'emr', label: '电子病历' },
    ],
  },
  { key: 'examination', icon: 'microscope', label: '检查检验' },
  {
    key: 'outpatient-group', icon: 'globe', label: '院外服务', children: [
      { key: 'outpatient', label: '在线问诊' },
      { key: 'followup', label: '随访管理' },
      { key: 'health-record', label: '健康档案' },
    ],
  },
  { key: 'notification', icon: 'bell', label: '消息通知' },
  { key: 'statistics', icon: 'bar chart', label: '数据统计' },
  { key: 'system', icon: 'setting', label: '系统设置' },
]

const pathLabels: Record<string, string> = {
  dashboard: '工作台',
  'user/patients': '患者管理',
  'user/employees': '员工管理',
  'user/departments': '科室管理',
  registration: '挂号管理',
  schedule: '排班管理',
  clinic: '门诊诊疗',
  prescription: '处方管理',
  pharmacy: '药房管理',
  billing: '收费结算',
  examination: '检查检验',
  inpatient: '住院管理',
  emr: '电子病历',
  outpatient: '在线问诊',
  followup: '随访管理',
  'health-record': '健康档案',
  notification: '消息通知',
  statistics: '数据统计',
  system: '系统设置',
}

export default function AdminLayout() {
  const collapsed = useAppStore((s) => s.collapsed)
  const toggleCollapsed = useAppStore((s) => s.toggleCollapsed)
  const username = useAuthStore((s) => s.username)
  const logout = useAuthStore((s) => s.logout)
  const navigate = useNavigate()
  const location = useLocation()
  const [openMenus, setOpenMenus] = useState<Set<string>>(new Set(['user']))

  const activeKey = useMemo(() => {
    const path = location.pathname.replace('/', '') || 'dashboard'
    return path
  }, [location.pathname])

  const breadcrumbs = useMemo(() => {
    return activeKey.split('/').map((part) => pathLabels[part] || part)
  }, [activeKey])

  const toggleMenu = (key: string) => {
    setOpenMenus((prev) => {
      const next = new Set(prev)
      if (next.has(key)) next.delete(key)
      else next.add(key)
      return next
    })
  }

  const onMenuClick = (key: string) => {
    navigate(`/${key}`)
  }

  const sidebarWidth = collapsed ? 60 : 220

  return (
    <Sidebar.Pushable as={Segment} style={{ minHeight: '100vh', margin: 0, border: 'none', borderRadius: 0 }}>
      <Sidebar
        as={Menu}
        animation="overlay"
        icon={collapsed ? 'labeled' : undefined}
        inverted
        vertical
        visible
        width="thin"
        style={{ width: sidebarWidth, overflowY: 'auto', transition: 'width 0.2s' }}
      >
        <Menu.Item header style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 54 }}>
          <Icon name="hospital" size="large" />
          {!collapsed && <span style={{ marginLeft: 10 }}>HIS-Go</span>}
        </Menu.Item>

        {menuItems.map((item) =>
          item.children ? (
            <Menu.Item key={item.key}>
              <Menu.Header onClick={() => toggleMenu(item.key)} style={{ cursor: 'pointer' }}>
                <Icon name={item.icon as any} />
                {!collapsed && item.label}
              </Menu.Header>
              {!collapsed && openMenus.has(item.key) && (
                <Menu.Menu>
                  {item.children.map((child) => (
                    <Menu.Item
                      key={child.key}
                      active={activeKey === child.key}
                      onClick={() => onMenuClick(child.key)}
                    >
                      {child.label}
                    </Menu.Item>
                  ))}
                </Menu.Menu>
              )}
            </Menu.Item>
          ) : (
            <Menu.Item
              key={item.key}
              active={activeKey === item.key}
              onClick={() => onMenuClick(item.key)}
            >
              <Icon name={item.icon as any} />
              {!collapsed && item.label}
            </Menu.Item>
          ),
        )}
      </Sidebar>

      <Sidebar.Pusher style={{ marginLeft: 0, minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
        <Menu attached="top" borderless style={{ margin: 0, borderRadius: 0, padding: '0 16px' }}>
          <Menu.Item onClick={toggleCollapsed}>
            <Icon name={collapsed ? 'content' : 'sidebar'} />
          </Menu.Item>
          <Menu.Menu position="right">
            <Dropdown item text={username || '用户'} icon="user circle">
              <Dropdown.Menu>
                <Dropdown.Item icon="user" text="个人信息" />
                <Dropdown.Divider />
                <Dropdown.Item icon="sign-out" text="退出登录" onClick={logout} />
              </Dropdown.Menu>
            </Dropdown>
          </Menu.Menu>
        </Menu>

        <Breadcrumb style={{ padding: '12px 20px', background: '#f8f8f8', borderBottom: '1px solid #e0e0e0' }}>
          <Breadcrumb.Section>
            <Icon name="home" />
          </Breadcrumb.Section>
          {breadcrumbs.map((label, i) => (
            <Fragment key={i}>
              <Breadcrumb.Divider icon="right angle" />
              <Breadcrumb.Section active={i === breadcrumbs.length - 1}>
                {label}
              </Breadcrumb.Section>
            </Fragment>
          ))}
        </Breadcrumb>

        <Container fluid style={{ flex: 1, padding: 24, background: '#f5f5f5' }}>
          <Outlet />
        </Container>
      </Sidebar.Pusher>
    </Sidebar.Pushable>
  )
}
