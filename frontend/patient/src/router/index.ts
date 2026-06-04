import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginView.vue'),
    meta: { title: '登录', noAuth: true },
  },
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/DashboardView.vue'),
        meta: { title: '首页', icon: 'HomeOutlined' },
      },
      {
        path: 'appointment',
        name: 'Appointment',
        component: () => import('@/views/appointment/AppointmentView.vue'),
        meta: { title: '预约挂号', icon: 'ScheduleOutlined' },
      },
      {
        path: 'consultation',
        name: 'Consultation',
        component: () => import('@/views/consultation/ConsultationView.vue'),
        meta: { title: '在线问诊', icon: 'MessageOutlined' },
      },
      {
        path: 'prescription',
        name: 'Prescription',
        component: () => import('@/views/prescription/PrescriptionView.vue'),
        meta: { title: '我的处方', icon: 'FileTextOutlined' },
      },
      {
        path: 'report',
        name: 'Report',
        component: () => import('@/views/report/ReportView.vue'),
        meta: { title: '检查报告', icon: 'FileSearchOutlined' },
      },
      {
        path: 'health-record',
        name: 'HealthRecord',
        component: () => import('@/views/health_record/HealthRecordView.vue'),
        meta: { title: '健康档案', icon: 'HeartOutlined' },
      },
      {
        path: 'chronic',
        name: 'Chronic',
        component: () => import('@/views/chronic/ChronicView.vue'),
        meta: { title: '慢病管理', icon: 'SafetyCertificateOutlined' },
      },
      {
        path: 'followup',
        name: 'Followup',
        component: () => import('@/views/followup/FollowupView.vue'),
        meta: { title: '我的随访', icon: 'PhoneOutlined' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory('/patient/'),
  routes,
})

router.beforeEach((to, _from, next) => {
  document.title = (to.meta.title as string) || 'HIS-Go 患者服务中心'
  if (to.meta.noAuth) { next(); return }
  const authStore = useAuthStore()
  authStore.restoreUserInfo()
  if (!authStore.isLoggedIn) { next('/login'); return }
  next()
})

export default router
