<script lang="ts" setup>
import { ref, onBeforeMount } from 'vue'
import { Icon } from '@iconify/vue'
import NewProjectModal from '../components/NewProjectModal.vue'
import ProjectCard from '../components/ProjectCard.vue'
import axios from 'axios'
import {  refreshTokens } from '../utils/helpers'

const showNewProjectModal = ref(false)

const projects = ref([])

const newProject = ref({
    name: 'New project',
    description: 'New project description',
    cover: "https://via.placeholder.com/150",
    file: null
})

const loadProjects = async () => {
    try {
        let { data } = await axios.get( `/api/projects`, {
            headers: {
                'Authorization': localStorage.getItem('Authorization'),
            }
        })
        projects.value = data.projects
    } catch (error) {
        console.log(error)
        if (error.response.status == 401) {
            await refreshTokens()
            loadProjects()
        }
        else console.log(error)
    }
}

onBeforeMount(() => {
    loadProjects()
})


</script>

<template>
    <section class="flex flex-col gap-5">
        <Transition>
            <div v-if="showNewProjectModal" @click.self="showNewProjectModal = false"
                class="wrapper w-full h-screen absolute left-0 z-40 top-0 backdrop-blur-md flex justify-center items-center">
                <NewProjectModal @reload="loadProjects(); showNewProjectModal = false" :project="newProject"/>
            </div>
        </Transition>
        <div class="flex justify-center gap-2">
            <h1 class="title">Projects</h1>
            <button @click="showNewProjectModal = true"  class="button w-fit">
                <Icon icon="mdi:plus" />
            </button>
        </div>
        <section class="projects grid grid-cols-1 lg:grid-cols-3 xl:grid-cols-4 md:grid-cols-2  gap-5">
            <ProjectCard v-for="(project, i) in projects" :key="i" :project="project" />
        </section>

    </section>
</template>


<style></style>