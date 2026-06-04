import { Card, Statistic, Grid, Icon, Segment, Label } from 'semantic-ui-react'

export default function StatisticsPage() {
  return (
    <div>
      <Statistic.Group widths={4} size="small">
        <Card fluid><Statistic color="blue"><Statistic.Value><Icon name="calendar check" /> 1,280</Statistic.Value><Statistic.Label>本月挂号</Statistic.Label></Statistic></Card>
        <Card fluid><Statistic color="green"><Statistic.Value><Icon name="stethoscope" /> 960</Statistic.Value><Statistic.Label>本月门诊</Statistic.Label></Statistic></Card>
        <Card fluid><Statistic color="orange"><Statistic.Value><Icon name="hospital" /> 320</Statistic.Value><Statistic.Label>住院人次</Statistic.Label></Statistic></Card>
        <Card fluid><Statistic color="red"><Statistic.Value><Icon name="dollar sign" /> ¥158,200</Statistic.Value><Statistic.Label>本月收入</Statistic.Label></Statistic></Card>
      </Statistic.Group>

      <Grid columns={2} style={{ marginTop: 24 }}>
        <Grid.Column>
          <Segment>
            <h3><Icon name="chart bar" /> 收入趋势</h3>
            <p style={{ color: '#999', padding: 40, textAlign: 'center' }}>
              （图表区域 — 接入 ECharts 后可渲染折线图）
            </p>
          </Segment>
        </Grid.Column>
        <Grid.Column>
          <Segment>
            <h3><Icon name="pie chart" /> 科室工作量</h3>
            <p style={{ color: '#999', padding: 40, textAlign: 'center' }}>
              （图表区域 — 接入 ECharts 后可渲染饼图）
            </p>
          </Segment>
        </Grid.Column>
      </Grid>
    </div>
  )
}
