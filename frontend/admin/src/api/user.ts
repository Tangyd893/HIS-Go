import http from './client'
import type { Patient, Employee, Department, PageData, PageQuery } from './types'

export const userApi = {
  /** 患者列表 */
  getPatients(params: PageQuery & { name?: string; phone?: string }): Promise<PageData<Patient>> {
    return http.get('/user/patients', { params })
  },

  /** 创建患者 */
  createPatient(data: Partial<Patient>): Promise<Patient> {
    return http.post('/user/patients', data)
  },

  /** 更新患者 */
  updatePatient(id: string, data: Partial<Patient>): Promise<Patient> {
    return http.put(`/user/patients/${id}`, data)
  },

  /** 删除患者 */
  deletePatient(id: string): Promise<void> {
    return http.delete(`/user/patients/${id}`)
  },

  /** 员工列表 */
  getEmployees(params: PageQuery & { deptId?: string; name?: string }): Promise<PageData<Employee>> {
    return http.get('/user/employees', { params })
  },

  /** 科室列表 */
  getDepartments(): Promise<Department[]> {
    return http.get('/user/departments')
  },
}
