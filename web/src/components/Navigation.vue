<template>
  <nav class="navigation">
    <div class="nav-container">
      <router-link to="/" class="logo">
        🎵 GO Music
      </router-link>
      
      <div class="nav-links">
        <template v-if="isAuthenticated">
          <!-- Для всех авторизованных -->
          <router-link to="/dashboard">📊 Панель</router-link>
          
          <!-- Динамическое меню из конфига -->
          <template v-if="availableSections.length > 0">
            <router-link 
              v-for="section in availableSections" 
              :key="section.url"
              :to="section.url"
              class="nav-section"
              :class="{ 'write-access': section.can_write }"
              :title="section.can_write ? 'Полный доступ' : 'Только чтение'"
            >
              {{ section.name }}
              <span v-if="section.can_write" class="write-indicator">✏️</span>
            </router-link>
          </template>

          <div class="user-info">
            <span class="user-role">({{ userRole }})</span>
            <button @click="logout" class="logout-btn">🚪 Выйти</button>
          </div>
        </template>

        <template v-else>
          <router-link to="/">🏠 Главная</router-link>
          <router-link to="/info">ℹ️ Информация</router-link>
          <router-link to="/login">🔑 Войти</router-link>
        </template>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { usePermissionsStore } from '@/stores/permissionsStore'
import { storeToRefs } from 'pinia'

const router = useRouter()
const permissionsStore = usePermissionsStore()

const { availableSections } = storeToRefs(permissionsStore)

const isAuthenticated = computed(() => {
  return !!localStorage.getItem('token')
})

const userRole = computed(() => {
  const token = localStorage.getItem('token')
  if (!token) return null
  
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.role || 'student'
  } catch {
    return null
  }
})

// Загружаем пермишены когда меняется роль
watch(userRole, (newRole) => {
  if (newRole) {
    permissionsStore.loadUserPermissions(newRole)
  }
}, { immediate: true })

// Также загружаем при монтировании
onMounted(() => {
  if (userRole.value) {
    permissionsStore.loadUserPermissions(userRole.value)
  }
})

const logout = (): void => {
  localStorage.removeItem('token')
  permissionsStore.clearPermissions()
  router.push('/')
}
</script>

<style scoped>
.navigation {
  background: #2c3e50;
  padding: 1rem 0;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.nav-container {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 2rem;
}

.logo {
  font-size: 1.5rem;
  font-weight: bold;
  text-decoration: none;
  color: white;
}

.nav-links {
  display: flex;
  gap: 1rem;
  align-items: center;
  flex-wrap: wrap;
}

.nav-links a {
  color: white;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  transition: all 0.3s ease;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.nav-links a:hover,
.nav-links a.router-link-active {
  background: rgba(255, 255, 255, 0.15);
  color: #42b983;
}

.nav-section.write-access {
  border-left: 3px solid #42b983;
}

.write-indicator {
  font-size: 0.75rem;
  opacity: 0.8;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-left: 1rem;
  padding-left: 1rem;
  border-left: 1px solid rgba(255, 255, 255, 0.3);
}

.user-role {
  color: #ccc;
  font-size: 0.875rem;
  font-style: italic;
}

.logout-btn {
  background: transparent;
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.5);
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.875rem;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: white;
}

@media (max-width: 768px) {
  .nav-container {
    flex-direction: column;
    gap: 1rem;
  }
  
  .nav-links {
    justify-content: center;
    text-align: center;
  }
  
  .user-info {
    margin-left: 0;
    padding-left: 0;
    border-left: none;
    border-top: 1px solid rgba(255, 255, 255, 0.3);
    padding-top: 0.5rem;
    width: 100%;
    justify-content: center;
  }
}
</style>