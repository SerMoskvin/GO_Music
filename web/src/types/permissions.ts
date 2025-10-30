export interface Section {
  name: string
  url: string
  can_read: boolean
  can_write: boolean
}

export interface RolePermissions {
  own_records_only: boolean
  sections: Section[]
}

export interface PermissionsConfig {
  roles: {
    [role: string]: RolePermissions
  }
}