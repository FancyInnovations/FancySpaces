<script lang="ts" setup>

import type {Issue} from "@/api/issues/types.ts";

const props = defineProps<{
  issue: Issue,
  density?: null | 'default' | 'comfortable' | 'compact'
}>();

const color = computed(() => {
  switch (props.issue.status.toLowerCase()) {
    case 'backlog':
      return 'yellow';
    case 'planned':
      return 'orange';
    case 'in_progress':
      return 'blue';
    case 'done':
      return 'green';
    case 'closed':
      return 'grey';
    default:
      return 'primary';
  }
});

const icon = computed(() => {
  switch (props.issue.status.toLowerCase()) {
    case 'backlog':
      return 'mdi-timer-sand-empty';
    case 'planned':
      return 'mdi-checkbox-blank-circle-outline';
    case 'in_progress':
      return 'mdi-progress-clock';
    case 'done':
      return 'mdi-check-circle-outline';
    case 'closed':
      return 'mdi-close-circle-outline';
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
    {{ props.issue.status.replace('_', ' ').toUpperCase() }}
  </v-chip>
</template>

<style scoped>

</style>
