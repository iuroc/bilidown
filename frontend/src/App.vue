<script setup lang="ts">
import { RouterView } from 'vue-router'
import TheHeader from './components/the-header.vue'
</script>

<template>
  <the-header></the-header>
  <router-view class="container py-3"></router-view>
</template>

<style scoped></style>
