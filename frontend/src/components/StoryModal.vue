<script setup>
import { Icon } from '@iconify/vue'
import {ref, onBeforeMount} from 'vue'
import axios from 'axios'

function getCookie(name) {
 const value = `; ${document.cookie}`;
 const parts = value.split(`; ${name}=`);
 if (parts.length === 2) return parts.pop().split(';').shift();
}

const props = defineProps(['story'])

const current = ref(0)


const user = ref({})

const loadUserInfo = async ()=>{
    try {
        let uid = props.story.creator
        let {data} = await axios.get(`http://localhost:3001/users/${uid}`, {
            headers:{
                'Authorization': getCookie('Authorization'),
            }
        })
        
        user.value = data.user
    } catch (error) {
        if (error.status == 401){
            await refreshTokens()
            loadUserInfo()
        }
        else console.log(error)
    }
}

onBeforeMount(()=>{
    loadUserInfo()
})


</script>

<template>
    <div  @click.self="$emit('close')" class="modal-wrapper w-screen h-screen absolute z-50 top-0 left-0 backdrop-blur-sm flex justify-center items-center">
        <section class=" w-3/4 overflow-hidden rounded-xl drop-shadow-xl h-3/4">

            <div class="content bg-white grid grid-cols-3 h-full">
                <div class="slider-container overflow-hidden relative h-full col-span-2">
                    <div v-if="story.banners.length > 1" class="button_container absolute z-20 h-full flex items-center">
                        <button @click="current = (current-1)%story.banners.length; current < 0 ? current = story.banners.length-1 : pass" type="button" class="arrow-button w-fit"><Icon icon="iconamoon:arrow-left-2-bold" /></button>
                    </div>
                    <div v-if="story.banners.length > 1" class="button_container  absolute z-20 right-0 h-full flex items-center">
                        <button @click="current = (current+1)%story.banners.length" type="button" class="arrow-button w-fit"><Icon icon="iconamoon:arrow-right-2-bold" /></button>
                    </div>
                    <div class="slider-track flex w-full duration-300 h-full" :style="`transform: translateX(calc(-100% * ${current}))`">
                        <div class="h-full banner  w-full min-w-full relative text-white flex items-end" v-for="(banner, i) of story.banners" :key="i" :style="`z-index:${i == current ? '100':'initial'}`">
                            <img :src="`http://localhost:3001/images/banners/banner${banner.id}.png`" alt="" class="w-full h-full object-contain duration-300 bg-black">
                        </div>
                    </div>
                </div>
                <div class="story_info flex p-5 gap-5 overflow-y-auto">
                    <img :src="user.avatar" alt="" class="w-16 h-16 object-cover rounded-full">
                    <div class=" h-full flex flex-col rounded-xl">
                        <i class="opacity-50 text-xs">{{ user.username }}</i>
                        <h2 class="font-bold text-xl">{{ story.banners[current].name }}</h2>
                        <p v-html="story.banners[current].description"></p>
                    </div>
                </div>
            </div>
        </section>
    </div>
</template>

<style>



</style>