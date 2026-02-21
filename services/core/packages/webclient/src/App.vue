<script lang="ts" setup>
import AppHeader from "@/components/AppHeader.vue";
import {useConfirmationStore} from "@/stores/confirmation.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import ConfirmationDialog from "@/components/common/ConfirmationDialog.vue";

const confirmationStore = useConfirmationStore();

const notifications = useNotificationStore();

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
