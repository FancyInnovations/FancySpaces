<script lang="ts" setup>

import AppHeader from "@/components/AppHeader.vue";
import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";

const spaces = ref<Space[]>();

onMounted(async () => {
  spaces.value = [];
  spaces.value.push(await getSpace("fn"));
  spaces.value.push(await getSpace("fh"));
  spaces.value.push(await getSpace("fd"));
  spaces.value.push(await getSpace("fc"));
  spaces.value.push(await getSpace("fa"));
  spaces.value.push(await getSpace("fv"));
});

</script>

<template>
  <AppHeader/>

  <v-container
    width="40%"
  >
    <v-row class="mt-16 mb-8" justify="center">
      <v-col>
        <h1 class="text-h3 text-center">All projects</h1>
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
