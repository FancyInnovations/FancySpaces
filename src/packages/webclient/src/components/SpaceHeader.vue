<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import type {SpaceVersion} from "@/api/versions/types.ts";

const props = defineProps<{
  space?: Space;
  latestVersion?: SpaceVersion;
  downloadCount?: number;
}>();

</script>

<template>
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
        <p class="text-body-2">Created {{ space?.created_at.toLocaleDateString() }}</p>
        <p class="text-body-2 mx-4">-</p>
        <p class="text-body-2">Updated {{ latestVersion?.published_at.toLocaleDateString() || space?.created_at.toLocaleDateString() }}</p>
        <p class="text-body-2 mx-4">-</p>
        <p class="text-body-2">{{ downloadCount }} downloads</p>
      </div>
    </div>

    <div class="d-flex flex-column justify-center">
      <v-btn
        v-if="latestVersion?.files.length != 1"
        :to="`/spaces/${space?.slug}/versions/latest`"
        class="sidebar__mobile"
        color="primary"
        prepend-icon="mdi-download"
        size="large"
        variant="tonal"
      >
        latest
      </v-btn>
      <v-btn
        v-else
        :href="latestVersion?.files[0]?.url"
        class="sidebar__mobile"
        color="primary"
        prepend-icon="mdi-download"
        size="large"
        variant="tonal"
      >
        latest
      </v-btn>
    </div>
  </div>
</template>

<style scoped>
@media (max-width: 960px) {
  .sidebar__mobile {
    display: none;
  }
}
</style>
