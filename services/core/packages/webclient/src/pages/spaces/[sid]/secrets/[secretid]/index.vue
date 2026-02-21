<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {getSecret, updateSecret} from "@/api/secrets/secrets.ts";
import type {SpaceSecret} from "@/api/secrets/types.ts";
import SpaceHeader from "@/components/SpaceHeader.vue";
import {useNotificationStore} from "@/stores/notifications.ts";
import {useUserStore} from "@/stores/user.ts";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
const notificationStore = useNotificationStore();

const isLoggedIn = ref(false);

const space = ref<Space>();
const secret = ref<SpaceSecret>();

const key = computed(() => {
  if (secret.value) {
    return secret.value.key;
  }
  return "";
})

const newValue = ref('');
const newDescription = ref('');

onMounted(async () => {
  isLoggedIn.value = await userStore.isAuthenticated;

  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!isLoggedIn || !space.value.secrets_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  const secretKey = (route.params as any).secretid as string;
  secret.value = await getSecret(space.value.id, secretKey);

  newDescription.value = secret.value.description;

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

async function updateSecretReq() {
  await updateSecret(
    space.value!.id,
    secret.value!.key,
    newValue.value.length > 0 ? newValue.value : "",
    newDescription.value !== secret.value?.description ? newDescription.value : "",
  );

  newValue.value = '';
  newDescription.value = '';

  notificationStore.info("Secret updated successfully");

  await router.push(`/spaces/${space.value?.slug}/secrets`);
}

const hasChanged = computed(() => {
  return newValue.value.length > 0 || (newDescription.value.length > 0 && newDescription.value !== secret.value?.description);
});

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
        <h1 class="text-center">Update secret</h1>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-text-field
          v-model="key"
          color="primary"
          disabled
          hide-details
          label="Key"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-textarea
          v-model="newValue"
          color="primary"
          hide-details
          label="New secret"
          rows="3"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-textarea
          v-model="newDescription"
          color="primary"
          hide-details
          label="New description"
          rows="6"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-btn
          :disabled="!hasChanged"
          color="primary"
          @click="updateSecretReq()"
        >
          Edit Secret
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
