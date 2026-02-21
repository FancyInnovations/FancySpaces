<script lang="ts" setup>
import AppHeader from "@/components/AppHeader.vue";
import {useConfirmationStore} from "@/stores/confirmation.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import ConfirmationDialog from "@/components/common/ConfirmationDialog.vue";
import {useUserStore} from "@/stores/user.ts";
import {refreshToken} from "@/api/auth/tokens.ts";

const confirmationStore = useConfirmationStore();
const notifications = useNotificationStore();
const userStore = useUserStore();

userStore.loadTokenFromStorage();
userStore.loadUserFromStorage();

// Refresh token if it's about to expire
if (userStore.token) {
  if (userStore.tokenTTL < 1000 * 60 * 60 * 24) {
    refreshToken(userStore.token).then(value => {
      userStore.setToken(value);
    }).catch(() => {
      userStore.clearUser();
      userStore.clearToken();
      // TODO send error message
    })
  }
}

function confirm() {
  confirmationStore.confirmation.onConfirm();
  confirmationStore.confirmation.shown = false;
}
</script>


<template>
  <v-app>
    <v-snackbar-queue
      v-model="notifications.queue"
      rounded="lg"
      timeout="1500"
    />

    <ConfirmationDialog
      :persistent="confirmationStore.confirmation.persistent"
      :shown="confirmationStore.confirmation.shown"
      :text="confirmationStore.confirmation.text"
      :title="confirmationStore.confirmation.title"
      :yesText="confirmationStore.confirmation.yesText"
      @clickedClose="confirmationStore.confirmation.shown = false"
      @clickedYes="confirm()"
    />

    <AppHeader/>

    <router-view />

    <AppFooter/>

    <IssueDialog/>
  </v-app>
</template>
