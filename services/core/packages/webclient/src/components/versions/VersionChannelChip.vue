<script lang="ts" setup>

import type {SpaceVersion} from "@/api/versions/types.ts";

const props = defineProps<{
  version: SpaceVersion,
  density?: null | 'default' | 'comfortable' | 'compact'
}>();

const color = computed(() => {
  switch (props.version.channel.toLowerCase()) {
    case 'alpha':
      return 'deep-orange-accent-3';
    case 'beta':
      return 'yellow-accent-4';
    case 'release':
      return 'light-green-accent-3';
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
    {{ props.version.channel.replace('_', ' ').toUpperCase() }}
  </v-chip>
</template>

<style scoped>

</style>
