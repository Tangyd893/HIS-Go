// 本地演示：先启动 frontend/patient (Vue 3)，默认 5174
//   cd frontend/patient && npm run dev -- --host
// 微信开发者工具 → 详情 → 本地设置 → 勾选「不校验合法域名、web-view、TLS…」

// ========== H5 地址配置 ==========
// 开发环境：模拟器用 127.0.0.1，真机需改为 PC 局域网 IP
const DEV_PATIENT_H5  = 'http://127.0.0.1:5174'
// 云演示环境：替换为实际域名
const CLOUD_PATIENT_H5 = 'https://your-domain.com/patient/'

// 切换方式：修改下面这行
const USE_CLOUD = false
const DEFAULT_PATIENT_H5 = USE_CLOUD ? CLOUD_PATIENT_H5 : DEV_PATIENT_H5

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