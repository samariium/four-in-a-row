import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// ✅ clean Vite configuration
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
  },
})
