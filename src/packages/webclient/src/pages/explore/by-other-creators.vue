<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";

useHead({
  title: 'Explore Projects by Other Creators - FancySpaces',
  meta: [
    {
      name: 'description',
      content: 'Discover and explore projects created by other talented creators on FancySpaces.'
    }
  ]
});

const spaces = ref<Space[]>();

onMounted(async () => {
  spaces.value = [];
  spaces.value.push(await getSpace("orbisguard"));
  spaces.value.push(await getSpace("orbismines"));
  spaces.value.push(await getSpace("wiflowscoreboard"));
});

</script>

<template>
  <v-container>
    <v-row class="my-4" justify="center">
      <v-col>
        <h1 class="text-h3 text-center">Projects by other creators</h1>
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
  </v-container>
</template>

<style scoped>

</style>
