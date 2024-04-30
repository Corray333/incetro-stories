<script setup>
import { ref } from 'vue'
import { Icon } from '@iconify/vue'
import axios from 'axios'

const emits = defineEmits(['reload'])


const name = ref('New project')
const description = ref('New project description')
const cover = ref("https://via.placeholder.com/150")

const file = ref(null)

const handleFileUpload = (event) => {
    if (event.target.files[0].size > 5000*1024) {
        alert('File is too big')
        return
    }
    file.value = event.target.files[0]
    const reader = new FileReader()

    reader.onload = (e) => {
        cover.value = e.target.result
    }
    reader.readAsDataURL(event.target.files[0])
}

const createProject = async ()=>{
    const formData = new FormData()
    if (file.value != null) formData.append('cover', file.value)
    formData.append('data', JSON.stringify({
        name: name.value,
        description: description.value
    }))

    try{
        await axios.post( `/api/projects`, formData, {
            headers: {
                'Content-Type':'multipart/form-data',
                'Authorization': localStorage.getItem('Authorization')
            }
        })
        emits('reload')
    } catch (error) {
        console.log(error)
    }
}

</script>

<template>
    <section class="bg-gray-900 text-white p-5 rounded-xl flex flex-col sm:flex-row gap-5 items-center">
        <div class="profile_photo relative">
            <input type="file" id="fileInput" class="hidden" @change="handleFileUpload" />
            <label for="fileInput"
                class="text-center absolute mx-auto bg-gray-900 bg-opacity-80 h-full w-full rounded-full flex items-center justify-center text-5xl text-green-400 opacity-0 duration-300 cursor-pointer border-green-400 border-8 hover:opacity-100">
                <Icon icon="mdi:user" />
            </label>
            <img :src="cover" alt="photo" class="w-48 h-48 rounded-full object-cover border-white border-8">
        </div>
        <div class="info flex flex-col gap-2">
            <div>
                <p>Name:</p>
                <input v-model="name" type="text" class="text-input">
            </div>
            <div>
                <p>Description:</p>
                <input v-model="description" type="text" class="text-input">
            </div>
            <button @click="createProject" class="button">Create project</button>
        </div>
    </section>
</template>


<style></style>