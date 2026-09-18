<script lang="ts" setup>

  import type { Issue } from '@/api/issues/types'
  import { getIssue } from '@/api/issues/issues'
  import { useIssueDialogStore } from '@/stores/issue-dialog'

  const issueDialogStore = useIssueDialogStore()
  const loading = ref(false)

  const props = defineProps<{
    issueName?: string
    issue?: Issue
    density?: null | 'default' | 'comfortable' | 'compact'
    withTitle?: boolean
  }>()

  const issueID = computed(() => {
    if (props.issueName) {
      const parts = props.issueName.split('-')
      if (parts.length === 0) return ''
      return parts[parts.length - 1]
    }
    return ''
  })

  const issueSpace = computed(() => {
    if (props.issueName) {
      const parts = props.issueName.split('-')
      if (parts.length <= 1) return ''
      parts.pop()
      return parts.join('-')
    }
    return ''
  })

  async function openDialog () {
    if (props.issue) {
      issueDialogStore.open(props.issue)
    } else if (props.issueName) {
      loading.value = true
      try {
        issueDialogStore.open(await getIssue(issueSpace.value, issueID.value))
      } catch {
        // The chip still remains useful as a non-interactive reference when the issue is unavailable.
      } finally {
        loading.value = false
      }
    }
  }

</script>

<template>
  <v-menu
    close-delay="100"
    location="top"
    offset="8"
    open-delay="100"
    open-on-hover
  >
    <template #activator="{ props: menuProps }">
      <v-chip
        :density="props.density || 'default'"
        color="primary"
        prepend-icon="mdi-sign-text"
        :loading="loading"
        rounded
        v-bind="menuProps"
        variant="tonal"
        @click="openDialog"
      >
        <template v-if="props.withTitle">
          {{ props.issue?.title || issueID }}
        </template>

        <template v-else>
          #{{ props.issue?.id || issueID }}
        </template>
      </v-chip>
    </template>

    <v-card min-width="220">
      <v-card-text v-if="props.issue">
        <p class="text-body-1 mb-2">{{ props.issue?.title }}</p>
        <strong>Type:</strong> {{ props.issue?.type.toUpperCase() }}<br>
        <strong>Priority:</strong> {{ props.issue?.priority.toUpperCase() }}<br>
        <strong>Status:</strong> {{ props.issue?.status.toUpperCase() }}
      </v-card-text>

      <v-card-text v-else>
        <strong>Space:</strong> {{ issueSpace }}<br>
        <strong>Issue ID:</strong> {{ issueID }}
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<style scoped>

</style>
