<script setup>
import { Icon } from '@iconify/vue'
import {useRouter} from 'vue-router'

const router = useRouter()
const logout = ()=>{
    document.cookie="Authorization=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"
    document.cookie="Refresh=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"
    router.push('/login')
}


</script>

<template>
    <aside class="p-5 flex flex-col gap-5">
        <router-link to="/home" class="button text-3xl flex"><Icon icon="mdi:home" /></router-link>
        <router-link to="/gallery" class="button text-3xl flex"><Icon icon="mdi:pencil" /></router-link>
        <router-link to="/profile" class="button text-3xl flex"><Icon icon="mdi:user" /></router-link>
        <button @click="logout" class="button w-min text-3xl"><Icon icon="mdi:logout" /></button>
    </aside>
</template>

