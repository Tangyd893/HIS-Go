import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Form, Button, Message, Grid, Header, Icon } from 'semantic-ui-react'
import { useAuthStore } from '../../store/authStore'

export default function LoginPage() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const login = useAuthStore((s) => s.login)
  const navigate = useNavigate()

  const handleSubmit = async () => {
    if (!username || !password) { setError('请输入用户名和密码'); return }
    setLoading(true); setError('')
    try { await login({ username, password }); navigate('/dashboard', { replace: true }) }
    catch (err: any) { setError(err.message || '登录失败') }
    finally { setLoading(false) }
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', background: 'linear-gradient(135deg, #2185d0 0%, #36cfc9 100%)' }}>
      <div style={{ width: 380 }}>
        <Grid textAlign="center" style={{ marginBottom: 0 }}>
          <Grid.Column>
            <Icon name="heart" size="big" inverted />
            <Header as="h1" inverted style={{ marginBottom: 8 }}>患者服务中心</Header>
            <p style={{ color: 'rgba(255,255,255,0.7)', marginBottom: 32 }}>HIS-Go Hospital Information System</p>
          </Grid.Column>
        </Grid>
        <div style={{ background: '#fff', padding: '40px 32px', borderRadius: 8, boxShadow: '0 10px 40px rgba(0,0,0,0.15)' }}>
          {error && <Message negative onDismiss={() => setError('')}><Message.Header>{error}</Message.Header></Message>}
          <Form onSubmit={handleSubmit} loading={loading}>
            <Form.Field><label>用户名</label><Form.Input icon="user" iconPosition="left" placeholder="用户名" value={username} onChange={(_, { value }) => setUsername(value)} autoFocus /></Form.Field>
            <Form.Field><label>密码</label><Form.Input icon="lock" iconPosition="left" type="password" placeholder="密码" value={password} onChange={(_, { value }) => setPassword(value)} /></Form.Field>
            <Button primary fluid size="large" type="submit" loading={loading}>登 录</Button>
          </Form>
          <p style={{ marginTop: 20, textAlign: 'center', color: '#999', fontSize: 11 }}>演示账号: demo-doctor / demo123</p>
        </div>
      </div>
    </div>
  )
}
