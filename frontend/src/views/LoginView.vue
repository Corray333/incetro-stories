<script setup>
import {ref} from 'vue'
import {useRouter} from 'vue-router'

const router = useRouter()
import axios from 'axios'

const action = ref("Log in")

const username = ref("")
const email = ref("")
const password = ref("")

const login = async ()=>{
    try {
        if (action.value == "Log in"){
            let {data} = await axios.post('http://localhost:3001/users/login', {
                email: email.value,
                password: password.value,
            })
            document.cookie = `Authorization=${data.authorization};`
            document.cookie = `Refresh=${data.refresh};`
            router.push('/home')
        } else if (action.value == "Sign up"){
            let {data} = await axios.post('http://localhost:3001/users/signup', {
                username: username.value,
                email: email.value,
                password: password.value,
            })
            document.cookie = `Authorization=${data.authorization};`
            document.cookie = `Refresh=${data.refresh};`
            router.push('/home')
        }
        else alert("Invalid action")
    } catch (error) {
        alert(error)   
    }
}

</script>

<template>
    <section class="h-screen w-full flex justify-center items-center bg-slate-100 text-white">
        <section class="flex flex-col w-min p-5 gap-2 rounded-xl bg-gray-900 items-center">
            <span class="flex gap-1">
                <p @click="action = 'Log in'" class="cursor-pointer" :class="action == 'Log in' ? 'text-green-400' : 'text-white'">log in</p>
                <p>|</p>
                <p @click="action = 'Sign up'" class="cursor-pointer" :class="action == 'Sign up' ? 'text-green-400': 'text-white'">sign up</p>
            </span>
            <h1 class="font-bold text-xl">{{ action }}</h1>
            <input v-if="action=='Sign up'" v-model="username" class="text-input" type="text" name="" id="" placeholder="username">
            <input class="text-input" v-model="email" type="text" name="" id="" placeholder="email">
            <input class="text-input" v-model="password" type="password" name="" id="" placeholder="password">
            <button type="button" class="button uppercase" @click="login">{{ action }}</button>
        </section>
    </section>
</template>

