<script setup>
import { onBeforeMount } from 'vue'
import NavBar from './components/NavBar.vue'
import { refreshTokens } from './utils/helpers'


onBeforeMount(()=>{
  if (localStorage.getItem('Refresh') != null){
    refreshTokens()
  }
})

</script>

<template>
  <div class="sm:flex block">
    <NavBar/>
  <div class="bg-slate-100 w-full min-h-screen p-5 text-gray-900">
    <router-view ></router-view>
  </div>
  </div>

</template>

<style scoped>

</style>
