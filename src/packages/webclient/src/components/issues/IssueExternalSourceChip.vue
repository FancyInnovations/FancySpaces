<script lang="ts" setup>

import type {Issue} from "@/api/issues/types.ts";

const props = defineProps<{
  issue: Issue,
  density?: null | 'default' | 'comfortable' | 'compact'
}>();

const color = computed(() => {
  if (!props.issue.external_source) {
    return 'primary';
  }

  switch (props.issue.external_source.toLowerCase()) {
    case 'github':
      return 'grey';
    case 'discord':
      return 'indigo';
    default:
      return 'primary';
  }
});

const icon = computed(() => {
  if (!props.issue.external_source) {
    return 'mdi-help-circle-outline';
  }

  switch (props.issue.external_source.toLowerCase()) {
    case 'github':
      return 'mdi-github';
    case 'discord':
      return 'mdi-chat';
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
    <p
      v-if="props.issue.external_source"
    >{{ props.issue.external_source.replace('_', ' ').toUpperCase() }}
    </p>
    <p v-else>UNKNOWN</p>
  </v-chip>
</template>

<style scoped>

</style>
