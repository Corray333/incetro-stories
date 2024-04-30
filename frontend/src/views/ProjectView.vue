<script setup>
import NewBanner from '../components/NewBanner.vue'
import { ref, onBeforeMount } from 'vue'
import { Icon } from '@iconify/vue'
import StoryCardEdit from '../components/StoryCardEdit.vue'
import StoryModalEdit from '../components/StoryModalEdit.vue'
import { jwtDecode } from "jwt-decode"
import axios from 'axios'
import { useRoute } from 'vue-router'
import {  refreshTokens } from '../utils/helpers'

const route = useRoute()
const project_id = route.params.project_id


const newBannerOpened = ref(false)
const storyId = ref(null)

const closeNewBanner = () => {
    newBannerOpened.value = false
    storyId.value = null
}

const pickStory = (id) => {
    storyId.value = id
    newBannerOpened.value = true
}



const stories = ref([])
const projects = ref([])


const loadContent = async () => {
    try {
        let { data } = await axios.get(`/api/projects/${project_id}/stories`, {
            headers: {
                'Authorization': localStorage.getItem('Authorization'),
            }
        })

        stories.value = data.stories
    } catch (error) {
        if (error.response.status == 401) {
            await refreshTokens()
            loadContent()
        }
        else console.log(error)
    }
}

const loadProject = async () => {
    try {
        let { data } = await axios.get(`/api/projects/${project_id}`, {
            headers: {
                'Authorization': localStorage.getItem('Authorization'),
            }
        })
        projects.value = data.projects
    } catch (error) {
        console.log(error)
        if (error.response.status == 401) {
            await refreshTokens()
            loadProject()
        }
        else console.log(error)
    }
}


const storyPick = ref(null)

onBeforeMount(() => {
    loadProject()
    loadContent()
})

</script>

<template>
    <section class="flex flex-col gap-5">
        <transition>
            <StoryModalEdit v-if="storyPick" :story="storyPick" :project_id="project_id" @reload="loadContent; storyPick=null" @close="storyPick = null" />
        </transition>
        <Transition>
            <NewBanner v-if="newBannerOpened" :story_id="storyId" :project_id="project_id" @close="closeNewBanner" @reload="loadContent" />
        </Transition>
        <div class="header w-full flex gap-2 justify-center text-gray-900 relative items-center">
            <img :src="projects[0]?.cover" alt="" class=" w-12 h-12 object-cover rounded-full">
            <h1 class="title">{{ projects[0]?.name }}</h1>
            <button class="button w-fit" @click="newBannerOpened = true">
                <Icon icon="mdi:plus" />
            </button>
        </div>
        <p class="text-center">{{ projects[0]?.description }}</p>
        <div class="stories grid grid-cols-4 gap-5 w-full">
            <StoryCardEdit @new-in-story="pickStory" v-for="(story, i) of stories" :key="i" :story="story"
                @pickStory="storyPick = story" />
        </div>
    </section>
</template>

<style></style>
