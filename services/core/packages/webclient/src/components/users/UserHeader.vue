<script lang="ts" setup>

import {getDownloadCountForSpace, getSpacesOfCreator} from "@/api/spaces/spaces.ts";
import type {Space} from "@/api/spaces/types.ts";

const props = defineProps<{
  userID?: string;
}>();

const spaces = ref<Space[]>([]);
const totalDownloads = ref(0);

onMounted(async () => {
  spaces.value = await getSpacesOfCreator(props.userID!);

  for (let sp of spaces.value) {
    totalDownloads.value += await getDownloadCountForSpace(sp.id);
  }
});

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
            :href="`/users/${userID}`"
            alt="User Profile Picture"
            height="100"
            max-height="100"
            max-width="100"
            min-height="100"
            min-width="100"
            src="/na-logo.png"
            width="100"
          />
        </div>

        <div class="mx-4 d-flex flex-column justify-space-between flex-grow-1">
          <div>
            <h1>{{ userID }}</h1>
            <p class="text-body-1 mt-2">The user bio belongs here</p>
          </div>

          <div class="d-flex mt-2 text-grey-lighten-1">
            <p class="text-body-2">Joined at: {{ '2026-01-01' }}</p>
            <p class="text-body-2 mx-4">-</p>
            <p class="text-body-2">{{ spaces.length }} spaces</p>
            <p class="text-body-2 mx-4">-</p>
            <p class="text-body-2">{{ totalDownloads }} downloads</p>
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
