<script setup>
import {ref, onBeforeMount} from 'vue'
import {useRoute} from 'vue-router'
import StoryCard from '../components/StoryCard.vue'
import StoryModal from '../components/StoryModal.vue'
import axios from 'axios'
import { getCookie, refreshTokens } from '../utils/helpers'

const route = useRoute()

const stories = ref([])
const storyPick = ref(null)

const pickStory = (story)=>{
    storyPick.value = story
}



const loadContent = async ()=>{
    try {
        let {data} = await axios.get(`http://localhost:3001/stories?creator=${route.params.uid}`, {
            headers:{
                'Authorization': getCookie('Authorization'),
            }
        })
        
        stories.value = data.stories
    } catch (error) {
        if (error.response.status == 401){
            await refreshTokens()
            loadContent()
        }
        else console.log(error)
    }
}

const user = ref({})

const loadUserInfo = async ()=>{
    try {
        let {data} = await axios.get(`http://localhost:3001/api/users/${route.params.uid}`, {
            headers:{
                'Authorization': getCookie('Authorization'),
            }
        })
        
        user.value = data.user
    } catch (error) {
        if (error.response.status == 401){
            await refreshTokens()
            loadUserInfo()
        }
        else console.log(error)
    }
}

onBeforeMount(()=>{
    loadUserInfo()
    loadContent()
})

</script>

<template>
    <transition>
        <StoryModal :story="storyPick" v-if="storyPick" @close="storyPick = null"/>
    </transition>
    <section class="w-full flex flex-col items-center gap-5">
        <img :src="user.avatar" alt="" class="w-40 h-40 object-cover rounded-full">
        <h1 class="title">{{ user.username }}</h1>
        <div class="stories grid grid-cols-4 gap-5 w-full">
            <StoryCard v-for="(story, i) of stories" :key="i" :story="story" @pickStory="pickStory"/>
        </div>
    </section>
</template>

