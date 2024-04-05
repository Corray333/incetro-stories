<script setup>
import NewBanner from '../components/NewBanner.vue'
import {ref,onBeforeMount} from 'vue'
import { Icon } from '@iconify/vue'
import StoryCardEdit from '../components/StoryCardEdit.vue'
import StoryModalEdit from '../components/StoryModalEdit.vue'
import { jwtDecode } from "jwt-decode"
import axios from 'axios'


const newBannerOpened = ref(false)
const storyId = ref(null)

const closeNewBanner = ()=>{
    newBannerOpened.value = false
    storyId.value = null
}

const pickStory = (id)=>{
    storyId.value = id
    newBannerOpened.value = true
}



const stories = ref([])

function getCookie(name) {
 const value = `; ${document.cookie}`;
 const parts = value.split(`; ${name}=`);
 if (parts.length === 2) return parts.pop().split(';').shift();
}

const refreshTokens = async()=>{
    try {
        let {data} = await axios.get('http://localhost:3001/users/refresh', {
            headers:{
                'Refresh': getCookie('Refresh'),
            }
        })

        document.cookie = `Authorization=${data.authorization};`
        document.cookie = `Refresh=${data.refresh};`
    } catch (error) {
        console.log(error)
    }
}

const loadContent = async ()=>{
    try {
        let uid = jwtDecode(getCookie('Authorization')).id
        console.log(uid)
        let {data} = await axios.get(`http://localhost:3001/stories?creator=${uid}`, {
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

const storyPick = ref(null)

onBeforeMount(()=>{
    loadContent()
})

</script>

<template>
    <section class="flex flex-col gap-5">
        <transition>
            <StoryModalEdit v-if="storyPick" :story="storyPick" @close="storyPick = null"/>
        </transition>
        <Transition>
            <NewBanner v-if="newBannerOpened" :story_id="storyId" @close="closeNewBanner" @reload="loadContent"/>
        </Transition>
       <div class="header w-full flex gap-2 justify-center text-gray-900 relative items-center">
        <h1 class="title">Your stories</h1>
        <button class="button w-fit" @click="newBannerOpened = true"><Icon icon="mdi:plus" /></button>
       </div> 
       <div class="stories grid grid-cols-4 gap-5 w-full">
            <StoryCardEdit @new-in-story="pickStory" v-for="(story, i) of stories" :key="i" :story="story" @pickStory="storyPick=story"/>
        </div>
    </section>
</template>

<style>



</style>

