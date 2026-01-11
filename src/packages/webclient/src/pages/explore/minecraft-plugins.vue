<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";

useHead({
  title: 'Explore Minecraft Plugins - FancySpaces',
  meta: [
    {
      name: 'description',
      content: 'Discover and explore Minecraft plugin project spaces on FancySpaces, your hub for innovative Minecraft creations.'
    }
  ]
});

const spaces = ref<Space[]>();

onMounted(async () => {
  spaces.value = [];
  spaces.value.push(await getSpace("fn"));
  spaces.value.push(await getSpace("fh"));
  spaces.value.push(await getSpace("fd"));
});
</script>

<template>
  <v-container
    width="40%"
  >
    <v-row class="my-4" justify="center">
      <v-col>
        <h1 class="text-h3 text-center">Minecraft plugins</h1>
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
