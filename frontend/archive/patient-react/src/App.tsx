import { Routes, Route, Navigate } from 'react-router-dom'
import { Suspense, lazy } from 'react'
import { Dimmer, Loader } from 'semantic-ui-react'
import PrivateRoute from './router/PrivateRoute'
import PatientLayout from './layouts/PatientLayout'
import LoginPage from './pages/login/LoginPage'

const Lazy = ({ children }: { children: React.ReactNode }) => (
  <Suspense fallback={<Dimmer active><Loader>加载中...</Loader></Dimmer>}>{children}</Suspense>
)

const DashboardPage = lazy(() => import('./pages/dashboard/DashboardPage'))
const AppointmentPage = lazy(() => import('./pages/appointment/AppointmentPage'))
const ConsultationPage = lazy(() => import('./pages/consultation/ConsultationPage'))
const PrescriptionPage = lazy(() => import('./pages/prescription/PrescriptionPage'))
const ReportPage = lazy(() => import('./pages/report/ReportPage'))
const HealthRecordPage = lazy(() => import('./pages/health_record/HealthRecordPage'))
const ChronicPage = lazy(() => import('./pages/chronic/ChronicPage'))
const FollowupPage = lazy(() => import('./pages/followup/FollowupPage'))

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/" element={<PrivateRoute><PatientLayout /></PrivateRoute>}>
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<Lazy><DashboardPage /></Lazy>} />
        <Route path="appointment" element={<Lazy><AppointmentPage /></Lazy>} />
        <Route path="consultation" element={<Lazy><ConsultationPage /></Lazy>} />
        <Route path="prescription" element={<Lazy><PrescriptionPage /></Lazy>} />
        <Route path="report" element={<Lazy><ReportPage /></Lazy>} />
        <Route path="health-record" element={<Lazy><HealthRecordPage /></Lazy>} />
        <Route path="chronic" element={<Lazy><ChronicPage /></Lazy>} />
        <Route path="followup" element={<Lazy><FollowupPage /></Lazy>} />
      </Route>
    </Routes>
  )
}
