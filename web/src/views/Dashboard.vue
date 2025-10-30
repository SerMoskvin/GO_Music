<template>
  <div class="dashboard">
    <div class="dashboard-container">
      <h1>📊 Панель управления</h1>
      <p>Добро пожаловать, {{ userRole }}!</p>
      
      <div v-if="permissionsLoading" class="loading">🔄 Загрузка прав доступа...</div>
      <div v-else-if="permissionsError" class="error">❌ {{ permissionsError }}</div>
      
      <div v-else class="stats-grid">
        <div 
          v-for="section in availableSections" 
          :key="section.url"
          class="stat-card"
          :class="{ 'write-access': section.can_write }"
        >
          <h3>{{ section.name }}</h3>
          <p>{{ section.can_write ? 'Полный доступ' : 'Только просмотр' }}</p>
          <router-link :to="section.url" class="btn btn-primary">
            {{ section.can_write ? 'Управлять' : 'Просмотреть' }}
          </router-link>
        </div>
      </div>

      <div class="dashboard-info">
        <p v-if="isOwnRecordsOnly">🔐 У вас доступ только к вашим собственным записям</p>
        <p v-else>🔓 У вас полный доступ ко всем записям системы</p>
      </div>

      <button @click="logout" class="btn btn-secondary">🚪 Выйти</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { usePermissionsStore } from '@/stores/permissionsStore'

const router = useRouter()
const permissionsStore = usePermissionsStore()

const { availableSections, isOwnRecordsOnly, loading: permissionsLoading, error: permissionsError } = storeToRefs(permissionsStore)

const userRole = computed(() => {
  const token = localStorage.getItem('token')
  if (!token) return 'guest'
  
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.role || 'student'
  } catch {
    return 'guest'
  }
})

const logout = (): void => {
  localStorage.removeItem('token')
  permissionsStore.clearPermissions()
  router.push('/')
}
</script>

<style scoped>
/* Стили остаются такими же, добавь только: */
.stat-card.write-access {
  border-left: 4px solid #42b983;
}

.dashboard-info {
  margin: 2rem 0;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  text-align: center;
}
</style>