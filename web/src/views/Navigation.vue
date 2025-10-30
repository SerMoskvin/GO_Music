<template>
  <nav class="navigation">
    <div class="nav-container">
      <router-link to="/" class="logo">
        🎵 GO Music
      </router-link>
      
      <div class="nav-links">
        <router-link to="/">Главная</router-link>
        <router-link to="/info">Информация</router-link>
        <router-link v-if="!isAuthenticated" to="/login">Войти</router-link>
        <router-link v-if="!isAuthenticated" to="/register">Регистрация</router-link>
        <button v-if="isAuthenticated" @click="logout" class="logout-btn">Выйти</button>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const isAuthenticated = computed(() => {
  return !!localStorage.getItem('token')
})

const logout = (): void => {
  localStorage.removeItem('token')
  router.push('/')
}
</script>

<style scoped>
.navigation {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  padding: 1rem 0;
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-container {
  max-width: 1200px;
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
  gap: 2rem;
  align-items: center;
}

.nav-links a {
  color: white;
  text-decoration: none;
  transition: opacity 0.3s ease;
}

.nav-links a:hover,
.nav-links a.router-link-active {
  opacity: 0.8;
}

.logout-btn {
  background: transparent;
  color: white;
  border: 1px solid white;
  padding: 0.5rem 1rem;
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.logout-btn:hover {
  background: white;
  color: #667eea;
}
</style>