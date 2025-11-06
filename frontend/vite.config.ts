import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: [
            { find: '@', replacement: '/src' },
            { find: '@tabler/icons-react', replacement: '@tabler/icons-react/dist/esm/icons/index.mjs' }
        ],
    },
})