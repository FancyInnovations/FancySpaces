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
    case 'discord_forum_post':
      return 'indigo';
    case 'discord_ticket_bot':
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
    case 'discord_forum_post':
      return 'mdi-chat';
    case 'discord_ticket_bot':
      return 'mdi-chat';
    default:
      return 'mdi-help-circle-outline';
  }
});

const externalLink = computed(() => {
  if (!props.issue.external_source) {
    return null;
  }

  switch (props.issue.external_source.toLowerCase()) {
    case 'github':
      return props.issue.extra_fields?.github_url || null;
    case 'discord_forum_post':
      return props.issue.extra_fields?.discord_forum_post_url || null;
    case 'discord_ticket_bot':
      return props.issue.extra_fields?.discord_ticket_url || null;
    default:
      return null;
  }
});

</script>

<template>
  <v-chip
    :color="color"
    :density="props.density || 'default'"
    :href="externalLink"
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
