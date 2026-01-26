<script lang="ts" setup>

import type {Issue} from "@/api/issues/types.ts";

const props = defineProps<{
  issue: Issue,
  density?: null | 'default' | 'comfortable' | 'compact'
}>();

const color = computed(() => {
  switch (props.issue.type.toLowerCase()) {
    case 'epic':
      return 'purple';
    case 'story':
      return 'blue';
    case 'task':
      return 'green';
    case 'bug':
      return 'red';
      case 'idea':
      return 'yellow';
    default:
      return 'primary';
  }
});

const icon = computed(() => {
  switch (props.issue.type.toLowerCase()) {
    case 'epic':
      return 'mdi-book-open-page-variant';
    case 'story':
      return 'mdi-book-variant';
    case 'task':
      return 'mdi-checkbox-marked-circle-outline';
    case 'bug':
      return 'mdi-bug-outline';
    case 'idea':
      return 'mdi-lightbulb-outline';
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
    {{ props.issue.type.replace('_', ' ').toUpperCase() }}
  </v-chip>
</template>

<style scoped>

</style>
