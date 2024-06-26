<script setup>
import { ref, computed } from 'vue'
import { Icon } from '@iconify/vue'
import axios from 'axios'
import LangPicker from './LangPicker.vue'
import { refreshTokens } from '../utils/helpers'

const btnLoading = ref(false)
const MaxFileSize = 5 * 1024 * 1024

const emit = defineEmits(['close', 'reload'])

const props = defineProps(['story_id', 'project_id', 'newBanner'])


const selected_lang = ref("en")



const showLangs = ref(false)



const selected_lang_id = computed(() => {
    return props.newBanner.langs.findIndex(lang => {
        return lang.lang == selected_lang.value
    })
})



const handleFileUpload = (event) => {
    // max file size is 5MB
    if (event.target.files[0].size > MaxFileSize) {
        alert('File size is too big')
        return
    }
    props.newBanner.file = event.target.files[0]
    const reader = new FileReader()

    reader.onload = (e) => {
        props.newBanner.newPhotoUrl = e.target.result
    }
    reader.readAsDataURL(event.target.files[0])
}


const createBanner = async () => {
    if (!props.newBanner.langs[selected_lang_id.value].title || !props.newBanner.langs[selected_lang_id.value].description || !props.newBanner.file) {
        alert('Please fill in all fields')
        return
    }
    const formData = new FormData()
    formData.append('file', props.newBanner.file)
    formData.append('langs', JSON.stringify(props.newBanner.langs))
    try {
        let url =  `/api/projects/${props.project_id}/banners`
        if (props.story_id) url += `?story_id=${props.story_id}`
        await axios.post(url, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
                'Authorization': localStorage.getItem('Authorization')
            }
        })
        emit('close')
        emit('reload')
    } catch (error) {
        btnLoading.value = false
        alert(error.response.data)
    }
}




</script>

<template>
    <div @click.self="$emit('close')"
        class="modal-wrapper w-screen h-screen fixed z-50 top-0 left-0 backdrop-blur-sm flex justify-center items-center">
        <section class=" w-3/4 overflow-hidden rounded-xl drop-shadow-xl h-3/4">

            <div class="content bg-white grid grid-cols-1 lg:grid-cols-2 grid-rows-2 lg:grid-rows-1 h-full">
                <div class="h-full banner  w-full min-w-full relative text-white flex items-end">
                    <input @input="changed = true" type="file" id="fileInput" class="hidden"
                        @change="handleFileUpload" />
                    <label for="fileInput"
                        class="text-center absolute mx-auto bg-gray-900 bg-opacity-80 h-full w-full flex items-center justify-center text-5xl text-green-400 opacity-0 duration-300 cursor-pointe hover:opacity-100">
                        <Icon icon="mdi:camera" />
                    </label>
                    <img :src="newBanner.file ? newBanner.newPhotoUrl : '/api/images/banners/no-image.jpg'" alt=""
                        class="w-full h-full object-contain duration-300 bg-black">
                </div>
                <div class="story_info w-full flex p-5 gap-5 overflow-y-auto relative">
                    <div class="w-full h-full flex flex-col rounded-xl">
                        <div class="dropdown relative w-full">
                            <button @click="showLangs = !showLangs" class="flex items-center">
                                <div class="duration-300" :style="showLangs ? '' : `transform:rotate(-90deg);`">
                                    <Icon icon="iconamoon:arrow-down-2-duotone" />
                                </div>
                                Language:{{ selected_lang }}
                            </button>
                            <Transition>
                                <div v-if="showLangs"
                                    class="dropdown-content w-full flex flex-col gap-2 absolute left-0 bg-gray-900 p-2 border-white border-2 rounded-lg">
                                    <LangPicker :langs="newBanner.langs" :selected_lang="selected_lang"
                                        @closeLangs="showLangs = false" @selectLang="lang => selected_lang = lang" />
                                </div>
                            </Transition>
                        </div>
                        <h2 class=" font-semibold">Title:</h2>
                        <input class=" bg-slate-100 p-2 rounded-md"
                            v-model="newBanner.langs[selected_lang_id].title" placeholder="Title">
                        <h2 class=" font-semibold">Description:</h2>
                        <textarea class="text-input text-black w-full bg-slate-100 p-2 rounded-md" rows="20"
                            v-model="newBanner.langs[selected_lang_id].description"
                            placeholder="Description"></textarea>
                        <button class="button absolute bottom-0 flex justify-center right-0 m-5 w-fit px-5"
                            @click.once="btnLoading = true; createBanner(current)"><Icon v-if="btnLoading" icon="line-md:loading-loop" /><p v-else>Create banner</p></button>
                    </div>
                </div>
            </div>
        </section>
    </div>
</template>