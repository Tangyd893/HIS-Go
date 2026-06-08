// 本地演示（Docker Compose + Nginx）：make demo-patient 或 compose up 后使用
//   http://127.0.0.1/patient/
// 可选 Vite 热更新：cd frontend/patient && npm run dev -- --host → :5174/patient/
// 微信开发者工具 → 详情 → 本地设置 → 勾选「不校验合法域名、web-view、TLS…」

// ========== H5 地址配置 ==========
// Docker 演示（推荐验收）：Nginx 托管患者 H5 + API
const DEV_PATIENT_H5  = 'http://127.0.0.1/patient/'
// Vite 开发服务器（热更新）
const VITE_PATIENT_H5 = 'http://127.0.0.1:5174/patient/'
// 云演示环境（已搁置）
const CLOUD_PATIENT_H5 = 'https://your-domain.com/patient/'

// 切换：docker | vite | cloud
const H5_MODE = 'docker'
const DEFAULT_PATIENT_H5 =
  H5_MODE === 'cloud' ? CLOUD_PATIENT_H5 :
  H5_MODE === 'vite'  ? VITE_PATIENT_H5  : DEV_PATIENT_H5

Page({
  data: {
    patientUrl: DEFAULT_PATIENT_H5,
    loading: true,
    loadError: false,
  },

  onLoad() {
    // 模拟加载延迟，避免白屏闪烁
    setTimeout(() => {
      this.setData({ loading: false })
    }, 500)
  },

  onWebviewError(e) {
    console.error('web-view 加载失败:', e)
    this.setData({
      loading: false,
      loadError: true,
    })
  },

  retryLoad() {
    this.setData({
      loading: true,
      loadError: false,
      patientUrl: '',
    })
    // 重新设置 URL 触发重新加载
    setTimeout(() => {
      this.setData({ patientUrl: DEFAULT_PATIENT_H5 })
    }, 100)
  },
})