<script lang="ts" setup>

  import type { Issue } from '@/api/issues/types'
  import type { Space } from '@/api/spaces/types'
  import Card from '@/components/common/Card.vue'
  import IssueCard from '@/components/issues/IssueCard.vue'

  const props = defineProps<{
    space: Space
    issues: Issue[]
  }>()

  const emit = defineEmits<{
    statusChange: [issue: Issue, status: Issue['status']]
  }>()

  const columns: Array<{ status: Issue['status'], title: string, empty: string }> = [
    { status: 'backlog', title: 'Backlog', empty: 'No backlog issues.' },
    { status: 'planned', title: 'Planned', empty: 'No planned issues.' },
    { status: 'in_progress', title: 'In Progress', empty: 'No issues in progress.' },
    { status: 'done', title: 'Done', empty: 'No done issues.' },
    { status: 'closed', title: 'Closed', empty: 'No closed issues.' },
  ]

  const draggedIssue = ref<Issue | null>(null)
  const dropStatus = ref<Issue['status'] | null>(null)

  function issuesFor (status: Issue['status']) {
    return props.issues.filter(issue => issue.status === status)
  }

  function startDrag (event: DragEvent, issue: Issue) {
    draggedIssue.value = issue
    if (event.dataTransfer) {
      event.dataTransfer.effectAllowed = 'move'
      event.dataTransfer.setData('text/plain', issue.id)
    }
  }

  function dragOver (event: DragEvent, status: Issue['status']) {
    event.preventDefault()
    if (event.dataTransfer) event.dataTransfer.dropEffect = 'move'
    dropStatus.value = status
  }

  function drop (status: Issue['status']) {
    if (draggedIssue.value && draggedIssue.value.status !== status) {
      emit('statusChange', draggedIssue.value, status)
    }
    draggedIssue.value = null
    dropStatus.value = null
  }

  function clearDrag () {
    draggedIssue.value = null
    dropStatus.value = null
  }

</script>

<template>
  <v-container class="issue-board">
    <v-row>
      <v-col
        v-for="column in columns"
        :key="column.status"
        class="pt-0"
        cols="12"
        lg="auto"
        md="6"
        @dragleave="dropStatus = null"
        @dragover="dragOver($event, column.status)"
        @drop="drop(column.status)"
      >
        <Card
          class="issue-column"
          :class="{ 'issue-column--drop-target': dropStatus === column.status }"
          min-height="600"
        >
          <v-card-title class="my-2 ml-2">
            {{ column.title }} ({{ issuesFor(column.status).length }})
          </v-card-title>

          <v-card-text>
            <div
              v-for="issue in issuesFor(column.status)"
              :key="issue.id"
              class="issue-drag-item"
              draggable="true"
              @dragend="clearDrag"
              @dragstart="startDrag($event, issue)"
            >
              <IssueCard
                :issue="issue"
                :space="space"
              />
            </div>

            <p v-if="issuesFor(column.status).length === 0" class="text-center my-4">
              {{ column.empty }}
            </p>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.issue-board :deep(.v-row) {
  flex-wrap: nowrap;
  overflow-x: auto;
}

.issue-board .v-col {
  flex: 1 0 280px;
  max-width: 360px;
}

.issue-column {
  transition: border-color 120ms ease, background-color 120ms ease;
}

.issue-column--drop-target {
  border: 1px dashed rgb(var(--v-theme-primary));
  background: rgba(var(--v-theme-primary), 0.08);
}

.issue-drag-item {
  cursor: grab;
}

.issue-drag-item:active {
  cursor: grabbing;
}
</style>
