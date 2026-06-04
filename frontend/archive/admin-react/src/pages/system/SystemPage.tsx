import { Segment, Header, Form, Button, Icon, Grid } from 'semantic-ui-react'
import { useState } from 'react'

export default function SystemPage() {
  const [siteName, setSiteName] = useState('HIS-Go 医院管理系统')

  return (
    <div>
      <Segment>
        <Header as="h3"><Icon name="setting" /> 基本设置</Header>
        <Form>
          <Form.Field>
            <label>系统名称</label>
            <Form.Input value={siteName} onChange={(_, { value }) => setSiteName(value)} />
          </Form.Field>
          <Form.Field>
            <label>系统描述</label>
            <Form.Input defaultValue="医院信息管理系统 v1.0.0" />
          </Form.Field>
          <Button primary onClick={() => alert('设置已保存')}>保存设置</Button>
        </Form>
      </Segment>

      <Segment style={{ marginTop: 16 }}>
        <Header as="h3"><Icon name="users" /> 角色管理</Header>
        <Grid columns={2}>
          {['超级管理员', '医生', '护士', '药师', '收费员', '管理员'].map((role) => (
            <Grid.Column key={role}>
              <Segment>
                <Icon name="user" /> {role}
                <Button floated="right" size="mini" basic>编辑</Button>
              </Segment>
            </Grid.Column>
          ))}
        </Grid>
      </Segment>

      <Segment style={{ marginTop: 16 }}>
        <Header as="h3"><Icon name="database" /> 系统信息</Header>
        <p><strong>后端框架：</strong>Go + Gin + gRPC</p>
        <p><strong>前端框架：</strong>React 19 + Semantic UI React</p>
        <p><strong>数据库：</strong>PostgreSQL 17</p>
        <p><strong>构建工具：</strong>Vite 6</p>
      </Segment>
    </div>
  )
}
