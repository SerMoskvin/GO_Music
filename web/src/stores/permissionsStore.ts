import { defineStore } from 'pinia'
import { permissionsService } from '@/services/permissionsService'
import type { RolePermissions, Section } from '@/types/permissions'

export const usePermissionsStore = defineStore('permissions', {
  state: () => ({
    userPermissions: null as RolePermissions | null,
    loading: false,
    error: null as string | null
  }),

  getters: {
    // Доступные секции для текущего пользователя
    availableSections: (state): Section[] => {
      return state.userPermissions?.sections.filter(section => section.can_read) || []
    },

    // Может ли пользователь писать в секцию
    canWriteToSection: (state) => (sectionUrl: string): boolean => {
      const section = state.userPermissions?.sections.find(s => s.url === sectionUrl)
      return section?.can_write || false
    },

    // Только свои записи?
    isOwnRecordsOnly: (state): boolean => {
      return state.userPermissions?.own_records_only || false
    }
  },

  actions: {
    async loadUserPermissions(role: string) {
      this.loading = true
      this.error = null
      try {
        this.userPermissions = await permissionsService.getUserPermissions(role)
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error'
        console.error('Failed to load permissions:', error)
      } finally {
        this.loading = false
      }
    },

    clearPermissions() {
      this.userPermissions = null
    }
  }
})