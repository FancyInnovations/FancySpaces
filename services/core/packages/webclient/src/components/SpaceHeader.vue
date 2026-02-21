<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";

const props = defineProps<{
  space?: Space;
}>();

</script>

<template>
  <Card
    class="mb-4"
    color="#4a2f0033"
  >
    <v-card-text>
      <div class="d-flex justify-space-between">
        <div class="d-flex flex-column justify-center">
          <v-img
            :href="`/spaces/${space?.slug}`"
            :src="space?.icon_url || '/logo.png'"
            alt="Space Icon"
            height="100"
            max-height="100"
            max-width="100"
            min-height="100"
            min-width="100"
            width="100"
          />
        </div>

        <div class="mx-4 d-flex flex-column justify-space-between flex-grow-1">
          <div>
            <h1>{{ space?.title }}</h1>
            <p class="text-body-1 mt-2">{{ space?.summary }}</p>
          </div>

          <div class="d-flex mt-2 text-grey-lighten-1">
            <p class="text-body-2">Creator: {{ space?.creator }}</p>
            <p class="text-body-2 mx-4">-</p>
            <p class="text-body-2">Created {{ space?.created_at.toLocaleDateString() }}</p>
            <slot name="metadata">
            </slot>
          </div>
        </div>

        <div class="d-flex flex-column justify-center">
          <slot name="quick-actions">
          </slot>
        </div>
      </div>
    </v-card-text>
  </Card>
</template>

<style scoped>
@media (max-width: 960px) {
  .sidebar__mobile {
    display: none;
  }
}
</style>
