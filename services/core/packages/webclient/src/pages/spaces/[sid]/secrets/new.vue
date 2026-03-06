<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {createSecret} from "@/api/secrets/secrets.ts";
import SpaceHeader from "@/components/SpaceHeader.vue";
import {useNotificationStore} from "@/stores/notifications.ts";
import {useUserStore} from "@/stores/user.ts";

const router = useRouter();
const route = useRoute();
const notificationStore = useNotificationStore();
const userStore = useUserStore();

const isLoggedIn = ref(false);

const space = ref<Space>();

const key = ref('');
const value = ref('');
const description = ref('');

onMounted(async () => {
  isLoggedIn.value = await userStore.isAuthenticated;

  const spaceID = (route.params as any).sid as string;
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

    <v-row justify="center">
      <v-col md="8">
        <Card>
          <v-card-title class="mt-2">
            New Secret
          </v-card-title>

          <v-card-text>
            <v-text-field
              v-model="key"
              class="mb-4"
              color="primary"
              hide-details
              label="Key"
              required
            />

            <v-textarea
              v-model="value"
              class="mb-4"
              color="primary"
              hide-details
              label="Secret"
              rows="3"
            />

            <v-textarea
              v-model="description"
              class="mb-4"
              color="primary"
              hide-details
              label="Description"
              rows="6"
            />

            <v-btn
              class="mt-4"
              color="primary"
              @click="createNewSecret()"
            >
              Create Secret
            </v-btn>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
