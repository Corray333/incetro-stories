<script setup>
import { Icon } from '@iconify/vue'
import {ref} from 'vue'

const props = defineProps(['story'])

const emit = defineEmits(['newInStory'])

const current = ref(0)
const hovered = ref(false)



</script>

<template>
    <section @mouseenter="hovered = true" @mouseleave="hovered = false" class="w-full relative overflow-hidden rounded-xl">
        <div v-if="story.banners.length > 1" class="button_container absolute z-20 h-full flex items-center">
            <button @click="current = (current-1)%story.banners.length; current < 0 ? current = story.banners.length-1 : pass" type="button" class="arrow-button w-fit"><Icon icon="iconamoon:arrow-left-2-bold" /></button>
        </div>
        <div v-if="story.banners.length > 1" class="button_container  absolute z-20 right-0 h-full flex items-center">
            <button @click="current = (current+1)%story.banners.length" type="button" class="arrow-button w-fit"><Icon icon="iconamoon:arrow-right-2-bold" /></button>
        </div>
        <div class="slider-track flex w-full duration-300" :style="`transform: translateX(calc(-100% * ${current}))`">
            <div class="banner h-96  w-full min-w-full relative text-white flex items-end" v-for="(banner, i) of story.banners" :key="i" :style="`z-index:${i == current ? '100':'initial'}`">
                <img :src="`http://localhost:3001/images/banners/banner${banner.id}.png`" alt="" :style="`transform: ${hovered ? 'scale(1.2)' : 'none'};`" class="w-full h-full object-cover absolute duration-300" style="z-index:0;">
                <div class="banner-info relative z-10 p-5 pt-10 h-full flex flex-col justify-end">
                    <h2 class="font-bold text-xl cursor-pointer hover:text-green-400 w-fit duration-300" @click="$emit('pickStory', story)">{{ banner.title }}</h2>
                    <p class="line-clamp-3" v-html="banner.description"></p>
                </div>
            </div>
        </div>
    </section>
</template>

<style>



</style>