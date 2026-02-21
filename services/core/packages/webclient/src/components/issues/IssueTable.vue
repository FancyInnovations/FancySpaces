<script lang="ts" setup>

import IssueIDChip from "@/components/issues/IssueIDChip.vue";
import type {Issue} from "@/api/issues/types.ts";
import type {Space} from "@/api/spaces/types.ts";
import Card from "@/components/common/Card.vue";

const router = useRouter();

const props = defineProps<{
  space: Space,
  issues: Issue[],
}>();

const sortedIssues = computed(() => {
  return props.issues.slice().sort((a, b) => b.updated_at.getTime() - a.updated_at.getTime());
});

const tableHeaders = [
  { title: 'ID', key: 'id', sortable: false },
  { title: 'Title', key: 'title' },
  { title: 'Type', key: 'type' },
  { title: 'Priority', key: 'priority' },
  { title: 'Status', key: 'status' },
  { title: 'Reporter', key: 'reporter' },
  { title: 'Created', key: 'created_at', value: (issue: Issue) => issue.created_at.toLocaleDateString() },
  { title: 'Updated', key: 'updated_at', value: (issue: Issue) => issue.updated_at.toLocaleDateString() },
];

function onRowClick(event: any, { item }: any) {
  router.push(`/spaces/${props.space.slug}/issues/${item.id}`);
}

</script>

<template>
  <Card
    width="100%"
  >
    <v-card-text>
      <v-data-table
        :headers="tableHeaders"
        :items="sortedIssues"
        class="bg-transparent"
        hover
        item-key="id"
        items-per-page="25"
        @click:row="onRowClick"
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
  </Card>
</template>

<style scoped>

</style>
