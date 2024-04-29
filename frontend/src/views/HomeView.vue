<script setup>
import {ref, onBeforeMount} from 'vue'
import StoryCard from '../components/StoryCard.vue';
import StoryModal from '../components/StoryModal.vue';
import axios from 'axios'
import {  refreshTokens } from '../utils/helpers'

const stories = ref([])
const storyPick = ref(null)

const pickStory = (story)=>{
    storyPick.value = story
}




const loadContent = async ()=>{
    try {
        let {data} = await axios.get('http://localhost:3001/api/stories', {
            headers:{
                'Authorization': localStorage.getItem('Authorization'),
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

onBeforeMount(()=>{
    loadContent()
})

</script>

<template>
    <transition>
        <StoryModal :story="storyPick" v-if="storyPick" @close="storyPick = null"/>
    </transition>
    <section class="w-full flex flex-col items-center gap-5">
        <h1 class="title">Home</h1>
        <div class="stories grid grid-cols-4 gap-5 w-full">
            <StoryCard v-for="(story, i) of stories" :key="i" :story="story" @pickStory="pickStory"/>
        </div>
    </section>
</template>

