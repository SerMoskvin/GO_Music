import { defineStore } from 'pinia'
import { assessmentApi } from '@/services/assessmentApi'
import type { StudentAssessment, StudentAssessmentCreateDTO, StudentAssessmentUpdateDTO } from '@/types/assessment'

export const useAssessmentStore = defineStore('assessment', {
  state: () => ({
    assessments: [] as StudentAssessment[],
    currentAssessment: null as StudentAssessment | null,
    loading: false,
    error: null as string | null
  }),

  getters: {
    getAssessmentById: (state) => (id: number): StudentAssessment | undefined => {
      return state.assessments.find((a: StudentAssessment) => a.id === id)
    }
  },

  actions: {
    async fetchAll(): Promise<void> {
      this.loading = true
      this.error = null
      try {
        this.assessments = await assessmentApi.getAll()
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error'
      } finally {
        this.loading = false
      }
    },

    async fetchById(id: number): Promise<void> {
      this.loading = true
      this.error = null
      try {
        this.currentAssessment = await assessmentApi.getById(id)
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error'
      } finally {
        this.loading = false
      }
    },

    async createAssessment(data: StudentAssessmentCreateDTO): Promise<StudentAssessment> {
      this.loading = true
      this.error = null
      try {
        const newAssessment = await assessmentApi.create(data)
        this.assessments.push(newAssessment)
        return newAssessment
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error'
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateAssessment(id: number, data: StudentAssessmentUpdateDTO): Promise<StudentAssessment> {
      this.loading = true
      this.error = null
      try {
        const updatedAssessment = await assessmentApi.update(id, data)
        const index = this.assessments.findIndex((a: StudentAssessment) => a.id === id)
        if (index !== -1) {
          this.assessments[index] = updatedAssessment
        }
        if (this.currentAssessment?.id === id) {
          this.currentAssessment = updatedAssessment
        }
        return updatedAssessment
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error'
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteAssessment(id: number): Promise<void> {
      this.loading = true
      this.error = null
      try {
        await assessmentApi.delete(id)
        this.assessments = this.assessments.filter((a: StudentAssessment) => a.id !== id)
        if (this.currentAssessment?.id === id) {
          this.currentAssessment = null
        }
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error'
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})