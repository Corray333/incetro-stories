<script setup>
import NavBar from './components/NavBar.vue'

</script>

<template>
  <div class="flex">
    <NavBar/>
  <div class="bg-slate-100 w-full min-h-screen p-5 text-gray-900">
    <router-view ></router-view>
  </div>
  </div>

</template>

<style scoped>

</style>
