<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";

useHead({
  title: 'Explore Hytale Plugins - FancySpaces',
  meta: [
    {
      name: 'description',
      content: 'Discover and explore Hytale plugin project spaces on FancySpaces, your hub for innovative Hytale creations.'
    }
  ]
});

const spaces = ref<Space[]>();
const plannedPlugins = ref<Space[]>();

onMounted(async () => {
  spaces.value = [];
  spaces.value.push(await getSpace("fc"));

  plannedPlugins.value = [];
  plannedPlugins.value.push({
    id: "fcw",
    slug: "fancycorewebsite",
    title: "FancyCore Website",
    summary: "Web frontend for FancyCore, allowing server admins to view and manage their server, plugin and player data through a user-friendly interface.",
    description: "",
    categories: ["hytale_plugin"],
    links: [],
    icon_url: "/soon-logo.png",
    status: "approved",
    created_at: new Date(),
    members: [],
  });
  plannedPlugins.value.push({
    id: "fcon",
    slug: "fancyconnect",
    title: "FancyConnect",
    summary: "Connect multiple Hytale servers together. Share chats, player data, and more across a network of servers.",
    description: "",
    categories: ["hytale_plugin"],
    links: [],
    icon_url: "/soon-logo.png",
    status: "approved",
    created_at: new Date(),
    members: [],
  });
  plannedPlugins.value.push({
    id: "fp",
    slug: "fancyplots",
    title: "FancyPlots",
    summary: "Generate plot world(s) and let users claim their own plots. Each plot is protected and can be managed by the owner only.",
    description: "",
    categories: ["hytale_plugin"],
    links: [],
    icon_url: "/soon-logo.png",
    status: "approved",
    created_at: new Date(),
    members: [],
  });
});

</script>

<template>
  <v-container
    width="40%"
  >
    <v-row class="my-4" justify="center">
      <v-col>
        <h1 class="text-h3 text-center">Hytale plugins</h1>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col v-for="space in spaces" :key="space.id" class="mb-4" md="12">
        <SpaceCard
          :space="space"
        />
      </v-col>
    </v-row>

    <v-row class="mt-16 mb-4" justify="center">
      <v-col>
        <h1 class="text-h3 text-center">Planned plugins</h1>
        <p class="text-body-1 text-center mt-2">These plugins are planned to be developed in the future.<br>Not final, and may be cancelled at any time. Stay tuned!</p>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col v-for="space in plannedPlugins" :key="space.id" class="mb-4" md="12">
        <SpaceCard
          :space="space"
          soon
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
