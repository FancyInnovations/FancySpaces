<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";

useHead({
  title: 'Explore Other Projects - FancySpaces',
  meta: [
    {
      name: 'description',
      content: 'Discover and explore various project spaces on FancySpaces, your hub for innovative creations.'
    }
  ]
});

const spaces = ref<Space[]>();

onMounted(async () => {
  spaces.value = [];
  spaces.value.push(await getSpace("fa"));
  spaces.value.push(await getSpace("fv"));
});

</script>

<template>
  <v-container
    width="40%"
  >
    <v-row class="my-4" justify="center">
      <v-col>
        <h1 class="text-h3 text-center">Other projects</h1>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col v-for="space in spaces" :key="space.id" class="mb-4" md="12">
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
