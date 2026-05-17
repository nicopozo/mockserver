import { fileURLToPath, URL } from 'node:url'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import vuetify, { transformAssetUrls } from 'vite-plugin-vuetify'

let appVersion = '0.0.0'
try {
  appVersion = readFileSync(resolve(__dirname, '../VERSION'), 'utf-8').trim()
} catch {
  console.warn('⚠️ Warning: Could not read VERSION file, falling back to 0.0.0')
}

// https://vite.dev/config/
export default defineConfig(({ mode }) => ({
  base: mode === 'production' ? '/mock-service/admin/' : '/admin/',
  plugins: [
    vue({
      template: { transformAssetUrls }
    }),
    vueDevTools(),
    vuetify({ autoImport: true }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  define: {
    __APP_VERSION__: JSON.stringify(appVersion)
  }
}))
