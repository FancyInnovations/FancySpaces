<script lang="ts" setup>

  import type { Issue } from '@/api/issues/types'
  import type { Space } from '@/api/spaces/types'
  import IssueIDChip from '@/components/issues/IssueIDChip.vue'
  import IssuePriorityChip from '@/components/issues/IssuePriorityChip.vue'
  import IssueTypeChip from '@/components/issues/IssueTypeChip.vue'
  const props = defineProps<{
    space: Space
    issue: Issue
  }>()

</script>

<template>
  <v-card
    class="mb-4"
    color="#382F2960"
    elevation="6"
    rounded="xl"
    :to="`/spaces/${space.slug}/issues/${issue.id}`"
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
            class="ma-1"
            :issue="issue"
          />

          <IssueTypeChip
            class="ma-1"
            :issue="issue"
          />

          <IssuePriorityChip
            class="ma-1"
            :issue="issue"
          />
        </div>
      </div>
    </v-card-text>
  </v-card>
</template>

<style scoped>

</style>
