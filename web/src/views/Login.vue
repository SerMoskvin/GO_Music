<template>
  <div class="login">
    <div class="login-container">
      <h2>Вход в систему</h2>
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="login">Логин</label>
          <input
            id="login"
            v-model="form.login"
            type="text"
            required
            placeholder="Введите ваш логин"
          />
        </div>
        
        <div class="form-group">
          <label for="password">Пароль</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            required
            placeholder="Введите ваш пароль"
          />
        </div>

        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? 'Вход...' : 'Войти' }}
        </button>

      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

interface LoginForm {
  login: string
  password: string
}

const router = useRouter()
const loading = ref(false)
const form = ref<LoginForm>({
  login: '',
  password: ''
})

const handleLogin = async (): Promise<void> => {
  loading.value = true
  
  try {
    const response = await fetch('http://localhost:8080/authentication/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(form.value)
    })

    const result = await response.json()
    
    if (result.status === 'success') {
      localStorage.setItem('token', result.data.token)
      router.push('/dashboard')
    } else {
      alert('Ошибка входа: ' + (result.error || 'Неверные учетные данные'))
    }
  } catch (error) {
    alert('Ошибка сети: ' + error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.login-container {
  background: rgba(255, 255, 255, 0.95);
  padding: 3rem;
  border-radius: 15px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  width: 100%;
  max-width: 400px;
  backdrop-filter: blur(10px);
}

.login-container h2 {
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

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.form-group input:focus {
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

.register-link {
  text-align: center;
  margin-top: 1.5rem;
  color: #666;
}

.register-link a {
  color: #667eea;
  text-decoration: none;
}

.register-link a:hover {
  text-decoration: underline;
}
</style>