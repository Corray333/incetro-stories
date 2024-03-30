<script setup>
import { onBeforeMount } from 'vue'
import axios from 'axios'
import NavBar from './components/NavBar.vue'

function getCookie(name) {
 const value = `; ${document.cookie}`;
 const parts = value.split(`; ${name}=`);
 if (parts.length === 2) return parts.pop().split(';').shift();
}

const refreshTokens = async()=>{
    try {
        let {data} = await axios.get('http://localhost:3001/users/refresh', {
            headers:{
                'Refresh': getCookie('Refresh'),
            }
        })

        document.cookie = `Authorization=${data.authorization};`
        document.cookie = `Refresh=${data.refresh};`
    } catch (error) {
        if (error.response.status == 401){
            await refreshTokens()
            loadContent()
        }
        else console.log(error)
    }
}
onBeforeMount(()=>{
  refreshTokens()
})

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
