<script lang="ts" setup>

import AppHeader from "@/components/AppHeader.vue";
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

onMounted(async () => {
  spaces.value = [];
  spaces.value.push(await getSpace("fc"));
});

</script>

<template>
  <AppHeader/>

  <v-container
    width="40%"
  >
    <v-row class="mt-16 mb-8" justify="center">
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
  </v-container>
</template>

<style scoped>

</style>
