<script lang="ts" setup>

import type {Issue} from "@/api/issues/types.ts";

const props = defineProps<{
  issue: Issue,
  density?: null | 'default' | 'comfortable' | 'compact'
}>();

const color = computed(() => {
  switch (props.issue.priority.toLowerCase()) {
    case 'low':
      return 'green';
    case 'medium':
      return 'blue';
    case 'high':
      return 'red';
    case 'critical':
      return 'red';
    default:
      return 'primary';
  }
});

const icon = computed(() => {
  switch (props.issue.priority.toLowerCase()) {
    case 'low':
      return 'mdi-arrow-down-circle-outline';
    case 'medium':
      return 'mdi-minus-circle-outline';
    case 'high':
      return 'mdi-arrow-up-circle-outline';
    case 'critical':
      return 'mdi-alert-circle-outline';
    default:
      return 'mdi-help-circle-outline';
  }
});

</script>

<template>
  <v-chip
    :color="color"
    :density="props.density || 'default'"
    :prepend-icon="icon"
    rounded
    variant="tonal"
  >
    {{ props.issue.priority.replace('_', ' ').toUpperCase() }}
  </v-chip>
</template>

<style scoped>

</style>
