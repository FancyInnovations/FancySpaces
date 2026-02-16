<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {mapEngineKeyToName, type SpaceDatabaseCollection} from "@/api/storage/types.ts";

const props = defineProps<{
  space: Space,
  coll: SpaceDatabaseCollection
  withoutActions?: boolean
}>();

</script>

<template>
  <v-card
    class="card__border"
    color="#19120D33"
    elevation="12"
    rounded="xl"
  >
    <v-card-title class="mt-2">
      {{ coll.name }}
    </v-card-title>

    <v-card-text>
      <p><strong>Database:</strong> {{ coll.database }}</p>
      <p><strong>Engine:</strong> {{ mapEngineKeyToName(coll.engine) }}</p>
      <p><strong>Created at:</strong> {{ coll.created_at.toLocaleString() }}</p>
    </v-card-text>

    <v-card-actions v-if="!withoutActions">
      <v-btn
        :to="`/spaces/${space.slug}/storage/${coll.database}/${coll.name}`"
        color="primary"
        variant="text"
      >
        View Data
      </v-btn>

      <v-btn
        :to="`/spaces/${space.slug}/storage/${coll.database}/${coll.name}/settings`"
        color="primary"
        variant="text"
      >
        Settings
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<style scoped>

</style>
