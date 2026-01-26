<script lang="ts" setup>

import type {Issue, IssueComment} from "@/api/issues/types.ts";
import IssueIDChip from "@/components/issues/IssueIDChip.vue";
import IssueStatusChip from "@/components/issues/IssueStatusChip.vue";
import IssuePriorityChip from "@/components/issues/IssuePriorityChip.vue";
import IssueTypeChip from "@/components/issues/IssueTypeChip.vue";

const props = defineProps<{
  issue: Issue,
  comments: IssueComment[]
}>();

</script>

<template>
  <v-list
    class="sidebar__background ma-4"
    elevation="12"
    height="max-content"
    location="right"
    min-width="250px"
    rounded="xl"
  >
    <v-list-item>
      <v-list-item-title class="text-h6 font-weight-bold">{{ props.issue.title }}</v-list-item-title>
    </v-list-item>

    <v-divider class="mt-2"/>

    <v-list-item>
      <v-list-item-title>
        ID:
        <IssueIDChip
          :issue="props.issue"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item
      v-if="props.issue.external_source"
      :title="'Source: ' + props.issue.external_source"
    />

    <v-list-item>
      <v-list-item-title>
        Type:
        <IssueTypeChip
          :issue="props.issue"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item>
      <v-list-item-title>
        Status:
        <IssueStatusChip
          :issue="props.issue"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item>
      <v-list-item-title>
        Priority:
        <IssuePriorityChip
          :issue="props.issue"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item
      :title="'Reported by: ' + props.issue.reporter"
    />

    <v-list-item
      :title="'Assigned to: ' + (props.issue.assignee || 'Unassigned')"
    />

    <v-list-item
      :title="'Created at: ' + new Date(props.issue.created_at).toLocaleDateString()"
    />

    <v-list-item
      :title="'Last updated: ' + new Date(props.issue.updated_at).toLocaleDateString()"
    />

    <v-list-item
      :title="'Comments: ' + props.comments.length"
    />
  </v-list>
</template>

<style scoped>
.sidebar__background {
  max-height: calc(100vh - 96px);
  background-color: rgba(21, 13, 25, 0.2) !important;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}
</style>
