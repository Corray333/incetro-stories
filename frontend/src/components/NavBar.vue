<script setup>
import { Icon } from '@iconify/vue'


</script>

<template>
    <aside class="p-5 flex flex-col gap-5">
        <router-link to="/home" class="button text-3xl flex"><Icon icon="mdi:home" /></router-link>
        <router-link to="/profile" class="button text-3xl flex"><Icon icon="mdi:user" /></router-link>
        <router-link to="/create" class="button text-3xl flex"><Icon icon="mdi:plus" /></router-link>
    </aside>
</template>

