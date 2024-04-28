<script setup>
import { Icon } from '@iconify/vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const logout = () => {
    document.cookie = "Authorization=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"
    document.cookie = "Refresh=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"
    router.push('/login')
}


</script>

<template>
    <aside class="p-5 flex flex-col gap-5">
        <router-link to="/home" class="button text-3xl flex items-center relative group">
            <Icon icon="mdi:home" />
            <p class="tip group-hover:scale-100 text-white">Home</p>
        </router-link>
        <router-link to="/projects" class="button text-3xl flex items-center relative group">
            <Icon icon="ant-design:project-outlined" />
            <p class="tip group-hover:scale-100 text-white">Projects</p>
        </router-link>
        <router-link to="/profile" class="button text-3xl flex items-center relative group">
            <Icon icon="mdi:user" />
            <p class="tip group-hover:scale-100 text-white">Profile</p>
        </router-link>
        <button @click="logout" class="button w-min text-3xl flex items-center relative group">
            <Icon icon="mdi:logout" />
            <p class="tip group-hover:scale-100 text-white">Exit</p>
        </button>
    </aside>
</template>
