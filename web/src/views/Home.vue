<template>
  <div class="home">
    <header class="hero">
      <h1>🎵 GO Music School</h1>
      <p>Управление музыкальной школой</p>
    </header>

    <main class="main-content">
      <section class="features">
        <div class="feature-card">
          <h3>🎼 Расписание</h3>
          <p>Управление занятиями и расписанием</p>
        </div>
        <div class="feature-card">
          <h3>👥 Ученики</h3>
          <p>Ведение базы учеников и групп</p>
        </div>
        <div class="feature-card">
          <h3>📊 Оценки</h3>
          <p>Система оценивания и посещаемости</p>
        </div>
      </section>

      <section class="auth-section">
        <h2>Начать работу</h2>
        <div class="auth-buttons">
          <button @click="goToLogin" class="btn btn-primary">Войти</button>
          <button @click="goToRegister" class="btn btn-secondary">Зарегистрироваться</button>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goToLogin = (): void => {
  router.push('/login')
}

const goToRegister = (): void => {
  router.push('/register')
}
</script>

<style scoped>
.home {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.hero {
  text-align: center;
  padding: 4rem 2rem;
}

.hero h1 {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.hero p {
  font-size: 1.5rem;
  opacity: 0.9;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem 4rem;
}

.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  margin-bottom: 4rem;
}

.feature-card {
  background: rgba(255, 255, 255, 0.1);
  padding: 2rem;
  border-radius: 10px;
  text-align: center;
  backdrop-filter: blur(10px);
}

.feature-card h3 {
  margin-bottom: 1rem;
  font-size: 1.5rem;
}

.auth-section {
  text-align: center;
  background: rgba(255, 255, 255, 0.1);
  padding: 3rem;
  border-radius: 15px;
  backdrop-filter: blur(10px);
}

.auth-section h2 {
  margin-bottom: 2rem;
  font-size: 2rem;
}

.auth-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}

.btn {
  padding: 1rem 2rem;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-primary {
  background: #4CAF50;
  color: white;
}

.btn-primary:hover {
  background: #45a049;
  transform: translateY(-2px);
}

.btn-secondary {
  background: transparent;
  color: white;
  border: 2px solid white;
}

.btn-secondary:hover {
  background: white;
  color: #667eea;
  transform: translateY(-2px);
}

@media (max-width: 768px) {
  .hero h1 {
    font-size: 2rem;
  }
  
  .auth-buttons {
    flex-direction: column;
  }
  
  .btn {
    width: 100%;
  }
}
</style>