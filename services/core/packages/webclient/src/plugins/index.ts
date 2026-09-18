/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 */

// Types
import type { App } from 'vue'
import { createHead } from '@vueuse/head'
import router from '../router'
import pinia from '../stores'
// Plugins
import vuetify from './vuetify'

export function registerPlugins (app: App) {
  const head = createHead()

  app
    .use(vuetify)
    .use(router)
    .use(pinia)
    .use(head)
}
