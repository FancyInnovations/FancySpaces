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

  if (localStorage.getItem("fs_api_key")) {
    spaces.value.push(await getSpace("fancyplots"));
    spaces.value.push(await getSpace("fancyconnect"));
    spaces.value.push(await getSpace("fancyaudits"));
    spaces.value.push(await getSpace("fancycorewebsite"));
    spaces.value.push(await getSpace("fancyshops"));
    spaces.value.push(await getSpace("citypass"));
    spaces.value.push(await getSpace("cityquests"));
    spaces.value.push(await getSpace("cityshops"));
  }

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
    creator: "admin",
    members: [],
    issue_settings: {
      enabled: false
    },
    release_settings: {
      enabled: false
    },
    maven_repository_settings: {
      enabled: false
    },
    storage_settings: {
      enabled: false
    },
    analytics_settings: {
      enabled: false,
      require_write_key: false
    }
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
    creator: "admin",
    members: [],
    issue_settings: {
      enabled: false
    },
    release_settings: {
      enabled: false
    },
    maven_repository_settings: {
      enabled: false
    },
    storage_settings: {
      enabled: false
    },
    analytics_settings: {
      enabled: false,
      require_write_key: false
    }
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
    creator: "admin",
    members: [],
    issue_settings: {
      enabled: false
    },
    release_settings: {
      enabled: false
    },
    maven_repository_settings: {
      enabled: false
    },
    storage_settings: {
      enabled: false
    },
    analytics_settings: {
      enabled: false,
      require_write_key: false
    }
  });
  plannedPlugins.value.push({
    id: "fancyaudits",
    slug: "fancyaudits",
    title: "FancyAudits",
    summary: "Log player and world related actions to a database for auditing and rollback purposes.",
    description: "",
    categories: ["hytale_plugin"],
    links: [],
    icon_url: "/soon-logo.png",
    status: "approved",
    created_at: new Date(),
    creator: "admin",
    members: [],
    issue_settings: {
      enabled: false
    },
    release_settings: {
      enabled: false
    },
    maven_repository_settings: {
      enabled: false
    },
    storage_settings: {
      enabled: false
    },
    analytics_settings: {
      enabled: false,
      require_write_key: false
    }
  });
});

</script>

<template>
  <v-container>
    <v-row class="my-4" justify="center">
      <v-col>
        <h1 class="text-h3 text-center">Hytale plugins</h1>
      </v-col>
    </v-row>

    <v-row v-for="space in spaces" :key="space.id" justify="center">
      <v-col class="mb-4" md="5">
        <SpaceCard
          :space="space"
          :with-badge="true"
        />
      </v-col>
    </v-row>

    <v-row class="mt-16 mb-4" justify="center">
      <v-col>
        <h1 class="text-h3 text-center">Planned plugins</h1>
        <p class="text-body-1 text-center mt-2">These plugins are planned to be developed in the future.<br>Not final, and may be cancelled at any time. Stay tuned!</p>
      </v-col>
    </v-row>

    <v-row v-for="space in plannedPlugins" :key="space.id" justify="center">
      <v-col class="mb-4" md="5">
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
