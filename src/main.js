import { createApp, h } from 'vue'
import './style.css'
import Notes from './components/Notes.vue'
import About from './components/About.vue'
import App from './App.vue'
import {createRouter, createWebHistory} from 'vue-router'

const routes = [
    { path: '/', component: Notes },
    { path: '/about', component: About },
]

const app = createApp({
    render: ()=>h(App)
})

app.use(createRouter({
    history: createWebHistory(),
    routes
}))

app.mount('#app')
