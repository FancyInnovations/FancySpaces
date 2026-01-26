<script lang="ts" setup>

import type {Issue} from "@/api/issues/types.ts";

const props = defineProps<{
  issue: Issue,
  density?: null | 'default' | 'comfortable' | 'compact'
}>();

</script>

<template>
  <v-chip
    :density="props.density || 'default'"
    color="primary"
    prepend-icon="mdi-sign-text"
    rounded
    variant="tonal"
  >
    #{{ props.issue.id }}
  </v-chip>
</template>

<style scoped>

</style>
