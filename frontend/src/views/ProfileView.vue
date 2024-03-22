<script setup>
import {useRouter} from 'vue-router'
import { Icon } from '@iconify/vue'

const router = useRouter()

const logout = ()=>{
    document.cookie="Authorization=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"
    document.cookie="Refresh=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"
    router.push('/login')
}

</script>

<template>
    <section>
       <div class="header w-full flex justify-center text-gray-900 relative items-center">
        <h1 class="title">Profile</h1>
        <button @click="logout" class="button w-min absolute right-0 text-2xl"><Icon icon="mdi:logout" /></button>
       </div> 
    </section>
</template>

