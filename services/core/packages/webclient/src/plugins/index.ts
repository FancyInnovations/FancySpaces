/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 */

// Plugins
import vuetify from './vuetify'
import pinia from '../stores'
import router from '../router'
import {createHead} from '@vueuse/head'

// Types
import type {App} from 'vue'

export function registerPlugins (app: App) {
  const head = createHead();

  app
    .use(vuetify)
    .use(router)
    .use(pinia)
    .use(head)
}
