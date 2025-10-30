import type { PermissionsConfig, RolePermissions } from '@/types/permissions'

// Загружаем конфиг с бэкенда
export const permissionsService = {
  async getPermissionsConfig(): Promise<PermissionsConfig> {
    try {
      const response = await fetch('/api/permissions/config') // или прямой путь к файлу
      if (!response.ok) {
        throw new Error('Failed to load permissions config')
      }
      return await response.json()
    } catch (error) {
      console.error('Error loading permissions config:', error)
      // Fallback конфиг на случай ошибки
      return getFallbackConfig()
    }
  },

  async getUserPermissions(role: string): Promise<RolePermissions | null> {
    const config = await this.getPermissionsConfig()
    return config.roles[role] || null
  }
}

// Fallback конфиг если не удалось загрузить
function getFallbackConfig(): PermissionsConfig {
  return {
    roles: {
      admin: {
        own_records_only: false,
        sections: [
          { name: "Расписание", url: "/schedules", can_read: true, can_write: true },
          { name: "Занятия", url: "/lessons", can_read: true, can_write: true },
          { name: "Сотрудники", url: "/employees", can_read: true, can_write: true },
          { name: "Аудитория", url: "/audiences", can_read: true, can_write: true },
          { name: "Инструмент", url: "/instruments", can_read: true, can_write: true },
          { name: "Пользователь", url: "/users", can_read: true, can_write: true },
          { name: "Ученики", url: "/students", can_read: true, can_write: true },
          { name: "Группы", url: "/study-groups", can_read: true, can_write: true },
          { name: "Оценки", url: "/assessments", can_read: true, can_write: true },
          { name: "Посещение", url: "/attendances", can_read: true, can_write: true },
          { name: "Программа", url: "/programms", can_read: true, can_write: true }
        ]
      },
      teacher: {
        own_records_only: true,
        sections: [
          { name: "Оценки", url: "/assessments", can_read: true, can_write: true },
          { name: "Посещение", url: "/attendances", can_read: true, can_write: true }
        ]
      },
      student: {
        own_records_only: true,
        sections: [
          { name: "Оценки", url: "/assessments", can_read: true, can_write: false },
          { name: "Посещение", url: "/attendances", can_read: true, can_write: false },
          { name: "Инструмент", url: "/instruments", can_read: true, can_write: false }
        ]
      },
      employee: {
        own_records_only: false,
        sections: [
          { name: "Аудитория", url: "/audiences", can_read: true, can_write: true },
          { name: "Инструмент", url: "/instruments", can_read: true, can_write: true }
        ]
      }
    }
  }
}