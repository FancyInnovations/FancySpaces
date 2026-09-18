<script lang="ts" setup>

  import type { Issue, IssueComment } from '@/api/issues/types'
  import IssueExternalSourceChip from '@/components/issues/IssueExternalSourceChip.vue'
  import IssueIDChip from '@/components/issues/IssueIDChip.vue'
  import IssuePriorityChip from '@/components/issues/IssuePriorityChip.vue'
  import IssueStatusChip from '@/components/issues/IssueStatusChip.vue'
  import IssueTypeChip from '@/components/issues/IssueTypeChip.vue'
  import VersionChip from '@/components/VersionChip.vue'

  const props = defineProps<{
    issue?: Issue
    comments: IssueComment[]
  }>()

  const formattedCreatedAt = ref('')
  const formattedUpdatedAt = ref('')
  const formattedResolvedAt = ref('')

  function formatDate (date: Date): string {
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const diffInHours = diff / (1000 * 60 * 60)
    if (diffInHours < 12) {
      if (diffInHours < 1) {
        const diffInMinutes = diff / (1000 * 60)
        return diffInMinutes < 1 ? `${Math.floor(diff / 1000)}s ago` : `${Math.floor(diffInMinutes)}min ago`
      } else {
        return `${Math.floor(diffInHours)}h ago`
      }
    } else {
      const createdAtDate = date.getDate()
      const nowDate = now.getDate()
      const createdAtMonth = date.getMonth()
      const nowMonth = now.getMonth()
      const createdAtYear = date.getFullYear()
      const nowYear = now.getFullYear()

      if (createdAtYear === nowYear && createdAtMonth === nowMonth && createdAtDate === nowDate - 1) {
        return 'Yesterday'
      } else if (createdAtYear === nowYear && createdAtMonth === nowMonth && createdAtDate === nowDate) {
        return 'Today'
      } else {
        return date.toLocaleDateString()
      }
    }
  }

  onMounted(() => {
    formattedCreatedAt.value = formatDate(new Date(props.issue?.created_at || ''))
    formattedUpdatedAt.value = formatDate(new Date(props.issue?.updated_at || ''))
    formattedResolvedAt.value = props.issue?.resolved_at
      ? formatDate(new Date(props.issue?.resolved_at))
      : 'Unresolved'

    setInterval(() => {
      formattedCreatedAt.value = formatDate(new Date(props.issue?.created_at || ''))
      formattedUpdatedAt.value = formatDate(new Date(props.issue?.updated_at || ''))
      formattedResolvedAt.value = props.issue?.resolved_at
        ? formatDate(new Date(props.issue?.resolved_at))
        : 'Unresolved'
    }, 1000)
  })

</script>

<template>
  <v-list
    v-if="props.issue"
    class="sidebar__background"
    elevation="12"
    location="right"
    min-width="250px"
    rounded="xl"
  >
    <v-list-item>
      <v-list-item-title>
        ID:
        <IssueIDChip
          class="ml-2"
          density="compact"
          :issue="props.issue"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item v-if="props.issue?.parent_issue">
      <v-list-item-title>
        Parent:
        <IssueIDChip
          class="ml-2"
          density="compact"
          :issue-name="props.issue?.parent_issue"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item v-if="props.issue?.external_source">
      <v-list-item-title>
        Source:
        <IssueExternalSourceChip
          class="ml-2"
          density="compact"
          :issue="props.issue!"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item>
      <v-list-item-title>
        Type:
        <IssueTypeChip
          class="ml-2"
          density="compact"
          :issue="props.issue!"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item>
      <v-list-item-title>
        Status:
        <IssueStatusChip
          class="ml-2"
          density="compact"
          :issue="props.issue!"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item>
      <v-list-item-title>
        Priority:
        <IssuePriorityChip
          class="ml-2"
          density="compact"
          :issue="props.issue!"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item v-if="props.issue?.affected_versions && props.issue?.affected_versions.length > 0">
      <v-list-item-title>
        Affected versions:
      </v-list-item-title>
    </v-list-item>

    <div
      v-if="props.issue?.affected_versions && props.issue?.affected_versions.length > 0"
      class="ml-4"
    >
      <VersionChip
        v-for="version in props.issue.affected_versions"
        :key="version"
        density="compact"
        :space-i-d="issue!.space"
        :version="version"
      />
    </div>

    <v-list-item v-if="props.issue?.fix_version">
      <v-list-item-title>
        Fix version:
        <VersionChip
          class="ml-2"
          density="compact"
          :space-i-d="issue!.space"
          :version="props.issue.fix_version"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item
      v-if="props.issue?.resolved_at"
      :title="'Resolved at: ' + formattedResolvedAt"
    />

    <v-list-item>
      <v-list-item-title>
        Reporter:
        <UserChip
          class="ml-2"
          density="compact"
          :user="props.issue!.reporter"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item v-if="props.issue?.assignee">
      <v-list-item-title>
        Assignee:
        <UserChip
          class="ml-2"
          density="compact"
          :user="props.issue!.assignee"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item
      :title="'Created at: ' + formattedCreatedAt"
    />

    <v-list-item
      :title="'Last updated: ' + formattedUpdatedAt"
    />

    <v-list-item
      :title="'Comments: ' + props.comments?.length"
    />
  </v-list>
</template>

<style scoped>
.sidebar__background {
  background-color: transparent !important;
  border: 1px solid rgba(255, 255, 255, 0.1);
  scrollbar-width: none;
}
</style>
