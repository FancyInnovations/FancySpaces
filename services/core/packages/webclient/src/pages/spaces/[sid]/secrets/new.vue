<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {createSecret} from "@/api/secrets/secrets.ts";
import SpaceHeader from "@/components/SpaceHeader.vue";
import {useNotificationStore} from "@/stores/notifications.ts";

const router = useRouter();
const notificationStore = useNotificationStore();

const isLoggedIn = computed(() => {
  return localStorage.getItem("fs_api_key") !== null;
});

const space = ref<Space>();

const key = ref('');
const value = ref('');
const description = ref('');

onMounted(async () => {
  const spaceID = (useRoute().params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!isLoggedIn || !space.value.secrets_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  useHead({
    title: `${space.value.title} - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || 'Create a new secret in this space on FancySpaces.'
      }
    ]
  });
});

async function createNewSecret() {
  await createSecret(space.value!.id, key.value, value.value, description.value);

  key.value = '';
  value.value = '';
  description.value = '';

  notificationStore.info("Secret created successfully");

  await router.push(`/spaces/${space.value?.slug}/secrets`);
}

</script>

<template>
  <v-container width="90%">
    <v-row>
      <v-col class="flex-grow-0 pa-0">
        <SpaceSidebar
          :space="space"
        />
      </v-col>

      <v-col>
        <SpaceHeader :space="space">
          <template #quick-actions>
            <v-btn
              :to="`/spaces/${space?.slug}/secrets`"
              class="sidebar__mobile"
              color="primary"
              size="large"
              variant="tonal"
            >
              View Secrets
            </v-btn>
          </template>
        </SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <h1 class="text-center">Create new secret for {{ space?.title }}</h1>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-text-field
          v-model="key"
          color="primary"
          hide-details
          label="Key"
          required
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-textarea
          v-model="value"
          color="primary"
          hide-details
          label="Secret"
          rows="3"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-textarea
          v-model="description"
          color="primary"
          hide-details
          label="Description"
          rows="6"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-btn
          color="primary"
          @click="createNewSecret()"
        >
          Create Secret
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.grey-border-color {
  border-color: rgba(0, 0, 0, 0.8);
}
</style>
