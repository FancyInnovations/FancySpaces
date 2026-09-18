<script lang="ts" setup>

  import type { Issue } from '@/api/issues/types'

  const model = defineModel<Partial<Issue>>({ required: true })

  const typeItems = [
    { title: 'Epic', value: 'epic' }, { title: 'Bug', value: 'bug' }, { title: 'Task', value: 'task' },
    { title: 'Story', value: 'story' }, { title: 'Idea', value: 'idea' },
  ]
  const priorityItems = [
    { title: 'Low', value: 'low' }, { title: 'Medium', value: 'medium' }, { title: 'High', value: 'high' }, { title: 'Critical', value: 'critical' },
  ]
  const statusItems = [
    { title: 'Backlog', value: 'backlog' }, { title: 'Planned', value: 'planned' }, { title: 'In Progress', value: 'in_progress' }, { title: 'Done', value: 'done' }, { title: 'Closed', value: 'closed' },
  ]
  const labelsText = computed({
    get: () => (model.value.labels || []).join(', '),
    set: value => {
      model.value.labels = value.split(',').map(label => label.trim()).filter(Boolean)
    },
  })

</script>

<template>
  <v-text-field
    v-model="model.title"
    class="mb-4"
    color="primary"
    hide-details
    label="Title"
    required
  />

  <v-textarea
    v-model="model.description"
    class="mb-4"
    color="primary"
    hide-details
    label="Description (Markdown supported)"
    rows="8"
  />

  <div class="d-flex flex-wrap mb-4 ga-3">
    <v-select
      v-model="model.type"
      color="primary"
      hide-details
      :items="typeItems"
      label="Type"
      min-width="180"
      required
    />

    <v-select
      v-model="model.priority"
      color="primary"
      hide-details
      :items="priorityItems"
      label="Priority"
      min-width="180"
      required
    />

    <v-select
      v-if="model.status"
      v-model="model.status"
      color="primary"
      hide-details
      :items="statusItems"
      label="Status"
      min-width="180"
    />
  </div>

  <v-text-field
    v-model="model.assignee"
    class="mb-4"
    color="primary"
    hide-details
    label="Assignee ID"
  />

  <v-text-field
    v-model="model.parent_issue"
    class="mb-4"
    color="primary"
    hide-details
    label="Parent issue ID"
  />

  <v-text-field
    v-model="labelsText"
    class="mb-4"
    color="primary"
    hint="Comma-separated; labels are shared in this space"
    label="Labels"
    persistent-hint
  />

  <v-text-field
    v-model="model.fix_version"
    class="mb-4"
    color="primary"
    hide-details
    label="Fix version"
  />
</template>
