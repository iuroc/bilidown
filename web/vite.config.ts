import { defineConfig } from 'vite'
import { resolve } from 'path'

export default defineConfig({
    base: './',
    build: {
        target: 'esnext',
        rollupOptions: {
            input: {
                work: resolve(import.meta.dirname, 'src/windows/work/index.html'),
                download: resolve(import.meta.dirname, 'src/windows/download/index.html'),
                settings: resolve(import.meta.dirname, 'src/windows/settings/index.html'),
            }
        },
        outDir: 'dist',
        emptyOutDir: true
    }
})