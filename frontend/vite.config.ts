import { defineConfig } from 'vite'

export default defineConfig({
    clearScreen: false,
    base: './',
    resolve: {
        alias: {
            '@': '/src'
        }
    }
})