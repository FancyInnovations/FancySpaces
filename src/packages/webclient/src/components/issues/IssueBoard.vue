<script lang="ts" setup>

import type {Issue} from "@/api/issues/types.ts";
import type {Space} from "@/api/spaces/types.ts";
import IssueCard from "@/components/issues/IssueCard.vue";

const plannedIssues = computed(() => {
  return props.issues.filter(issue => issue.status === 'planned');
});

const inProgressIssues = computed(() => {
  return props.issues.filter(issue => issue.status === 'in_progress');
});

const doneIssues = computed(() => {
  return props.issues.filter(issue => issue.status === 'done');
});

const props = defineProps<{
  space: Space,
  issues: Issue[],
}>();

</script>

<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="4">
        <v-card
          class="card__border"
          color="#19120D33"
          elevation="12"
          min-height="600"
          rounded="xl"
        >
          <v-card-title class="my-2 ml-2">Planned</v-card-title>

          <v-card-text>
            <IssueCard
              v-for="issue in plannedIssues"
              v-if="plannedIssues.length > 0"
              :key="issue.id"
              :issue="issue"
            />
            <p v-else class="text-center my-4">No planned issues.</p>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="4">
        <v-card
          class="card__border"
          color="#19120D33"
          elevation="12"
          min-height="600"
          rounded="xl"
        >
          <v-card-title class="my-2 ml-2">In Progress</v-card-title>

          <v-card-text>
            <IssueCard
              v-for="issue in inProgressIssues"
              v-if="inProgressIssues.length > 0"
              :key="issue.id"
              :issue="issue"
            />
            <p v-else class="text-center my-4">No issues in progress.</p>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="4">
        <v-card
          class="card__border"
          color="#19120D33"
          elevation="12"
          min-height="600"
          rounded="xl"
        >
          <v-card-title class="my-2 ml-2">Done</v-card-title>

          <v-card-text>
            <IssueCard
              v-for="issue in doneIssues"
              v-if="doneIssues.length > 0"
              :key="issue.id"
              :issue="issue"
            />
            <p v-else class="text-center my-4">No done issues.</p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
