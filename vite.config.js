import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

process.env = {...process.env, ...loadEnv(process.env.BUILD_MODE ?? "dev", process.cwd())};

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  base: process.env.VITE_FE_BASE_URI ? process.env.VITE_FE_BASE_URI : "/",
  build: {
    outDir: "backend/views/fe"
  }
})
