import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  base: '/patient/',
  resolve: {
    alias: { '@': '/src' },
  },
  css: {
    transformer: 'postcss',
  },
  build: {
    cssMinify: 'esbuild',
  },
  server: {
    port: 5174,
    host: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
