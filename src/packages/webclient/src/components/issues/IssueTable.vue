<script lang="ts" setup>

import IssueIDChip from "@/components/issues/IssueIDChip.vue";
import type {Issue} from "@/api/issues/types.ts";
import type {Space} from "@/api/spaces/types.ts";

const props = defineProps<{
  space: Space,
  issues: Issue[],
}>();

const tableHeaders = [
  { title: 'ID', key: 'id', sortable: false },
  { title: 'Title', key: 'title' },
  { title: 'Type', key: 'type' },
  { title: 'Priority', key: 'priority' },
  { title: 'Status', key: 'status' },
  { title: 'Reporter', key: 'reporter' },
  { title: 'Created at', key: 'created_at', value: (issue: Issue) => issue.created_at.toLocaleString() },
  { title: 'Updated at', key: 'updated_at', value: (issue: Issue) => issue.updated_at.toLocaleString() },
]

</script>

<template>
  <v-card
    class="card__border"
    color="#19120D33"
    elevation="12"
    rounded="xl"
    width="100%"
  >
    <v-card-text>
      <v-data-table
        :headers="tableHeaders"
        :items="issues"
        class="bg-transparent"
        hover
        item-key="id"
      >
        <template v-slot:item.id="{ item }">
          <IssueIDChip :issue="item" />
        </template>

        <template v-slot:item.title="{ item }">
          <div>{{ item.title }}</div>
        </template>

        <template v-slot:item.type="{ item }">
          <IssueTypeChip :issue="item"/>
        </template>

        <template v-slot:item.priority="{ item }">
          <IssuePriorityChip :issue="item"/>
        </template>

        <template v-slot:item.status="{ item }">
          <IssueStatusChip :issue="item"/>
        </template>

        <template v-slot:item.reporter="{ item }">
          <UserChip :user="item.reporter" />
        </template>
      </v-data-table>
    </v-card-text>
  </v-card>
</template>

<style scoped>

</style>
