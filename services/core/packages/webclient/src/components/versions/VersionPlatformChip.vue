<script lang="ts" setup>

import type {SpaceVersion} from "@/api/versions/types.ts";

const props = defineProps<{
  version: SpaceVersion,
  density?: null | 'default' | 'comfortable' | 'compact'
}>();

const color = computed(() => {
  switch (props.version.platform.toLowerCase()) {
    case 'bukkit':
      return 'light-blue';
    case 'spigot':
      return 'orange';
    case 'paper':
      return 'light-blue-lighten-3';
    case 'purpur':
      return 'purple';
    case 'folia':
      return 'pink';
    case 'bungeecord':
      return 'cyan';
    case 'waterfall':
      return 'cyan-lighten-3';
    case 'velocity':
      return 'teal';
    case 'hytale_plugin':
      return 'green';
    case 'executable':
      return 'grey-3';
    default:
      return 'primary';
  }
});

</script>

<template>
  <v-chip
    :color="color"
    :density="props.density || 'default'"
    rounded
    variant="tonal"
  >
    {{ props.version.platform.replace('_', ' ').toUpperCase() }}
  </v-chip>
</template>

<style scoped>

</style>
