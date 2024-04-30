<script setup>
import { ref, computed } from 'vue'
import { Icon } from '@iconify/vue'
import axios from 'axios'
import LangPicker from './LangPicker.vue'
import { refreshTokens } from '../utils/helpers'


const emit = defineEmits(['close', 'reload'])

const props = defineProps(['story_id', 'project_id'])

const file = ref(null)
const fileMsg = ref("Upload file")
const newPhotoUrl = ref(null)


const showLangs = ref(false)


const selected_lang = ref("eng")
const langs = ref([
    {
        lang: "eng",
        title: "",
        description: ""
    }
])

const selected_lang_id = computed(() => {
    return langs.value.findIndex(lang => {
        return lang.lang == selected_lang.value
    })
})



const handleFileUpload = (event) => {
    if (event.target.files[0].size > 5000 * 1024) {
        fileMsg.value = "File is too large"
        return
    }
    file.value = event.target.files[0]
    const reader = new FileReader()

    reader.onload = (e) => {
        newPhotoUrl.value = e.target.result
    }
    reader.readAsDataURL(event.target.files[0])
}


const createBanner = async () => {
    const formData = new FormData()
    formData.append('file', file.value)
    formData.append('langs', JSON.stringify(langs.value))
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
        console.log(error)
    }
}




</script>

<!-- <template>
    <div @click.self="$emit('close')"
        class="modal-wrapper w-screen h-screen absolute z-50 top-0 left-0 backdrop-blur-sm flex justify-center items-center">
        <div class="modal flex flex-col bg-gray-900 text-white p-5 rounded-lg items-center gap-2">
            <h2 class="title">New banner</h2>

            <div class="dropdown relative w-full">
                <button @click="showLangs = !showLangs" class="flex items-center">
                    <div class="duration-300" :style="showLangs ? '' : `transform:rotate(-90deg);`">
                        <Icon icon="iconamoon:arrow-down-2-duotone" />
                    </div>
                    Language:{{ selected_lang }}
                </button>
                <Transition>
                    <div v-if="showLangs"
                        class="dropdown-content flex flex-col gap-2 absolute -left-2 bg-gray-900 p-2 border-white border-2 rounded-lg">
                        <LangPicker :langs="langs" :selected_lang="selected_lang" @closeLangs="showLangs=false" @selectLang="lang => selected_lang = lang"/>
                    </div>
                </Transition>
            </div>
            <div v-for="(lang, i) of langs" :key="i">
                <div v-if="lang.lang == selected_lang" class="flex flex-col gap-2">
                    <input v-model="lang.title" type="text" class="text-input" placeholder="Title">
                    <textarea v-model="lang.description" type="text" class="text-input"
                        placeholder="Description"></textarea>
                </div>
            </div>
            <input type="file" id="fileInput" class="hidden" @change="handleFileUpload" />
            <label for="fileInput" class="button text-center w-72">
                <p class="truncate">{{ file != null ? file.name : fileMsg }}</p>
            </label>

            <button @click="createBanner" class="button">Create</button>
        </div>
    </div>
</template> -->

<template>
    <div @click.self="$emit('close')"
        class="modal-wrapper w-screen h-screen absolute z-50 top-0 left-0 backdrop-blur-sm flex justify-center items-center">
        <section class=" w-3/4 overflow-hidden rounded-xl drop-shadow-xl h-3/4">

            <div class="content bg-white grid grid-cols-1 lg:grid-cols-2 h-full">
                <div class="h-full banner  w-full min-w-full relative text-white flex items-end">
                    <input @input="changed = true" type="file" id="fileInput" class="hidden"
                        @change="handleFileUpload" />
                    <label for="fileInput"
                        class="text-center absolute mx-auto bg-gray-900 bg-opacity-80 h-full w-full flex items-center justify-center text-5xl text-green-400 opacity-0 duration-300 cursor-pointe hover:opacity-100">
                        <Icon icon="mdi:camera" />
                    </label>
                    <img :src="file ? newPhotoUrl : 'http://localhost:3001/images/banners/no-image.jpg'" alt=""
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
                                    class="dropdown-content flex flex-col gap-2 absolute -left-2 bg-gray-900 p-2 border-white border-2 rounded-lg">
                                    <LangPicker :langs="langs" :selected_lang="selected_lang"
                                        @closeLangs="showLangs = false" @selectLang="lang => selected_lang = lang" />
                                </div>
                            </Transition>
                        </div>
                        <h2 class=" font-semibold">Title:</h2>
                        <input class=" bg-slate-100 p-2 rounded-md"
                            v-model="langs[selected_lang_id].title" placeholder="Title">
                        <h2 class=" font-semibold">Description:</h2>
                        <textarea class="text-input text-black w-full bg-slate-100 p-2 rounded-md" rows="20"
                            v-model="langs[selected_lang_id].description"
                            placeholder="Description"></textarea>
                        <button class="button absolute bottom-0 right-0 m-5 w-fit px-5"
                            @click="createBanner(current)">Save</button>
                    </div>
                </div>
            </div>
        </section>
    </div>
</template>