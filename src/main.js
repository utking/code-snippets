import { createApp, h } from 'vue'
import './style.css'
import Notes from './components/Notes.vue'
import About from './components/About.vue'
import App from './App.vue'
import {createRouter, createWebHistory} from 'vue-router'
import axios, { AxiosInstance } from 'axios'
import VueAxios from 'vue-axios'

const routes = [
    { path: '/', component: Notes },
    { path: '/about', component: About },
]

const apiClient = axios.create({
    baseURL: 'http://localhost:8080',
    headers: {
        'Accept': 'application/json',
        'Content-type': 'application/json',
    },
})

const app = createApp({
    render: ()=>h(App)
})

app.provide('axios', apiClient)

// app.use(VueAxios, axios)
app.use(createRouter({
    history: createWebHistory(),
    routes
}))

app.mount('#app')
