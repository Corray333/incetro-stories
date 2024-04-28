<script setup>

const props = defineProps(['project'])

</script>

<template>
    <router-link :to="`/projects/${project.id}`">
        <article class="flex flex-col items-center text-white relative group cursor-pointer">
            <img :src="project.cover" alt="" class=" aspect-square object-cover rounded-xl">
            <div class="gradient bottom-0   absolute h-1/2 w-full rounded-xl duration-300 group-hover:h-full"></div>
            <h1 class="title absolute bottom-0 duration-300 group-hover:bottom-1/2 group-hover:text-4xl">{{ project.name }}</h1>
        </article>
    </router-link>
</template>


<style></style>