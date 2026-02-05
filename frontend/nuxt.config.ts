// https://nuxt.com/docs/api/configuration/nuxt-config

import { ENV_DEVELOPMENT } from './constant/env'

export default defineNuxtConfig({
  // Make ~/@ resolve from project root to keep custom modules in root directories.
  srcDir: '.',
  appDir: '.',
  css: ['~/assets/css/tailwind.css'],
  modules: ['@nuxtjs/tailwindcss'],
  runtimeConfig: {
    public: {
      env: process.env.ENV ?? ENV_DEVELOPMENT,
    }
  },
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true }
})
