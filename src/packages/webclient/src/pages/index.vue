<script lang="ts" setup>

import AppHeader from "@/components/AppHeader.vue";
import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";

const spaces = ref<Space[]>();

onMounted(async () => {
  spaces.value = [];
  spaces.value.push(await getSpace("fn"));
  spaces.value.push(await getSpace("fc"));
  spaces.value.push(await getSpace("fh"));
  spaces.value.push(await getSpace("fa"));
});
</script>

<template>
  <AppHeader/>

  <v-container
    width="60%"
  >
    <v-row class="mt-16" justify="center">
      <v-col>
        <v-img
          src="/banner.png"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col>
        <p class="text-h5 text-center text-primary">Download platform for products by FancyInnovations</p>
      </v-col>
    </v-row>

    <v-row class="mt-16" justify="center">
      <v-col md="4">
        <v-card
          class="card__border"
          color="#29152550"
          elevation="12"
          height="100%"
          rounded="xl"
          to="/explore/minecraft-plugins"
        >
          <v-card-title class="mt-2">
            <div class="d-flex align-center">
              <v-icon class="mr-2" color="green lighten-1">mdi-minecraft</v-icon>
              <p>Minecraft plugins</p>
            </div>
          </v-card-title>

          <v-card-text>
            FancyNpcs, FancyHolograms and more fancy Minecraft plugins made by FancyInnovations.
          </v-card-text>
        </v-card>
      </v-col>

      <v-col md="4">
        <v-card
          class="card__border"
          color="#29152550"
          elevation="12"
          height="100%"
          rounded="xl"
          to="/explore/hytale-plugins"
        >
          <v-card-title class="mt-2">
            <v-icon class="mr-2" color="blue lighten-1">mdi-sword-cross</v-icon>
            Hytale plugins
          </v-card-title>

          <v-card-text>
            Our first Hytale plugin FancyCore is a must-have for every Hytale server.
          </v-card-text>
        </v-card>
      </v-col>

      <v-col md="4">
        <v-card
          class="card__border"
          color="#29152550"
          elevation="12"
          height="100%"
          rounded="xl"
          to="/explore/other-projects"
        >
          <v-card-title class="mt-2">
            <v-icon class="mr-2" color="purple lighten-1">mdi-cube-outline</v-icon>
            Other projects
          </v-card-title>

          <v-card-text>
            Explore our other exciting projects and tools like FancyAnalytics or FancyVerteiler.
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <h1 class="text-h4 mt-16 text-center">Popular projects</h1>
      </v-col>
    </v-row>

    <v-row v-if="spaces && spaces.length > 0" justify="center">
      <v-col md="8">
        <v-carousel
          :show-arrows="false"
          cycle
          height="fit-content"
          interval="5000"
        >
          <v-carousel-item v-for="space in spaces" :key="space.id">
            <SpaceCard
              :space="space"
              :with-badge="true"
              class="pb-10"
            />
          </v-carousel-item>
        </v-carousel>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
