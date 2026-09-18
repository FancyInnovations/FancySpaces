<script lang="ts" setup>

  import type { Issue } from '@/api/issues/types'
  import type { Space } from '@/api/spaces/types'
  import Card from '@/components/common/Card.vue'
  import IssueIDChip from '@/components/issues/IssueIDChip.vue'

  const router = useRouter()

  const props = defineProps<{
    space: Space
    issues: Issue[]
  }>()
  const selected = defineModel<string[]>('selected', { default: () => [] })

  const sortedIssues = computed(() => {
    return props.issues.slice().sort((a, b) => b.updated_at.getTime() - a.updated_at.getTime())
  })

  const tableHeaders = [
    { title: 'ID', key: 'id', sortable: false },
    { title: 'Title', key: 'title' },
    { title: 'Type', key: 'type' },
    { title: 'Priority', key: 'priority' },
    { title: 'Status', key: 'status' },
    { title: 'Reporter', key: 'reporter' },
    { title: 'Created', key: 'created_at', value: (issue: Issue) => issue.created_at.toLocaleDateString() },
    { title: 'Updated', key: 'updated_at', value: (issue: Issue) => issue.updated_at.toLocaleDateString() },
  ]

  function onRowClick (event: any, { item }: any) {
    router.push(`/spaces/${props.space.slug}/issues/${item.id}`)
  }

</script>

<template>
  <Card
    width="100%"
  >
    <v-card-text>
      <v-data-table
        v-model="selected"
        class="bg-transparent"
        :headers="tableHeaders"
        hover
        item-key="id"
        item-value="id"
        :items="sortedIssues"
        items-per-page="25"
        show-select
        @click:row="onRowClick"
      >
        <template #item.id="{ item }">
          <IssueIDChip :issue="item" />
        </template>

        <template #item.title="{ item }">
          <div>{{ item.title }}</div>
        </template>

        <template #item.type="{ item }">
          <IssueTypeChip :issue="item" />
        </template>

        <template #item.priority="{ item }">
          <IssuePriorityChip :issue="item" />
        </template>

        <template #item.status="{ item }">
          <IssueStatusChip :issue="item" />
        </template>

        <template #item.reporter="{ item }">
          <UserChip :user="item.reporter" />
        </template>
      </v-data-table>
    </v-card-text>
  </Card>
</template>

<style scoped>

</style>
