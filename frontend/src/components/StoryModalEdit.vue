<script setup>
import { Icon } from '@iconify/vue'
import { ref, onBeforeMount, computed } from 'vue'
import axios from 'axios'
import LangPicker from './LangPicker.vue'
import { getCookie, refreshTokens } from '../utils/helpers'



const props = defineProps(['story', 'project_id'])

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
        let { data } = await axios.get(`http://localhost:3001/api/users/${uid}`, {
            headers: {
                'Authorization': getCookie('Authorization'),
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
})

const updateBanner = async (id) => {
    try {
        await axios.put(`http://localhost:3001/api/projects/${project_id}banners/${props.story.banners[id].id}`, {
            data:JSON.stringify(props.story.banners[id])
        }, {
            headers: {
                'Authorization': getCookie('Authorization'),
            }
        })
        location.reload()
    } catch (error) {
        if (error.status == 401) {
            await refreshTokens()
            loadUserInfo()
        }
        else console.log(error)
    }
}


</script>

<template>
    <div @click.self="$emit('close')"
        class="modal-wrapper w-screen h-screen absolute z-50 top-0 left-0 backdrop-blur-sm flex justify-center items-center">
        <section class=" w-3/4 overflow-hidden rounded-xl drop-shadow-xl h-3/4">

            <div class="content bg-white grid grid-cols-2 h-full">
                <div class="slider-container overflow-hidden relative h-full">
                    <div v-if="story.banners.length > 1"
                        class="button_container absolute z-20 h-full flex items-center">
                        <button
                            @click="current = (current - 1) % story.banners.length; current < 0 ? current = story.banners.length - 1 : pass"
                            type="button" class="arrow-button w-fit">
                            <Icon icon="iconamoon:arrow-left-2-bold" />
                        </button>
                    </div>
                    <div v-if="story.banners.length > 1"
                        class="button_container  absolute z-20 right-0 h-full flex items-center">
                        <button @click="current = (current + 1) % story.banners.length" type="button"
                            class="arrow-button w-fit">
                            <Icon icon="iconamoon:arrow-right-2-bold" />
                        </button>
                    </div>
                    <div class="slider-track flex w-full duration-300 h-full"
                        :style="`transform: translateX(calc(-100% * ${current}))`">
                        <div class="h-full banner  w-full min-w-full relative text-white flex items-end"
                            v-for="(banner, i) of story.banners" :key="i"
                            :style="`z-index:${i == current ? '100' : 'initial'}`">
                            <img :src="banner.media_url" alt=""
                                class="w-full h-full object-contain duration-300 bg-black">
                        </div>
                    </div>
                </div>
                <div class="story_info w-full flex p-5 gap-5 overflow-y-auto relative">
                    <img :src="user.avatar" alt="" class="w-16 h-16 object-cover rounded-full">
                    <div class="w-full h-full flex flex-col rounded-xl">
                        <i class="opacity-50 text-xs">{{ user.username }}</i>
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
                        <input class=" bg-slate-100 p-2 rounded-md" v-model="story.banners[current].langs[selected_lang_id].title" placeholder="Title">
                        <h2 class=" font-semibold">Description:</h2>
                        <textarea class="text-input text-black w-full bg-slate-100 p-2 rounded-md" rows="20"
                            v-model="story.banners[current].langs[selected_lang_id].description" placeholder="Description"></textarea>
                        <button class="button absolute bottom-0 right-0 m-5 w-fit px-5"
                            @click="updateBanner(current)">Save</button>
                    </div>
                </div>
            </div>
        </section>
    </div>
</template>

<style></style>