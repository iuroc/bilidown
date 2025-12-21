import { defineConfig } from 'vite'

export default defineConfig({
    clearScreen: false,
    base: './',
    resolve: {
        alias: {
            '@': '/src'
        }
    },
    server: {
        open: true,
        proxy: {
            '/api': `http://localhost:8080`
        }
    }
})