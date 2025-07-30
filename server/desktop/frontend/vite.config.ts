import { defineConfig } from 'vite'
import { resolve } from 'path'
import react from '@vitejs/plugin-react'
import prism from 'vite-plugin-prismjs'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), prism({
    languages: ['json', 'javascript'],
    plugins: ['line-numbers', 'show-language', 'copy-to-clipboard'],
    theme: 'tomorrow-night'
  })],
  build: {
    rollupOptions: {
      input: {
        index: resolve(__dirname, 'index.html'),
        builder: resolve(__dirname, 'builder.html'),
      }
    }
  }
})
