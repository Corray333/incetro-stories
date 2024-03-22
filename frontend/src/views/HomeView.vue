<script setup>
import {ref, onBeforeMount} from 'vue'
import StoryCard from '../components/StoryCard.vue';
import axios from 'axios'

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
        alert(error)
    }
}

const loadContent = async ()=>{
    try {
        let {data} = await axios.get('http://localhost:3001/stories', {
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
        else alert(error)
    }
}

onBeforeMount(()=>{
    loadContent()
})

</script>

<template>
    <section class="w-full flex flex-col items-center">
        <h1 class="title">Home</h1>
        <div class="stories grid grid-cols-4 gap-5 w-full">
            <StoryCard v-for="(story, i) of stories" :key="i" :story="story" />
        </div>
    </section>
</template>

