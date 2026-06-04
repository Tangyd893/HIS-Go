import http from './client'
import type { Patient, Employee, Department, PageData } from './types'

export const userApi = {
  // ---- 患者 ----
  getPatients(params?: Record<string, unknown>): Promise<PageData<Patient>> {
    return http.get('/user/patients', { params })
  },
  createPatient(data: Partial<Patient>): Promise<Patient> {
    return http.post('/user/patients', data)
  },
  updatePatient(id: string, data: Partial<Patient>): Promise<Patient> {
    return http.put(`/user/patients/${id}`, data)
  },
  deletePatient(id: string): Promise<void> {
    return http.delete(`/user/patients/${id}`)
  },
  // ---- 员工 ----
  getEmployees(params?: Record<string, unknown>): Promise<PageData<Employee>> {
    return http.get('/user/employees', { params })
  },
  createEmployee(data: Partial<Employee>): Promise<Employee> {
    return http.post('/user/employees', data)
  },
  updateEmployee(id: string, data: Partial<Employee>): Promise<Employee> {
    return http.put(`/user/employees/${id}`, data)
  },
  deleteEmployee(id: string): Promise<void> {
    return http.delete(`/user/employees/${id}`)
  },
  // ---- 科室 ----
  getDepartments(params?: Record<string, unknown>): Promise<PageData<Department>> {
    return http.get('/user/departments', { params })
  },
  createDepartment(data: Partial<Department>): Promise<Department> {
    return http.post('/user/departments', data)
  },
  updateDepartment(id: string, data: Partial<Department>): Promise<Department> {
    return http.put(`/user/departments/${id}`, data)
  },
  deleteDepartment(id: string): Promise<void> {
    return http.delete(`/user/departments/${id}`)
  },
}
