<template>
  <div class="info">
    <div class="container">
      <h1>ℹ️ О системе</h1>
      
      <section class="info-section">
        <h2>GO Music School</h2>
        <p>Система управления музыкальной школой, разработанная на Go и Vue.js</p>
      </section>

      <section class="info-section">
        <h3>🎯 Функциональность</h3>
        <ul>
          <li>Управление расписанием занятий</li>
          <li>Ведение базы учеников и преподавателей</li>
          <li>Система оценивания и посещаемости</li>
          <li>Управление аудиториями и инструментами</li>
          <li>Разграничение прав доступа по ролям</li>
        </ul>
      </section>

      <section class="info-section">
        <h3>👥 Роли пользователей</h3>
        <div class="roles-grid">
          <div class="role-card">
            <h4>Администратор</h4>
            <p>Полный доступ ко всем функциям системы</p>
          </div>
          <div class="role-card">
            <h4>Преподаватель</h4>
            <p>Управление занятиями, оценками, посещаемостью</p>
          </div>
          <div class="role-card">
            <h4>Ученик</h4>
            <p>Просмотр расписания, оценок, посещаемости</p>
          </div>
        </div>
      </section>

      <div class="back-home">
        <button @click="goHome" class="btn btn-primary">← На главную</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Info',
  methods: {
    goHome() {
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.info {
  min-height: 100vh;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
  padding: 2rem 0;
}

.container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 2rem;
}

.info h1 {
  text-align: center;
  font-size: 2.5rem;
  margin-bottom: 3rem;
}

.info-section {
  background: rgba(255, 255, 255, 0.1);
  padding: 2rem;
  border-radius: 10px;
  margin-bottom: 2rem;
  backdrop-filter: blur(10px);
}

.info-section h2, .info-section h3 {
  margin-bottom: 1rem;
}

.info-section ul {
  list-style: none;
  padding: 0;
}

.info-section li {
  padding: 0.5rem 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
}

.info-section li:last-child {
  border-bottom: none;
}

.roles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.role-card {
  background: rgba(255, 255, 255, 0.15);
  padding: 1.5rem;
  border-radius: 8px;
  text-align: center;
}

.role-card h4 {
  margin-bottom: 0.5rem;
  color: #ffeb3b;
}

.back-home {
  text-align: center;
  margin-top: 3rem;
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
</style>