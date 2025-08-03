import path from "path"
import tailwindcss from "@tailwindcss/vite"
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

const timestamp = new Date().getTime()/1000
// https://vite.dev/config/
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        entryFileNames: `[name]` + timestamp + `.js`,
        chunkFileNames: `[name]` + timestamp + `.js`,
        assetFileNames: `[name]` + timestamp + `.[ext]`
      }
    }
  },
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      '/api/': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      }
    }
  },
})

