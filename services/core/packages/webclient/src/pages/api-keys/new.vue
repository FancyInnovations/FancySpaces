<script lang="ts" setup>

import {onMounted} from "vue";
import {useHead} from "@vueuse/head";
import router from "@/router";
import {useUserStore} from "@/stores/user.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import {createApiKey} from "@/api/auth/api-keys.ts";

const userStore = useUserStore();
const notificationStore = useNotificationStore();

const description = ref('');
const key = ref('');

onMounted(async () => {
  useHead({
    title: `FancySpaces - API Keys`,
    meta: [
      {
        name: 'description',
        content: 'Manage your API keys for accessing the FancySpaces API. Create, view, and revoke API keys to control access to your account.'
      }
    ]
  });

  if (!(await userStore.isAuthenticated)) {
    await router.push("/");
  }
});

async function createApiKeyReq() {
  key.value = await createApiKey(description.value);
  notificationStore.info("API key created. Make sure to copy the API key now, as it will not be shown again for security reasons.");
}

</script>

<template>
  <v-container width="60%">
    <v-row>
      <v-col cols="12">
        <h1>API Keys</h1>
        <p>Create a new API key to access the FancySpaces API. Provide a description for your API key to easily identify its purpose.</p>
      </v-col>
    </v-row>

    <v-row v-if="key.length == 0">
      <v-col cols="12">
        <Card>
          <v-card-text>
            <v-text-field
              v-model="description"
              color="primary"
              hide-details
              label="Description (optional)"
              placeholder="E.g. 'My API key for accessing the API from my script'"
              required
            />

            <v-btn
              class="mt-4"
              color="primary"
              @click="createApiKeyReq"
            >
              Create API Key
            </v-btn>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
    <v-row v-if="key.length > 0">
      <v-col cols="12">
        <Card>
          <v-card-text>
            <p>
              <strong>API Key:</strong>
              <span class="api-key ml-2">{{ key }}</span>
            </p>

            <p class="mt-4 text-red">Make sure to copy the API key now, as it will not be shown again for security reasons.</p>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-btn
          class="mt-4"
          to="/api-keys"
          variant="text"
        >
          Back to API Keys
        </v-btn>
      </v-col>
    </v-row>

  </v-container>
</template>

<style scoped>
.api-key {
  font-family: "Courier New", Courier, monospace;
  background-color: rgba(236, 76, 61, 0.2) !important;
  padding: 2px 4px;
  border-radius: 4px;
  user-select: all;
}
</style>
