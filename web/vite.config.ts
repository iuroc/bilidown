import { defineConfig } from 'vite'
import { resolve } from 'path'

export default defineConfig({
    build: {
        rollupOptions: {
            input: {
                work: resolve(import.meta.dirname, 'src/windows/work/index.html'),
                download: resolve(import.meta.dirname, 'src/windows/download/index.html'),
                settings: resolve(import.meta.dirname, 'src/windows/settings/index.html'),
                // 继续添加更多窗口页面
            }
        },
        outDir: 'dist',
        emptyOutDir: true
    }
})