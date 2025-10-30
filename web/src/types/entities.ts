export interface BaseEntity {
  id: number
  createdAt: string
  updatedAt: string
}

export interface Student extends BaseEntity {
  surname: string
  name: string
  fatherName?: string
  birthday: string
  phoneNumber: string
  groupId: number
  userId?: number
}

export interface Teacher extends BaseEntity {
  surname: string
  name: string
  fatherName?: string
  birthday: string
  phoneNumber: string
  specialization: string
  experience: number
  userId?: number
}

export interface Employee extends BaseEntity {
  surname: string
  name: string
  fatherName?: string
  birthday: string
  phoneNumber: string
  job: string
  workExperience: number
  userId?: number
}

export interface Group extends BaseEntity {
  name: string
  teacherId: number
  instrument: string
  level: string
  schedule: string
}

export interface Lesson extends BaseEntity {
  groupId: number
  teacherId: number
  date: string
  startTime: string
  endTime: string
  topic?: string
  homework?: string
}

export interface Assessment extends BaseEntity {
  lessonId: number
  studentId: number
  taskType: string
  grade: number
  assessmentDate: string
  comments?: string
}