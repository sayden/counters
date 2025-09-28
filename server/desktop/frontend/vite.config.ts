import { defineConfig } from 'vite'
import { resolve } from 'path'
import react from '@vitejs/plugin-react'
import prism from 'vite-plugin-prismjs'
import { execSync } from 'child_process';

function astroBuildBeforeVitePlugin() {
  return {
    name: 'astro-build-before-vite',
    buildStart() {
      execSync('bun run build', { cwd: "../documentation", stdio: 'inherit' }); // adjust as necessary
    }
  };
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    prism({
      languages: ['json', 'javascript'],
      plugins: ['line-numbers', 'show-language', 'copy-to-clipboard'],
      theme: 'tomorrow-night'
    }),
    // astroBuildBeforeVitePlugin(),
  ],
  build: {
    rollupOptions: {
      input: {
        index: resolve(__dirname, 'index.html'),
        builder: resolve(__dirname, 'src/views/builder/index.html'),
      }
    }
  }
})
