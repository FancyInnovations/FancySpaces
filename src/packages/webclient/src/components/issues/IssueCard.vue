<script lang="ts" setup>

import IssueTypeChip from "@/components/issues/IssueTypeChip.vue";
import IssueIDChip from "@/components/issues/IssueIDChip.vue";
import IssuePriorityChip from "@/components/issues/IssuePriorityChip.vue";
import type {Issue} from "@/api/issues/types.ts";
import {useIssueDialogStore} from "@/stores/issue-dialog.ts";

const issueDialogStore = useIssueDialogStore();

const props = defineProps<{
  issue: Issue
}>();

function openIssueDialog() {
  issueDialogStore.open(props.issue);
}

</script>

<template>
  <v-card
    class="mb-4"
    color="#382F2960"
    elevation="6"
    rounded="xl"
    @click="openIssueDialog"
  >
    <v-card-text class="d-flex">
      <div class="d-flex flex-column justify-space-between">
        <p
          class="text-body-1 mb-2"
        >
          {{ issue.title }}
        </p>

        <div class="d-inline-block">
          <IssueIDChip
            :issue="issue"
          />

          <IssueTypeChip
            :issue="issue"
            class="ml-2 mr-2"
          />

          <IssuePriorityChip
            :issue="issue"
            class="mr-2"
          />
        </div>
      </div>
    </v-card-text>
  </v-card>
</template>

<style scoped>

</style>
