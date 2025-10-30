import { api } from './api'
import type { StudentAssessment, StudentAssessmentCreateDTO, StudentAssessmentUpdateDTO } from '@/types/assessment'

export const assessmentApi = {
  // GET /assessments
  async getAll(): Promise<StudentAssessment[]> {
    return api.get('/assessments')
  },

  // GET /assessments/{id}
  async getById(id: number): Promise<StudentAssessment> {
    return api.get(`/assessments/${id}`)
  },

  // POST /assessments
  async create(data: StudentAssessmentCreateDTO): Promise<StudentAssessment> {
    return api.post('/assessments', data)
  },

  // PUT /assessments/{id}
  async update(id: number, data: StudentAssessmentUpdateDTO): Promise<StudentAssessment> {
    return api.put(`/assessments/${id}`, data)
  },

  // PATCH /assessments/{id}
  async partialUpdate(id: number, data: StudentAssessmentUpdateDTO): Promise<StudentAssessment> {
    return api.patch(`/assessments/${id}`, data)
  },

  // DELETE /assessments/{id}
  async delete(id: number): Promise<void> {
    return api.delete(`/assessments/${id}`)
  }
}