<template>
  <div class="assessment-list">
    <div class="header">
      <h2>Оценки студентов</h2>
      <button @click="showCreateForm = true" class="btn btn-primary">
        Добавить оценку
      </button>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    
    <table v-else class="data-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Урок</th>
          <th>Студент</th>
          <th>Тип задания</th>
          <th>Оценка</th>
          <th>Дата оценки</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="assessment in assessments" :key="assessment.id">
          <td>{{ assessment.id }}</td>
          <td>{{ assessment.lesson_id }}</td>
          <td>{{ assessment.student_id }}</td>
          <td>{{ assessment.task_type }}</td>
          <td>{{ assessment.grade }}</td>
          <td>{{ assessment.assessment_date }}</td>
          <td class="actions">
            <button @click="viewAssessment(assessment.id)" class="btn btn-sm btn-info">
              Просмотр
            </button>
            <button @click="editAssessment(assessment.id)" class="btn btn-sm btn-warning">
              Редакт.
            </button>
            <button @click="deleteAssessment(assessment.id)" class="btn btn-sm btn-danger">
              Удалить
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Модальное окно создания -->
    <AssessmentForm 
      v-if="showCreateForm"
      :assessment="null"
      @save="handleCreate"
      @cancel="showCreateForm = false"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAssessmentStore } from '@stores/assessmentStore'
import AssessmentForm from '@/components/assessments/AssessmentForm.vue'
import { storeToRefs } from 'pinia'

const router = useRouter()
const assessmentStore = useAssessmentStore()
const showCreateForm = ref(false)

const { assessments, loading, error } = storeToRefs(assessmentStore)

onMounted(() => {
  assessmentStore.fetchAll()
})

const viewAssessment = (id: number) => {
  router.push(`/assessments/${id}`)
}

const editAssessment = (id: number) => {
  router.push(`/assessments/${id}/edit`)
}

const deleteAssessment = async (id: number) => {
  if (confirm('Вы уверены что хотите удалить эту оценку?')) {
    await assessmentStore.deleteAssessment(id)
  }
}

const handleCreate = async (data: any) => {
  await assessmentStore.createAssessment(data)
  showCreateForm.value = false
}
</script>

<style scoped>
.assessment-list {
  padding: 2rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 0.75rem;
  border: 1px solid #ddd;
  text-align: left;
}

.data-table th {
  background-color: #f5f5f5;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.btn-primary { background: #007bff; color: white; }
.btn-info { background: #17a2b8; color: white; }
.btn-warning { background: #ffc107; color: black; }
.btn-danger { background: #dc3545; color: white; }
.btn-sm { padding: 0.25rem 0.5rem; font-size: 0.875rem; }

.loading, .error {
  padding: 2rem;
  text-align: center;
}

.error {
  color: #dc3545;
}
</style>