<script lang="ts" setup>

import IssueIDChip from "@/components/issues/IssueIDChip.vue";
import type {Issue} from "@/api/issues/types.ts";
import type {Space} from "@/api/spaces/types.ts";

const issues = ref<Issue[]>([]);

onMounted(() => {
  for (let i = 0; i < 20; i++) {
    issues.value.push({
      id: '7G5B1',
      space: 'fc',
      title: 'NPE in some module causing crashes',
      description: '## 🐞 Bug Description\n'+
        'A defect was identified that causes unexpected behavior in the application. Further investigation is required to determine the root cause and scope of impact.\n'+
        '\n'+
        '## 🔁 Steps to Reproduce\n'+
        '1. Navigate to `[page / feature]`\n'+
        '2. Perform `[action]`\n'+
        '3. Observe the result\n'+
        '\n'+
        '## ✅ Expected Result\n'+
        'The system should `[expected behavior]`.\n'+
        '\n'+
        '## ❌ Actual Result\n'+
        'The system instead `[actual behavior]`.\n'+
        '\n'+
        '## 🌍 Environment\n'+
        '- App version: `[version]`\n'+
        '- Environment: `[dev / staging / prod]`\n'+
        '- Browser / Device: `[if applicable]`\n'+
        '\n'+
        '## 📎 Additional Notes\n'+
        '- Frequency: `[always / intermittent / once]`\n'+
        '- Severity: `[low / medium / high / critical]`\n'+
        '- Attachments: `[logs / screenshots / videos if any]`',
      type: 'bug',
      priority: 'low',
      status: 'done',
      assignee: 'user123',
      reporter: 'user456',
      created_at: new Date(2026, 0, 26, 21, 0, 0, 0),
      updated_at: new Date(),
      external_source: 'github',
      affected_versions: ['1.2.2', '1.2.3', '1.2.4'],
      fix_version: '2.0.0',
      resolved_at: new Date(2026, 0, 27, 15, 30, 0, 0),
      parent_issue: 'fc-ABC123',
      extra_fields: {
        github_url: 'https://github.com/FancyInnovations/FancyPlugins/issues/195'
      }
    });
  }
})

const props = defineProps<{
  space: Space,
  typeFilter: string|undefined,
  priorityFilter: string|undefined,
  statusFilter: string|undefined
}>();

const tableHeaders = [
  { title: 'ID', key: 'id', sortable: false },
  { title: 'Title', key: 'title', sortable: false },
  { title: 'Type', key: 'type', sortable: false },
  { title: 'Priority', key: 'priority', sortable: false },
  { title: 'Status', key: 'status', sortable: false },
  { title: 'Reporter', key: 'reporter', sortable: false },
  { title: 'Created at', key: 'created_at', sortable: false, value: (issue: Issue) => issue.created_at.toLocaleString() },
  { title: 'Updated at', key: 'updated_at', sortable: false, value: (issue: Issue) => issue.updated_at.toLocaleString() },
]

</script>

<template>
  <v-data-table
    :headers="tableHeaders"
    :items="issues"
    class="elevation-1"
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
</template>

<style scoped>

</style>
