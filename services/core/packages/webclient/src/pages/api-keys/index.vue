<script lang="ts" setup>

import {onMounted} from "vue";
import {useHead} from "@vueuse/head";
import router from "@/router";
import {useUserStore} from "@/stores/user.ts";
import {deleteApiKey, getApiKeys} from "@/api/auth/api-keys.ts";
import type {ApiKey} from "@/api/auth/types.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import {useConfirmationStore} from "@/stores/confirmation.ts";

const userStore = useUserStore();
const notificationStore = useNotificationStore();
const confirmationStore = useConfirmationStore();

const apiKeys = ref<ApiKey[]>([]);

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

  apiKeys.value = await getApiKeys(userStore.user!.id);
});

function revokeApiKey(apiKeyId: string) {
  confirmationStore.confirmation = {
    shown: true,
    persistent: true,
    title: "Revoke API Key",
    text: "Are you sure you want to revoke this API key? This action cannot be undone.",
    yesText: "Revoke",
    onConfirm: async () => {
      await deleteApiKey(apiKeyId)

      apiKeys.value = apiKeys.value.filter(key => key.key_id != apiKeyId);
      notificationStore.info("API key revoked.");
    }
  };
}

</script>

<template>
  <v-container width="60%">
    <v-row>
      <v-col cols="12">
        <h1>API Keys</h1>
        <p>Manage your API keys for accessing the FancySpaces API. Create, view, and revoke API keys to control access to your account.</p>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <Card>
          <v-card-text>
            <v-data-table :headers="[
              { title: 'ID', value: 'key_id' },
              { title: 'Description', value: 'description' },
              { title: 'Created At', value: 'created_at' },
              { title: 'Last Used', value: 'last_used_at' },
              { title: 'Actions', value: 'actions', sortable: false }
            ]"
            :items="apiKeys"
            :items-per-page="5"
            class="bg-transparent"
            >
              <template #item.created_at="{ item }">
                {{ item.created_at.toLocaleString() }}
              </template>

              <template #item.last_used_at="{ item }">
                {{ item.last_used_at ? item.last_used_at.toLocaleString() : "Never" }}
              </template>

              <template #item.actions="{ item }">
                <v-btn
                  color="red"
                  size="small"
                  @click="revokeApiKey(item.key_id)"
                >
                  Revoke
                </v-btn>
              </template>
            </v-data-table>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <v-btn
          color="primary"
          to="/api-keys/new"
        >
          Create New API Key
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
