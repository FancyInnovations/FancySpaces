<script lang="ts" setup>

import IssueTypeChip from "@/components/issues/IssueTypeChip.vue";
import IssueIDChip from "@/components/issues/IssueIDChip.vue";
import IssuePriorityChip from "@/components/issues/IssuePriorityChip.vue";
import type {Issue} from "@/api/issues/types.ts";

const props = defineProps<{
  issue: Issue
}>();

</script>

<template>
  <v-card
    v-for="i in 6"
    :key="i"
    class="mb-4"
    color="#382F2960"
    elevation="6"
    rounded="xl"
  >
    <v-card-text class="d-flex">
      <div class="d-flex flex-column justify-space-between">
        <p class="text-body-1 mb-2">This is a placeholder issue card. Issue functionality is coming soon!</p>
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
