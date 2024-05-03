<script setup>
import { ref } from 'vue'
import { Icon } from '@iconify/vue'
import axios from 'axios'

const emits = defineEmits(['reload'])

const MaxFileSize = 5 * 1024 * 1024

const btnLoading = ref(false)

const props = defineProps(['project'])


const file = ref(null)

const handleFileUpload = (event) => {
    if (event.target.files[0].size > MaxFileSize) {
        alert('File is too big')
        return
    }
    props.project.file = event.target.files[0]
    const reader = new FileReader()

    reader.onload = (e) => {
        props.project.cover = e.target.result
    }
    reader.readAsDataURL(event.target.files[0])
}

const createProject = async ()=>{
    if (!props.project.name || !props.project.description || !props.project.file) {
        alert('Please fill in all fields')
        return
    }
    const formData = new FormData()
    formData.append('cover', props.project.file)
    formData.append('data', JSON.stringify({
        name: props.project.name,
        description: props.project.description
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
        btnLoading.value = false
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
            <img :src="project.cover" alt="photo" class="w-48 h-48 rounded-full object-cover border-white border-8">
        </div>
        <div class="info flex flex-col gap-2">
            <div>
                <p>Name:</p>
                <input v-model="project.name" type="text" class="text-input">
            </div>
            <div>
                <p>Description:</p>
                <input v-model="project.description" type="text" class="text-input">
            </div>
            <button @click.once="btnLoading = true; createProject()" class="button flex justify-center"><Icon v-if="btnLoading" icon="line-md:loading-loop" /><p v-else>Create project</p> </button>
        </div>
    </section>
</template>


<style></style>