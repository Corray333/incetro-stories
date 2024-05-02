<script setup>
import { Icon } from '@iconify/vue'
import { ref, onBeforeMount } from 'vue'
import { jwtDecode } from "jwt-decode"
import axios from 'axios'
import {  refreshTokens } from '../utils/helpers'


const avatarUrl = ref('')

const changed = ref(false)


const user = ref({})

const loadUserInfo = async () => {
    try {
        let uid = jwtDecode(localStorage.getItem('Authorization')).id
        console.log(uid)
        let { data } = await axios.get( `/api/users/${uid}`, {
            headers: {
                'Authorization': localStorage.getItem('Authorization'),
            }
        })

        user.value = data.user
    } catch (error) {
        if (error.response.status == 401) {
            await refreshTokens()
            loadUserInfo()
        }
        else console.log(error)
    }
}

onBeforeMount(() => {
    loadUserInfo()
})

const file = ref(null)

const handleFileUpload = (event) => {
    console.log('test')
    if (event.target.files[0].size > 500 * 1024) {
        fileMsg.value = "File is too large"
        return
    }
    file.value = event.target.files[0]
    const reader = new FileReader()

    reader.onload = (e) => {
        avatarUrl.value = e.target.result
    }
    reader.readAsDataURL(event.target.files[0])
}

const saveChanges = async () => {
    const formData = new FormData()
    if (file.value != null) formData.append('avatar', file.value)
    formData.append('username', user.value.username)

    try {
        let url =  `/api/users/` + jwtDecode(localStorage.getItem('Authorization')).id
        await axios.put(url, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
                'Authorization': localStorage.getItem('Authorization')
            }
        })
        location.reload()
    } catch (error) {
        console.log(error)
    }
}

</script>

<template>
    <section class="flex flex-col gap-5 items-center">
        <h1 class="title">Profile</h1>
        <div class="profile_card flex flex-col md:flex-row items-center gap-5 bg-gray-900 rounded-xl p-5">
            <div class="profile_photo relative">
                <input @input="changed = true" type="file" id="fileInput" class="hidden" @change="handleFileUpload" />
                <label for="fileInput"
                    class="text-center absolute mx-auto bg-gray-900 bg-opacity-80 h-full w-full rounded-full flex items-center justify-center text-5xl text-green-400 opacity-0 duration-300 cursor-pointer border-green-400 border-8 hover:opacity-100">
                    <Icon icon="mdi:user" />
                </label>
                <img :src="file ? avatarUrl : user.avatar" alt=""
                    class="w-48 h-48 rounded-full object-cover border-white border-8">
            </div>
            <div class="profile_info flex flex-col gap-2">
                <input @input="changed = true" type="text" v-model="user.username" class="text-input">
                <input @input="changed = true" type="text" v-model="user.email" class="text-input opacity-75" disabled>
                <button @click="saveChanges" class="button" :class="changed ? '' : 'disabled'">Save</button>
            </div>
        </div>
    </section>
</template>

<style></style>
