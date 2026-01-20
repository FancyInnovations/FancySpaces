<script lang="ts" setup>

import {mapCategoryToDisplayname, mapLinkToDisplayname, mapLinkToIcon, type Space} from "@/api/spaces/types.ts";

const props = defineProps<{
  space?: Space
}>();
</script>

<template>
  <v-navigation-drawer
    class="sidebar__mobile sidebar__background ma-4"
    elevation="12"
    rounded="xl"
  >
    <v-list>
      <v-list-item>
        <v-list-item-title class="text-h6 font-weight-bold">{{ space?.title }}</v-list-item-title>
        <v-list-item-subtitle>{{ mapCategoryToDisplayname(space?.categories[0]) }}</v-list-item-subtitle>
      </v-list-item>

      <v-divider class="mt-2"/>

      <v-list-item
        :to="`/spaces/${space?.slug}`"
        exact
        link
        prepend-icon="mdi-information-slab-circle-outline"
        title="Information"
      />

      <v-list-item
        :to="`/spaces/${space?.slug}/versions`"
        link
        prepend-icon="mdi-file-download-outline"
        title="Versions"
      />

      <v-list-item
        :to="`/spaces/${space?.slug}/stats`"
        disabled
        link
        prepend-icon="mdi-chart-box-outline"
        title="Stats"
      />

      <v-divider />
      <v-list-subheader>External links</v-list-subheader>

      <v-list-item
        v-for="(link) in space?.links" :key="link.name"
        :href="link.url"
        :prepend-icon="mapLinkToIcon(link.name)"
        :title="mapLinkToDisplayname(link.name)"
        link
        target="_blank"
      />

    </v-list>
  </v-navigation-drawer>
</template>

<style scoped>
.sidebar__background {
  max-height: calc(100vh - 96px);
  background-color: #29152550 !important;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

@media (max-width: 960px) {
  .sidebar__mobile {
    display: none;
  }
}
</style>
