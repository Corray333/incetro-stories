<script setup>
import { ref, onBeforeMount } from 'vue'
import { Icon } from '@iconify/vue'
import axios from 'axios'
import LangPicker from './LangPicker.vue'
import { getCookie, refreshTokens } from '../utils/helpers'

const emit = defineEmits(['close', 'reload'])

const props = defineProps(['story_id', 'project_id'])

const file = ref(null)
const fileMsg = ref("Upload file")

const showLangs = ref(false)


const selected_lang = ref("eng")
const langs = ref([
    {
        lang: "eng",
        title: "",
        description: ""
    }
])



const handleFileUpload = (event) => {
    if (event.target.files[0].size > 500 * 1024) {
        fileMsg.value = "File is too large"
        return
    }
    file.value = event.target.files[0]
}


const createBanner = async () => {
    const formData = new FormData()
    formData.append('file', file.value)
    formData.append('langs', JSON.stringify(langs.value))
    console.log(getCookie("Authorization"))
    try {
        let url = `http://localhost:3001/api/projects/${props.project_id}/banners`
        if (props.story_id) url += `?story_id=${props.story_id}`
        await axios.post(url, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
                'Authorization': getCookie("Authorization")
            }
        })
        emit('close')
        emit('reload')
    } catch (error) {
        console.log(error)
    }
}




</script>

<template>
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


</template>