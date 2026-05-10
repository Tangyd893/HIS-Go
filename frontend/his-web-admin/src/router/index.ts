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
        meta: { title: '工作台', icon: 'DashboardOutlined' },
      },
      {
        path: 'user/patients',
        name: 'PatientList',
        component: () => import('@/views/user/PatientList.vue'),
        meta: { title: '患者管理', icon: 'TeamOutlined' },
      },
      {
        path: 'user/employees',
        name: 'EmployeeList',
        component: () => import('@/views/user/EmployeeList.vue'),
        meta: { title: '员工管理', icon: 'IdcardOutlined' },
      },
      {
        path: 'user/departments',
        name: 'DepartmentList',
        component: () => import('@/views/user/DepartmentList.vue'),
        meta: { title: '科室管理', icon: 'ClusterOutlined' },
      },
      {
        path: 'registration',
        name: 'Registration',
        component: () => import('@/views/registration/RegistrationView.vue'),
        meta: { title: '挂号管理', icon: 'FormOutlined' },
      },
      {
        path: 'schedule',
        name: 'Schedule',
        component: () => import('@/views/schedule/ScheduleView.vue'),
        meta: { title: '排班管理', icon: 'ScheduleOutlined' },
      },
      {
        path: 'clinic',
        name: 'Clinic',
        component: () => import('@/views/clinic/ClinicView.vue'),
        meta: { title: '门诊诊疗', icon: 'MedicineBoxOutlined' },
      },
      {
        path: 'prescription',
        name: 'Prescription',
        component: () => import('@/views/prescription/PrescriptionView.vue'),
        meta: { title: '处方管理', icon: 'FileTextOutlined' },
      },
      {
        path: 'billing',
        name: 'Billing',
        component: () => import('@/views/billing/BillingView.vue'),
        meta: { title: '收费结算', icon: 'DollarOutlined' },
      },
      {
        path: 'pharmacy',
        name: 'Pharmacy',
        component: () => import('@/views/pharmacy/PharmacyView.vue'),
        meta: { title: '药房管理', icon: 'ExperimentOutlined' },
      },
      {
        path: 'examination',
        name: 'Examination',
        component: () => import('@/views/examination/ExaminationView.vue'),
        meta: { title: '检查检验', icon: 'FundProjectionScreenOutlined' },
      },
      {
        path: 'inpatient',
        name: 'Inpatient',
        component: () => import('@/views/inpatient/InpatientView.vue'),
        meta: { title: '住院管理', icon: 'HomeOutlined' },
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: () => import('@/views/statistics/StatisticsView.vue'),
        meta: { title: '数据统计', icon: 'BarChartOutlined' },
      },
      {
        path: 'system',
        name: 'System',
        component: () => import('@/views/system/SystemView.vue'),
        meta: { title: '系统设置', icon: 'SettingOutlined' },
      },
      {
        path: 'emr',
        name: 'EMR',
        component: () => import('@/views/emr/EMRView.vue'),
        meta: { title: '电子病历', icon: 'FileProtectOutlined' },
      },
      {
        path: 'outpatient',
        name: 'Outpatient',
        component: () => import('@/views/outpatient/OutpatientView.vue'),
        meta: { title: '院外服务', icon: 'GlobalOutlined' },
      },
      {
        path: 'followup',
        name: 'Followup',
        component: () => import('@/views/followup/FollowupView.vue'),
        meta: { title: '随访管理', icon: 'PhoneOutlined' },
      },
      {
        path: 'health-record',
        name: 'HealthRecord',
        component: () => import('@/views/health_record/HealthRecordView.vue'),
        meta: { title: '健康档案', icon: 'HeartOutlined' },
      },
      {
        path: 'notification',
        name: 'Notification',
        component: () => import('@/views/notification/NotificationView.vue'),
        meta: { title: '消息通知', icon: 'BellOutlined' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  document.title = (to.meta.title as string) || 'HIS-Go 医院管理系统'

  if (to.meta.noAuth) {
    next()
    return
  }

  const authStore = useAuthStore()
  authStore.restoreUserInfo()

  if (!authStore.isLoggedIn) {
    next('/login')
    return
  }

  next()
})

export default router
