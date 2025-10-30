export interface StudentAssessment {
  id: number
  lesson_id: number
  student_id: number
  task_type: string
  grade: number
  assessment_date: string // DD.MM.YYYY
}

export interface StudentAssessmentCreateDTO {
  lesson_id: number
  student_id: number
  task_type: string
  grade: number
  assessment_date: string // DD.MM.YYYY
}

export interface StudentAssessmentUpdateDTO {
  lesson_id?: number
  student_id?: number
  task_type?: string
  grade?: number
  assessment_date?: string // DD.MM.YYYY
}