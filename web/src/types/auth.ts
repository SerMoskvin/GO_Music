export interface User {
  id: number
  username: string
  email: string
  role: UserRole
  firstName: string
  lastName: string
}

export type UserRole = 'student' | 'employee' | 'admin' | 'teacher'

export interface PermissionSection {
  name: string
  url: string
  canRead: boolean
  canWrite: boolean
}

export interface RolePermissions {
  role: UserRole
  sections: PermissionSection[]
  ownRecordsOnly: boolean
}

export interface PermissionConfig {
  roles: Record<UserRole, RolePermissions>
}

export interface LoginResponse {
  token: string
  user: User
}

export interface LoginCredentials {
  username: string
  password: string
}