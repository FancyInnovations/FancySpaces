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
  spaces.value.push(await getSpace("cpp"));
});

</script>

<template>
  <v-container
    width="40%"
  >
    <v-row class="my-4" justify="center">
      <v-col>
        <h1 class="text-h3 text-center">Projects by other creators</h1>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col v-for="space in spaces" :key="space.id" class="mb-4" md="12">
        <SpaceCard
          :space="space"
          with-author
          with-badge
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
