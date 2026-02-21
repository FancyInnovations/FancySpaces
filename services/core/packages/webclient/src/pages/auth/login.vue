<script lang="ts" setup>

import {useNotificationStore} from "@/stores/notifications.ts";

const notificationStore = useNotificationStore();

const apiKey = ref('');

function login() {
  localStorage.setItem('fs_api_key', apiKey.value);
  notificationStore.info('Logged in successfully!');
}

function logout() {
  localStorage.removeItem('fs_api_key');
  notificationStore.info('Logged out successfully!');
}
</script>

<template>
  <v-container>
    <v-row>
      <v-col>
        <h1 class="text-center">Login</h1>
        <p class="text-center">Paste your API key to continue</p>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="4">
        <v-text-field
          v-model="apiKey"
          color="primary"
          hide-details
          label="API key"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <div class="d-flex justify-center">
          <v-btn
            color="primary"
            large
            variant="tonal"
            @click="login"
          >
            Login
          </v-btn>

          <v-btn
            class="ml-4"
            color="primary"
            large
            variant="tonal"
            @click="logout"
          >
            Logout
          </v-btn>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
