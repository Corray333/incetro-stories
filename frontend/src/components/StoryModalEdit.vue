<script setup>
import { Icon } from '@iconify/vue'
import { ref, onBeforeMount, computed } from 'vue'
import axios from 'axios'
import LangPicker from './LangPicker.vue'
import { refreshTokens } from '../utils/helpers'

import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';



const props = defineProps(['story', 'project_id'])
const emits = defineEmits(['reload'])
const expires_at = ref(new Date())

const current = ref(0)
const selected_lang = ref("")
const showLangs = ref(false)

const selected_lang_id = computed(() => {
    return props.story.banners[current.value].langs.findIndex(lang => {
        return lang.lang == selected_lang.value
    })
})


const user = ref({})

const loadUserInfo = async () => {
    try {
        let uid = props.story.creator
        let { data } = await axios.get(`/api/users/${uid}`, {
            headers: {
                'Authorization': localStorage.getItem('Authorization'),
            }
        })

        user.value = data.user
    } catch (error) {
        if (error.status == 401) {
            await refreshTokens()
            loadUserInfo()
        }
        else console.log(error)
    }
}

onBeforeMount(() => {
    loadUserInfo()
    selected_lang.value = props.story.banners[current.value].langs[0].lang
    expires_at.value = new Date(props.story.expires_at * 1000)
})

const updateBanner = async (id) => {
    const formData = new FormData()
    formData.append('file', file.value)
    formData.append('expires_at', Math.floor(expires_at.value.getTime() / 1000))
    formData.append('banner', JSON.stringify(props.story.banners[id]))

    try {
        await axios.put(`/api/banners`,
            formData, {
            headers: {
                'Authorization': localStorage.getItem('Authorization'),
            }
        })

        emits('reload')
    } catch (error) {
        if (error.status == 401) {
            await refreshTokens()
            updateBanner()
        }
        else console.log(error)
    }
}

const newPhotoUrl = ref('')
const file = ref(null)

const handleFileUpload = (event) => {
    console.log('test')
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


</script>

<template>
    <div @click.self="$emit('close')"
        class="modal-wrapper w-screen h-screen fixed z-50 top-0 left-0 backdrop-blur-sm flex justify-center items-center">
        <section class=" w-3/4 overflow-hidden rounded-xl drop-shadow-xl h-3/4">

            <div class="content bg-white grid grid-cols-1 lg:grid-cols-2 h-full">
                <div class="slider-container overflow-hidden relative h-full">
                    <div v-if="story.banners.length > 1"
                        class="button_container absolute z-20 h-full flex items-center">
                        <button
                            @click="file = null; selected_lang = story.banners[(current - 1) % story.banners.length].langs[0].lang; current = (current - 1) % story.banners.length; current < 0 ? current = story.banners.length - 1 : pass"
                            type="button" class="arrow-button w-fit">
                            <Icon icon="iconamoon:arrow-left-2-bold" />
                        </button>
                    </div>
                    <div v-if="story.banners.length > 1"
                        class="button_container  absolute z-20 right-0 h-full flex items-center">
                        <button
                            @click="file = null; selected_lang = story.banners[(current + 1) % story.banners.length].langs[0].lang; current = (current + 1) % story.banners.length"
                            type="button" class="arrow-button w-fit">
                            <Icon icon="iconamoon:arrow-right-2-bold" />
                        </button>
                    </div>
                    <div class="slider-track flex w-full duration-300 h-full"
                        :style="`transform: translateX(calc(-100% * ${current}))`">
                        <div class="h-full banner  w-full min-w-full relative text-white flex items-end"
                            v-for="(banner, i) of story.banners" :key="i"
                            :style="`z-index:${i == current ? '100' : 'initial'}`">
                            <input @input="changed = true" type="file" id="fileInput" class="hidden"
                                @change="handleFileUpload" />
                            <label for="fileInput"
                                class="text-center absolute mx-auto bg-gray-900 bg-opacity-80 h-full w-full flex items-center justify-center text-5xl text-green-400 opacity-0 duration-300 cursor-pointe hover:opacity-100">
                                <Icon icon="mdi:camera" />
                            </label>
                            <img :src="file ? newPhotoUrl : banner.media_url" alt=""
                                class="w-full h-full object-contain duration-300 bg-black">
                        </div>
                    </div>
                </div>
                <div class="story_info w-full flex p-5 gap-5 overflow-y-auto relative">
                    <div class="w-full h-full flex flex-col rounded-xl">
                        <i class="opacity-50 text-xs">Creator: {{ user.username }}</i>
                        <i class="opacity-50 text-xs">Created at:{{ new Date(story.created_at * 1000) }}</i>
                        <span class=" gap-2 items-center">
                            <i class="opacity-50 text-xs">Expires at:</i>
                            <div class="w-fit">
                                <VueDatePicker v-model="expires_at"/>
                            </div>
                        </span>
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
                                    <LangPicker :langs="story.banners[current].langs" :selected_lang="selected_lang"
                                        @closeLangs="showLangs = false" @selectLang="lang => selected_lang = lang" />
                                </div>
                            </Transition>
                        </div>
                        <h2 class=" font-semibold">Title:</h2>
                        <input class=" bg-slate-100 p-2 rounded-md"
                            v-model="story.banners[current].langs[selected_lang_id].title" placeholder="Title">
                        <h2 class=" font-semibold">Description:</h2>
                        <textarea class="text-input text-black w-full bg-slate-100 p-2 rounded-md" rows="20"
                            v-model="story.banners[current].langs[selected_lang_id].description"
                            placeholder="Description"></textarea>
                        <button class="button absolute bottom-0 right-0 m-5 w-fit px-5"
                            @click="updateBanner(current)">Save</button>
                    </div>
                </div>
            </div>
        </section>
    </div>
</template>

<style></style>