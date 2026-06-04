import { Routes, Route, Navigate } from 'react-router-dom'
import PrivateRoute from './router/PrivateRoute'
import AdminLayout from './layouts/AdminLayout'
import LoginPage from './pages/login/LoginPage'
import DashboardPage from './pages/dashboard/DashboardPage'
import { Suspense, lazy } from 'react'
import { Dimmer, Loader } from 'semantic-ui-react'

const Lazy = ({ children }: { children: React.ReactNode }) => (
  <Suspense fallback={<Dimmer active><Loader>加载中...</Loader></Dimmer>}>{children}</Suspense>
)

const PatientList = lazy(() => import('./pages/user/PatientList'))
const EmployeeList = lazy(() => import('./pages/user/EmployeeList'))
const DepartmentList = lazy(() => import('./pages/user/DepartmentList'))
const RegistrationPage = lazy(() => import('./pages/registration/RegistrationPage'))
const SchedulePage = lazy(() => import('./pages/schedule/SchedulePage'))
const ClinicPage = lazy(() => import('./pages/clinic/ClinicPage'))
const PrescriptionPage = lazy(() => import('./pages/prescription/PrescriptionPage'))
const BillingPage = lazy(() => import('./pages/billing/BillingPage'))
const PharmacyPage = lazy(() => import('./pages/pharmacy/PharmacyPage'))
const ExaminationPage = lazy(() => import('./pages/examination/ExaminationPage'))
const InpatientPage = lazy(() => import('./pages/inpatient/InpatientPage'))
const StatisticsPage = lazy(() => import('./pages/statistics/StatisticsPage'))
const SystemPage = lazy(() => import('./pages/system/SystemPage'))
const EMRPage = lazy(() => import('./pages/emr/EMRPage'))
const OutpatientPage = lazy(() => import('./pages/outpatient/OutpatientPage'))
const FollowupPage = lazy(() => import('./pages/followup/FollowupPage'))
const HealthRecordPage = lazy(() => import('./pages/health_record/HealthRecordPage'))
const NotificationPage = lazy(() => import('./pages/notification/NotificationPage'))

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/" element={<PrivateRoute><AdminLayout /></PrivateRoute>}>
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<Lazy><DashboardPage /></Lazy>} />
        <Route path="user/patients" element={<Lazy><PatientList /></Lazy>} />
        <Route path="user/employees" element={<Lazy><EmployeeList /></Lazy>} />
        <Route path="user/departments" element={<Lazy><DepartmentList /></Lazy>} />
        <Route path="registration" element={<Lazy><RegistrationPage /></Lazy>} />
        <Route path="schedule" element={<Lazy><SchedulePage /></Lazy>} />
        <Route path="clinic" element={<Lazy><ClinicPage /></Lazy>} />
        <Route path="prescription" element={<Lazy><PrescriptionPage /></Lazy>} />
        <Route path="billing" element={<Lazy><BillingPage /></Lazy>} />
        <Route path="pharmacy" element={<Lazy><PharmacyPage /></Lazy>} />
        <Route path="examination" element={<Lazy><ExaminationPage /></Lazy>} />
        <Route path="inpatient" element={<Lazy><InpatientPage /></Lazy>} />
        <Route path="statistics" element={<Lazy><StatisticsPage /></Lazy>} />
        <Route path="system" element={<Lazy><SystemPage /></Lazy>} />
        <Route path="emr" element={<Lazy><EMRPage /></Lazy>} />
        <Route path="outpatient" element={<Lazy><OutpatientPage /></Lazy>} />
        <Route path="followup" element={<Lazy><FollowupPage /></Lazy>} />
        <Route path="health-record" element={<Lazy><HealthRecordPage /></Lazy>} />
        <Route path="notification" element={<Lazy><NotificationPage /></Lazy>} />
      </Route>
    </Routes>
  )
}
