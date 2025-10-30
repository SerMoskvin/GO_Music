<template>
  <div class="register">
    <div class="register-container">
      <h2>Регистрация</h2>
      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-group">
          <label for="login">Логин</label>
          <input
            id="login"
            v-model="form.login"
            type="text"
            required
            placeholder="Придумайте логин"
          />
        </div>
        
        <div class="form-group">
          <label for="password">Пароль</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            required
            placeholder="Придумайте пароль"
          />
        </div>

        <div class="form-group">
          <label for="surname">Фамилия</label>
          <input
            id="surname"
            v-model="form.surname"
            type="text"
            required
            placeholder="Введите фамилию"
          />
        </div>

        <div class="form-group">
          <label for="name">Имя</label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            required
            placeholder="Введите имя"
          />
        </div>

        <div class="form-group">
          <label for="email">Email</label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            required
            placeholder="Введите email"
          />
        </div>

        <div class="form-group">
          <label for="role">Роль</label>
          <select id="role" v-model="form.role" required>
            <option value="">Выберите роль</option>
            <option value="student">Ученик</option>
            <option value="teacher">Преподаватель</option>
            <option value="admin">Администратор</option>
          </select>
        </div>

        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
        </button>

        <div class="login-link">
          Уже есть аккаунт? <router-link to="/login">Войти</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

interface RegisterForm {
  login: string
  password: string
  surname: string
  name: string
  email: string
  role: string
}

const router = useRouter()
const loading = ref(false)
const form = ref<RegisterForm>({
  login: '',
  password: '',
  surname: '',
  name: '',
  email: '',
  role: ''
})

const handleRegister = async (): Promise<void> => {
  loading.value = true
  
  try {
    const response = await fetch('http://localhost:8080/authentication/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(form.value)
    })

    const result = await response.json()
    
    if (response.ok) {
      alert('Регистрация успешна! Теперь вы можете войти.')
      router.push('/login')
    } else {
      alert('Ошибка регистрации: ' + (result.error || 'Неизвестная ошибка'))
    }
  } catch (error) {
    alert('Ошибка сети: ' + error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.register-container {
  background: rgba(255, 255, 255, 0.95);
  padding: 3rem;
  border-radius: 15px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  width: 100%;
  max-width: 400px;
  backdrop-filter: blur(10px);
}

.register-container h2 {
  text-align: center;
  margin-bottom: 2rem;
  color: #333;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #333;
  font-weight: 500;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
}

.btn {
  width: 100%;
  padding: 1rem;
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

.btn-primary:hover:not(:disabled) {
  background: #45a049;
  transform: translateY(-2px);
}

.btn-primary:disabled {
  background: #cccccc;
  cursor: not-allowed;
  transform: none;
}

.login-link {
  text-align: center;
  margin-top: 1.5rem;
  color: #666;
}

.login-link a {
  color: #667eea;
  text-decoration: none;
}

.login-link a:hover {
  text-decoration: underline;
}
</style>