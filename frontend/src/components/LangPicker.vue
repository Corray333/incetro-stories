<script setup>
import { ref, defineProps } from 'vue'

const props = defineProps(["langs"])

const new_lang = ref("")

const addLang = (lang) => {
    if (props.langs.find(l => l.lang == lang)) {
        alert("Language already exists")
        return
    }
    props.langs.push({
        lang: lang,
        title: "",
        description: ""
    })
}

</script>

<template>
    <div class="flex flex-col gap-2">
        <div class="new_lang flex flex-col">
            <input v-model="new_lang" type="text" class="text-input" placeholder="New lang">
            <button @click="if (new_lang != '') addLang(new_lang); new_lang = ''" class="button">Add language</button>
        </div>
        <button v-for="(lang, i) of langs" :key="i" @click="$emit('selectLang', lang.lang); $emit('closeLangs')"
            class="button">{{ lang.lang }}</button>
    </div>

</template>


<style></style>