import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/health': 'http://localhost:2500',
      '/send': 'http://localhost:2500',
      '/receivers': 'http://localhost:2500',
      '/groups': 'http://localhost:2500',
      '/templates': 'http://localhost:2500',
      '/swagger': 'http://localhost:2500',
    },
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
  },
})
