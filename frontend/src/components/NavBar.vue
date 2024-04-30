<script setup>
import { Icon } from '@iconify/vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const logout = () => {
    localStorage.removeItem('Authorization')
    localStorage.removeItem('Refresh')
    router.push('/login')
}


</script>

<template>
    <aside class="p-5 flex flex-row sm:flex-col justify-center sm:justify-start gap-5">
        <router-link to="/home" class="button w-fit text-3xl flex items-center relative group">
            <Icon icon="mdi:home" />
            <p class="tip group-hover:scale-100 text-white">Home</p>
        </router-link>
        <router-link to="/projects" class="button w-fit text-3xl flex items-center relative group">
            <Icon icon="ant-design:project-outlined" />
            <p class="tip group-hover:scale-100 text-white">Projects</p>
        </router-link>
        <router-link to="/profile" class="button w-fit text-3xl flex items-center relative group">
            <Icon icon="mdi:user" />
            <p class="tip group-hover:scale-100 text-white">Profile</p>
        </router-link>
        <button @click="logout" class="button w-fit w-min text-3xl flex items-center relative group">
            <Icon icon="mdi:logout" />
            <p class="tip group-hover:scale-100 text-white">Exit</p>
        </button>
    </aside>
</template>
