<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import {useUserStore} from "@/stores/user.ts";

const userStore = useUserStore();

const isLoggedIn = ref(false);

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
  isLoggedIn.value = await userStore.isAuthenticated;

  spaces.value = [];
  spaces.value.push(await getSpace("fc"));

  if (isLoggedIn.value) {
    spaces.value.push(await getSpace("fancyplots"));
    spaces.value.push(await getSpace("fancyconnect"));
    spaces.value.push(await getSpace("fancyaudits"));
    spaces.value.push(await getSpace("fancycorewebsite"));
    spaces.value.push(await getSpace("fancyshops"));
    spaces.value.push(await getSpace("citypass"));
    spaces.value.push(await getSpace("cityquests"));
    spaces.value.push(await getSpace("cityshops"));
  }
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
  </v-container>
</template>

<style scoped>

</style>
