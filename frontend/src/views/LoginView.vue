<script setup>
import { Icon } from '@iconify/vue'
import {ref} from 'vue'
import {useRouter} from 'vue-router'

const router = useRouter()
import axios from 'axios'

let loading = ref(false)

const action = ref("Log in")

const username = ref("")
const email = ref("")
const password = ref("")

const passwordRequirements = ref([
    {text: "contain ≥ 8 charachters", valid: false},
    {text: "contain 1 digit", valid: false},
    {text: "contain 1 uppercase letter", valid: false},
    {text: "contain 1 lowercase letter", valid: false},
])

const checkPassword = ()=>{
    passwordRequirements.value[0].valid = password.value.length >= 8
    passwordRequirements.value[1].valid = /\d/.test(password.value)
    passwordRequirements.value[2].valid = /[A-Z]/.test(password.value)
    passwordRequirements.value[3].valid = /[a-z]/.test(password.value)
}
const emailRequirements = ref([
    {text: "contain @ symbol", valid: false},
    {text: "contain . symbol", valid: false},
    {text: "contain only one @ symbol", valid: false},
    {text: "contain at least one character before and after @", valid: false},
    {text: "contain at least two characters after . (dot)", valid: false},
])

const checkEmail = ()=>{
    emailRequirements.value[0].valid = /@/.test(email.value);
    emailRequirements.value[1].valid = /\./.test(email.value);
    emailRequirements.value[2].valid = (email.value.match(/@/g) || []).length === 1;
    emailRequirements.value[3].valid = /^[^@]*@[^@]*$/.test(email.value);
    emailRequirements.value[4].valid = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/.test(email.value) && email.value.split('.').pop().length >= 2;
}

const login = async ()=>{
    loading.value = true
    try {
        let valid = true
        for (let i = 0; i < passwordRequirements.value.length; i++) valid = valid & passwordRequirements.value[i].valid
        for (let i = 0; i < emailRequirements.value.length; i++) valid = valid & emailRequirements.value[i].valid
        if (!valid) {
            alert("Введите корректные данные!")
            loading.value = false
            return
        }
        if (action.value == "Log in"){
            let {data} = await axios.post('http://localhost:3001/api/users/login', {
                email: email.value,
                password: password.value,
            })
            localStorage.setItem('Authorization', data.authorization)
            localStorage.setItem('Refresh', data.refresh)
            router.push('/home')
        } else if (action.value == "Sign up"){
            let {data} = await axios.post('http://localhost:3001/api/users/signup', {
                username: username.value,
                email: email.value,
                password: password.value,
            })
            localStorage.setItem('Authorization', data.authorization)
            localStorage.setItem('Refresh', data.refresh)
            router.push('/projects')
        }
        else console.log("Invalid action")
    } catch (error) {
        loading.value = false
        alert("Вход не выполнен")   
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
            <div class="input_box relative flex items-center group">
                <input @input="checkEmail" class="text-input" v-model="email" type="text" name="" id="" placeholder="email">
                <div class="input_tip flex items-center gap-2 group-focus-within:scale-100">
                    <Icon icon="mdi:bulb" class="text-green-400 text-4xl" />
                    <ul class="list-disc list-inside">
                        <strong>Email must:</strong>
                        <li v-for="(req, i) of emailRequirements" :key="i" :class="req.valid?'text-green-400':''">{{ req.text }}</li>
                    </ul>
                </div>

            </div>
            <div class="input_box relative flex items-center group">
                <input @input="checkPassword" class="text-input" v-model="password" type="password" name="" id="" placeholder="password">
                <div class="input_tip flex items-center gap-2 group-focus-within:scale-100">
                    <Icon icon="mdi:bulb" class="text-green-400 text-2xl" />
                    <ul class="list-disc list-inside">
                        <strong>Password must:</strong>
                        <li v-for="(req, i) of passwordRequirements" :key="i" :class="req.valid?'text-green-400':''">{{ req.text }}</li>
                    </ul>
                </div>
            </div>
            <button type="button" class="button uppercase flex justify-center hover:text-gray-900 text-sm" @click="login">
                <Icon v-if="loading" icon="line-md:loading-loop" />
                <p v-else>{{ action }}</p>
            </button>
        </section>
    </section>
</template>

