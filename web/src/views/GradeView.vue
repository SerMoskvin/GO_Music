<template>
  <div class="grades-view">
    <div class="page-header">
      <h1>📊 Мои оценки</h1>
      <p>История ваших оценок по всем предметам</p>
    </div>

    <div class="content-grid">
      <!-- Статистика -->
      <div class="stats-cards">
        <div class="stat-card primary">
          <div class="stat-icon">⭐</div>
          <div class="stat-content">
            <h3>Средний балл</h3>
            <span class="stat-value">4.5</span>
            <span class="stat-change positive">+0.2 с начала месяца</span>
          </div>
        </div>
        
        <div class="stat-card success">
          <div class="stat-icon">📈</div>
          <div class="stat-content">
            <h3>Лучший предмет</h3>
            <span class="stat-value">Сольфеджио</span>
            <span class="stat-change">Средний балл: 4.8</span>
          </div>
        </div>
        
        <div class="stat-card warning">
          <div class="stat-icon">🎯</div>
          <div class="stat-content">
            <h3>Цель на месяц</h3>
            <span class="stat-value">4.7</span>
            <span class="stat-change">Осталось: +0.2</span>
          </div>
        </div>
      </div>

      <!-- Фильтры -->
      <div class="filters-card">
        <h3>Фильтры</h3>
        <div class="filters">
          <select v-model="selectedSubject" class="filter-select">
            <option value="">Все предметы</option>
            <option value="solfeggio">Сольфеджио</option>
            <option value="guitar">Гитара</option>
            <option value="theory">Теория музыки</option>
            <option value="vocal">Вокал</option>
          </select>
          
          <select v-model="selectedMonth" class="filter-select">
            <option value="">Все время</option>
            <option value="december">Декабрь 2024</option>
            <option value="november">Ноябрь 2024</option>
            <option value="october">Октябрь 2024</option>
          </select>
          
          <button class="filter-btn active">Все оценки</button>
          <button class="filter-btn">Только 4 и 5</button>
        </div>
      </div>

      <!-- Таблица оценок -->
      <div class="grades-table-card">
        <div class="table-header">
          <h3>История оценок</h3>
          <div class="table-actions">
            <button class="btn btn-outline">
              <span>📥</span> Экспорт
            </button>
          </div>
        </div>

        <div class="table-container">
          <table class="grades-table">
            <thead>
              <tr>
                <th>Предмет</th>
                <th>Тип задания</th>
                <th>Оценка</th>
                <th>Дата</th>
                <th>Преподаватель</th>
                <th>Комментарий</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="grade in filteredGrades" :key="grade.id" class="grade-row">
                <td>
                  <div class="subject-cell">
                    <span class="subject-icon">{{ grade.subjectIcon }}</span>
                    {{ grade.subject }}
                  </div>
                </td>
                <td>{{ grade.taskType }}</td>
                <td>
                  <span :class="['grade-badge', grade.gradeClass]">
                    {{ grade.grade }}
                  </span>
                </td>
                <td>{{ grade.date }}</td>
                <td>
                  <div class="teacher-cell">
                    <span class="teacher-avatar">ИП</span>
                    {{ grade.teacher }}
                  </div>
                </td>
                <td>
                  <span class="comment" :title="grade.comment">
                    {{ grade.comment }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- График успеваемости -->
      <div class="chart-card">
        <h3>📈 Динамика успеваемости</h3>
        <div class="chart-placeholder">
          <div class="chart-bars">
            <div v-for="week in performanceData" :key="week.week" class="chart-bar-container">
              <div class="chart-bar" :style="{ height: week.height + '%' }"></div>
              <span class="chart-label">{{ week.week }}</span>
            </div>
          </div>
        </div>
        <div class="chart-legend">
          <div class="legend-item">
            <span class="legend-color excellent"></span>
            <span>Отлично (5)</span>
          </div>
          <div class="legend-item">
            <span class="legend-color good"></span>
            <span>Хорошо (4)</span>
          </div>
          <div class="legend-item">
            <span class="legend-color average"></span>
            <span>Удовлетворительно (3)</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const selectedSubject = ref('')
const selectedMonth = ref('')

const grades = ref([
  {
    id: 1,
    subject: 'Сольфеджио',
    subjectIcon: '🎼',
    taskType: 'Интервалы',
    grade: 5,
    gradeClass: 'excellent',
    date: '15.12.2024',
    teacher: 'Иван Петров',
    comment: 'Отличное исполнение!'
  },
  {
    id: 2,
    subject: 'Гитара',
    subjectIcon: '🎸',
    taskType: 'Аккорды',
    grade: 4,
    gradeClass: 'good',
    date: '14.12.2024',
    teacher: 'Мария Сидорова',
    comment: 'Хорошо, но нужно поработать над переходом'
  },
  {
    id: 3,
    subject: 'Теория музыки',
    subjectIcon: '📚',
    taskType: 'Тест',
    grade: 5,
    gradeClass: 'excellent',
    date: '12.12.2024',
    teacher: 'Алексей Козлов',
    comment: 'Превосходные знания'
  },
  {
    id: 4,
    subject: 'Вокал',
    subjectIcon: '🎤',
    taskType: 'Распевка',
    grade: 4,
    gradeClass: 'good',
    date: '10.12.2024',
    teacher: 'Елена Николаева',
    comment: 'Хороший прогресс'
  },
  {
    id: 5,
    subject: 'Сольфеджио',
    subjectIcon: '🎼',
    taskType: 'Ритм',
    grade: 5,
    gradeClass: 'excellent',
    date: '08.12.2024',
    teacher: 'Иван Петров',
    comment: 'Идеально!'
  }
])

const performanceData = ref([
  { week: 'Нед. 1', height: 85 },
  { week: 'Нед. 2', height: 78 },
  { week: 'Нед. 3', height: 92 },
  { week: 'Нед. 4', height: 88 },
  { week: 'Нед. 5', height: 95 }
])

const filteredGrades = computed(() => {
  let filtered = grades.value
  
  if (selectedSubject.value) {
    filtered = filtered.filter(grade => 
      grade.subject.toLowerCase().includes(selectedSubject.value.toLowerCase())
    )
  }
  
  return filtered
})
</script>

<style scoped>
.grades-view {
  padding: 2rem;
  background: #f8f9fa;
  min-height: 100vh;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2.5rem;
  margin: 0 0 0.5rem 0;
  color: #333;
}

.page-header p {
  color: #666;
  margin: 0;
  font-size: 1.1rem;
}

.content-grid {
  display: grid;
  gap: 1.5rem;
  grid-template-columns: 1fr 1fr 1fr;
  grid-template-areas: 
    "stats stats stats"
    "filters table table"
    "chart chart chart";
}

.stats-cards {
  grid-area: stats;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
}

.stat-card {
  background: white;
  padding: 1.5rem;
  border-radius: 15px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  display: flex;
  align-items: center;
  gap: 1rem;
  border-left: 4px solid;
}

.stat-card.primary {
  border-left-color: #667eea;
}

.stat-card.success {
  border-left-color: #4CAF50;
}

.stat-card.warning {
  border-left-color: #FFC107;
}

.stat-icon {
  font-size: 2.5rem;
}

.stat-content h3 {
  margin: 0 0 0.5rem 0;
  font-size: 0.9rem;
  color: #666;
  text-transform: uppercase;
}

.stat-value {
  font-size: 1.8rem;
  font-weight: bold;
  display: block;
  color: #333;
}

.stat-change {
  font-size: 0.8rem;
  color: #666;
}

.stat-change.positive {
  color: #4CAF50;
}

.filters-card {
  grid-area: filters;
  background: white;
  padding: 1.5rem;
  border-radius: 15px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.filters-card h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.filters {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.filter-select {
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
}

.filter-btn {
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
}

.filter-btn.active,
.filter-btn:hover {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.grades-table-card {
  grid-area: table;
  background: white;
  border-radius: 15px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  overflow: hidden;
}

.table-header {
  padding: 1.5rem;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.table-header h3 {
  margin: 0;
  color: #333;
}

.btn {
  padding: 0.5rem 1rem;
  border: 1px solid #667eea;
  border-radius: 8px;
  background: white;
  color: #667eea;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
}

.btn:hover {
  background: #667eea;
  color: white;
}

.btn-outline {
  background: transparent;
  border: 1px solid #667eea;
  color: #667eea;
}

.table-container {
  overflow-x: auto;
}

.grades-table {
  width: 100%;
  border-collapse: collapse;
}

.grades-table th {
  background: #f8f9fa;
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #333;
  border-bottom: 1px solid #eee;
}

.grades-table td {
  padding: 1rem;
  border-bottom: 1px solid #f0f0f0;
}

.grade-row:hover {
  background: #f8f9fa;
}

.subject-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.subject-icon {
  font-size: 1.2rem;
}

.grade-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-weight: bold;
  font-size: 0.9rem;
}

.grade-badge.excellent {
  background: #4CAF50;
  color: white;
}

.grade-badge.good {
  background: #FFC107;
  color: black;
}

.teacher-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.teacher-avatar {
  width: 30px;
  height: 30px;
  background: #667eea;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 0.8rem;
  font-weight: bold;
}

.comment {
  display: block;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #666;
}

.chart-card {
  grid-area: chart;
  background: white;
  padding: 1.5rem;
  border-radius: 15px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.chart-card h3 {
  margin: 0 0 1.5rem 0;
  color: #333;
}

.chart-placeholder {
  height: 200px;
  background: #f8f9fa;
  border-radius: 10px;
  display: flex;
  align-items: end;
  justify-content: center;
  padding: 1rem;
  margin-bottom: 1rem;
}

.chart-bars {
  display: flex;
  align-items: end;
  gap: 2rem;
  height: 100%;
}

.chart-bar-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.chart-bar {
  width: 30px;
  background: linear-gradient(to top, #667eea, #764ba2);
  border-radius: 5px 5px 0 0;
  transition: height 0.3s ease;
}

.chart-label {
  font-size: 0.8rem;
  color: #666;
}

.chart-legend {
  display: flex;
  gap: 2rem;
  justify-content: center;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  color: #666;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.legend-color.excellent {
  background: #4CAF50;
}

.legend-color.good {
  background: #FFC107;
}

.legend-color.average {
  background: #FF9800;
}

@media (max-width: 1200px) {
  .content-grid {
    grid-template-columns: 1fr;
    grid-template-areas: 
      "stats"
      "filters"
      "table"
      "chart";
  }
  
  .stats-cards {
    grid-template-columns: 1fr;
  }
}
</style>