<script lang="ts" setup>

import IssueTypeChip from "@/components/issues/IssueTypeChip.vue";
import IssueIDChip from "@/components/issues/IssueIDChip.vue";
import IssuePriorityChip from "@/components/issues/IssuePriorityChip.vue";
import type {Issue} from "@/api/issues/types.ts";
import {useIssueDialogStore} from "@/stores/issue-dialog.ts";
import type {Space} from "@/api/spaces/types.ts";

const issueDialogStore = useIssueDialogStore();

const props = defineProps<{
  space: Space,
  issue: Issue
}>();

</script>

<template>
  <v-card
    :to="`/spaces/${space.slug}/issues/${issue.id}`"
    class="mb-4"
    color="#382F2960"
    elevation="6"
    rounded="xl"
  >
    <v-card-text class="d-flex">
      <div class="d-flex flex-column justify-space-between">
        <p
          class="text-body-1"
        >
          {{ issue.title }}
        </p>

        <div class="d-inline-block">
          <IssueIDChip
            :issue="issue"
            class="ma-1"
          />

          <IssueTypeChip
            :issue="issue"
            class="ma-1"
          />

          <IssuePriorityChip
            :issue="issue"
            class="ma-1"
          />
        </div>
      </div>
    </v-card-text>
  </v-card>
</template>

<style scoped>

</style>
