<script lang="ts" setup>

import {getDownloadCountForSpace, getSpacesOfCreator} from "@/api/spaces/spaces.ts";
import type {Space} from "@/api/spaces/types.ts";
import type {User} from "@/api/auth/types.ts";

const props = defineProps<{
  user?: User;
}>();

const spaces = ref<Space[]>([]);
const totalDownloads = ref(0);

const bio = computed(() => {
  if (!props.user?.metadata || !props.user?.metadata['public_biography']) {
    return "No biography set.";
  }

  return props.user?.metadata['public_biography'];
});

const profilePicture = computed(() => {
  if (!props.user?.metadata || !props.user?.metadata['public_profile_picture']) {
    return "/na-logo.png";
  }

  return props.user?.metadata['public_profile_picture'];
});

watch(() => props.user, async () => {
  await updateSpaces();
}, {immediate: true});

async function updateSpaces() {
  if (!props.user) {
    return;
  }

  spaces.value = await getSpacesOfCreator(props.user!.id);

  for (let sp of spaces.value) {
    totalDownloads.value += await getDownloadCountForSpace(sp.id);
  }
}

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
            :src="profilePicture"
            alt="User Profile Picture"
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
            <h1>{{ user?.name }}</h1>
            <p class="text-body-1 mt-2">{{ bio }}</p>
          </div>

          <div class="d-flex mt-2 text-grey-lighten-1">
            <p class="text-body-2">Joined at: {{ user?.created_at.toLocaleDateString() }}</p>
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
