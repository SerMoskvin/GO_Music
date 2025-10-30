<template>
  <div id="app">
    <h1 v-if="!$route.meta.hideNav">Временная навигация</h1>
    <nav v-if="!$route.meta.hideNav">
      <router-link to="/">Главная</router-link> | 
      <router-link to="/info">Информация</router-link> |
      <router-link to="/login">Войти</router-link> |
      <router-link to="/register">Регистрация</router-link>
    </nav>
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const $route = useRoute()
</script>

<style>
body {
  margin: 0;
  font-family: Arial, sans-serif;
}
nav {
  padding: 1rem;
  background: #f0f0f0;
}
nav a {
  margin: 0 1rem;
  text-decoration: none;
  color: #333;
}
</style>