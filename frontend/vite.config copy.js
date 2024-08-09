// vite.config.js
import { resolve } from 'path'
import { defineConfig } from 'vite'
import { libInjectCss } from 'vite-plugin-lib-inject-css'

import tailwind from "tailwindcss";
import autoprefixer from "autoprefixer";
export default defineConfig({
  plugins: [
    tailwind, autoprefixer,
    libInjectCss()
  ],
  build: {
    lib: {
      // Could also be a dictionary or array of multiple entry points
      entry: resolve(__dirname, 'src/index.js'),
      name: 'mmp',
      // the proper extensions will be added
      fileName: 'main',
      formats: ['es'],
    },
    rollupOptions: {
      // make sure to externalize deps that shouldn't be bundled
      // into your library
      external: [],
      output: {
        assetFileNames: 'assets/[name][extname]',
        entryFileNames: '[name].js',
        globals: {
        },
      },
    },
  },
})