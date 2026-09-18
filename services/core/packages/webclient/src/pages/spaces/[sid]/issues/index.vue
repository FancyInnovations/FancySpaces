<script lang="ts" setup>

  import type { Issue } from '@/api/issues/types'
  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { bulkUpdateIssues, getIssues, updateIssue } from '@/api/issues/issues'
  import { getSpace } from '@/api/spaces/spaces'
  import Card from '@/components/common/Card.vue'
  import SpaceHeader from '@/components/SpaceHeader.vue'
  import SpaceSidebar from '@/components/SpaceSidebar.vue'
  import { useNotificationStore } from '@/stores/notifications'
  import { useUserStore } from '@/stores/user'

  const router = useRouter()
  const route = useRoute()
  const userStore = useUserStore()
  const notificationStore = useNotificationStore()

  const space = ref<Space>()
  const issues = ref<Issue[]>([])
  const loading = ref(true)
  const errorMessage = ref('')
  const total = ref(0)
  const offset = ref(0)
  const selected = ref<string[]>([])
  const bulkStatus = ref<Issue['status']>()
  const savedViews = ref<Array<{ name: string, filters: Record<string, string> }>>([])
  const savedViewName = ref('')
  const searchField = ref<any>()

  const openIssues = computed(() => {
    return issues.value.filter(issue => issue.status !== 'closed')
  })

  const closedIssues = computed(() => {
    return issues.value.filter(issue => issue.status === 'done' || issue.status === 'closed')
  })

  const filteredIssues = computed(() => issues.value)

  const displayType = ref<'board' | 'list'>('board')
  const searchQuery = ref('')
  const typeFilter = ref()
  const priorityFilter = ref()
  const statusFilter = ref()
  const labelFilter = ref('')

  const canWrite = computed(() => {
    const userID = userStore.user?.id
    if (!userID || !space.value) return false
    return space.value.creator === userID || space.value.members.some(member => member.user_id === userID && ['member', 'admin'].includes(member.role))
  })

  async function load () {
    loading.value = true
    errorMessage.value = ''
    try {
      await userStore.isAuthenticated
      const spaceID = (route.params as any).sid as string
      space.value = await getSpace(spaceID)

      if (!space.value.issue_settings.enabled) {
        await router.push(`/spaces/${space.value.slug}`)
        return
      }

      searchQuery.value = String(route.query.q || '')
      typeFilter.value = route.query.type || undefined
      priorityFilter.value = route.query.priority || undefined
      statusFilter.value = route.query.status || undefined
      labelFilter.value = String(route.query.label || '')
      const result = await getIssues(space.value.id, { limit: 50, offset: offset.value, q: searchQuery.value, type: typeFilter.value, priority: priorityFilter.value, status: statusFilter.value, label: labelFilter.value })
      issues.value = result.items
      total.value = result.total
      const savedDisplayType = localStorage.getItem(`issues_display_type_${space.value.id}`)
      if (savedDisplayType === 'board' || savedDisplayType === 'list') displayType.value = savedDisplayType
      savedViews.value = JSON.parse(localStorage.getItem(`issues_saved_views_${space.value.id}`) || '[]')

      useHead({
        title: `${space.value.title} issues - FancySpaces`,
        meta: [{ name: 'description', content: space.value.summary || 'View issues for this space on FancySpaces.' }],
      })
    } catch (error) {
      errorMessage.value = error instanceof Error ? error.message : 'Failed to load issues.'
    } finally {
      loading.value = false
    }
  }

  onMounted(load)
  function keyboardShortcut (event: KeyboardEvent) {
    const target = event.target as HTMLElement | null
    if (target?.matches('input, textarea, [contenteditable="true"]')) return
    if (event.key === '/' && searchField.value) {
      event.preventDefault(); searchField.value.focus?.()
    }
    if (event.key.toLowerCase() === 'c' && canWrite.value && space.value) router.push(`/spaces/${space.value.slug}/issues/new`)
  }
  onMounted(() => window.addEventListener('keydown', keyboardShortcut))
  onBeforeUnmount(() => window.removeEventListener('keydown', keyboardShortcut))

  // Watch for changes in displayType and save to localStorage
  watch(displayType, newType => {
    if (space.value) {
      localStorage.setItem(`issues_display_type_${space.value.id}`, newType)
    }
  })

  let filterTimer: ReturnType<typeof setTimeout> | undefined
  watch([searchQuery, typeFilter, priorityFilter, statusFilter, labelFilter], () => {
    if (!space.value) return
    clearTimeout(filterTimer)
    filterTimer = setTimeout(async () => {
      offset.value = 0
      await router.replace({ query: { q: searchQuery.value || undefined, type: typeFilter.value || undefined, priority: priorityFilter.value || undefined, status: statusFilter.value || undefined, label: labelFilter.value || undefined } })
      await loadPage()
    }, 250)
  })

  async function loadPage () {
    if (!space.value) return
    try {
      const result = await getIssues(space.value.id, { limit: 50, offset: offset.value, q: searchQuery.value, type: typeFilter.value, priority: priorityFilter.value, status: statusFilter.value, label: labelFilter.value })
      issues.value = result.items
      total.value = result.total
      selected.value = []
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to load issues.')
    }
  }

  function saveView () {
    if (!space.value || !savedViewName.value.trim()) return
    savedViews.value = [...savedViews.value.filter(view => view.name !== savedViewName.value.trim()), { name: savedViewName.value.trim(), filters: { q: searchQuery.value, type: typeFilter.value || '', priority: priorityFilter.value || '', status: statusFilter.value || '', label: labelFilter.value } }]
    localStorage.setItem(`issues_saved_views_${space.value.id}`, JSON.stringify(savedViews.value))
    savedViewName.value = ''
  }

  function applyView (view: { filters: Record<string, string> }) {
    searchQuery.value = view.filters.q || ''; typeFilter.value = view.filters.type || undefined; priorityFilter.value = view.filters.priority || undefined; statusFilter.value = view.filters.status || undefined; labelFilter.value = view.filters.label || ''
  }

  async function applyBulkStatus () {
    if (!space.value || !bulkStatus.value || selected.value.length === 0) return
    try {
      const updated = await bulkUpdateIssues(space.value.id, selected.value, { status: bulkStatus.value })
      const byID = new Map(updated.map(issue => [issue.id, issue]))
      issues.value = issues.value.map(issue => byID.get(issue.id) || issue)
      selected.value = []
      notificationStore.info('Issues updated successfully')
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to update issues.')
    }
  }

  async function statusChanged (issue: Issue, newStatus: Issue['status']) {
    const previousStatus = issue.status
    issue.status = newStatus
    try {
      const saved = await updateIssue(issue.space, issue.id, { status: newStatus })
      Object.assign(issue, saved)
      notificationStore.info('Issue status updated successfully')
    } catch (error) {
      issue.status = previousStatus
      notificationStore.error(error instanceof Error ? error.message : 'Failed to update issue status.')
    }
  }

</script>

<template>
  <v-container width="90%">
    <v-row>
      <v-col class="flex-grow-0 pa-0">
        <SpaceSidebar
          :space="space"
        />
      </v-col>

      <v-col>
        <SpaceHeader :space="space">
          <template #metadata>
            <p class="text-body-2 mx-4">-</p>
            <p class="text-body-2">{{ openIssues.length }} open issues</p>
            <p class="text-body-2 mx-4">-</p>
            <p class="text-body-2">{{ closedIssues.length }} resolved issues</p>
          </template>

          <template #quick-actions>
            <v-btn
              v-if="canWrite"
              color="primary"
              size="large"
              :to="`/spaces/${space?.slug}/issues/new`"
              variant="tonal"
            >
              New Issue
            </v-btn>
          </template>
        </SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        >
      </v-col>
    </v-row>

    <v-row v-if="loading" justify="center">
      <v-col cols="12" md="8">
        <Card><v-card-text class="text-center">Loading issues…</v-card-text></Card>
      </v-col>
    </v-row>

    <v-row v-else-if="errorMessage" justify="center">
      <v-col cols="12" md="8">
        <Card>
          <v-card-text class="text-center">
            <p class="mb-4">{{ errorMessage }}</p>
            <v-btn color="primary" @click="load">Retry</v-btn>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>

    <v-row v-else>
      <v-col
        class="d-flex align-center justify-space-between flex-wrap"
      >
        <Card>
          <v-card-text>
            <div class="d-flex align-center justify-space-between">
              <div class="d-flex align-center flex-wrap">
                <v-text-field
                  ref="searchField"
                  v-model="searchQuery"
                  class="ma-2"
                  clearable
                  color="primary"
                  density="compact"
                  hide-details
                  label="Search issues"
                  min-width="300"
                  prepend-inner-icon="mdi-magnify"
                />

                <v-select
                  v-model="typeFilter"
                  class="ma-2"
                  clearable
                  color="primary"
                  density="compact"
                  hide-details
                  :items="[
                    { title: 'Epic', value: 'epic' },
                    { title: 'Bug', value: 'bug' },
                    { title: 'Task', value: 'task' },
                    { title: 'Story', value: 'story' },
                    { title: 'Idea', value: 'idea' },

                  ]"
                  label="Type"
                  min-width="200"
                />

                <v-text-field
                  v-model="labelFilter"
                  class="ma-2"
                  clearable
                  color="primary"
                  density="compact"
                  hide-details
                  label="Label"
                  min-width="160"
                  prepend-inner-icon="mdi-tag"
                />

                <v-select
                  v-model="priorityFilter"
                  class="ma-2"
                  clearable
                  color="primary"
                  density="compact"
                  hide-details
                  :items="[
                    { title: 'Low', value: 'low' },
                    { title: 'Medium', value: 'medium' },
                    { title: 'High', value: 'high' },
                    { title: 'Critical', value: 'critical' },
                  ]"
                  label="Priority"
                  min-width="200"
                />

                <v-select
                  v-model="statusFilter"
                  class="ma-2"
                  clearable
                  color="primary"
                  density="compact"
                  hide-details
                  :items="[
                    { title: 'Backlog', value: 'backlog' },
                    { title: 'Planned', value: 'planned' },
                    { title: 'In Progress', value: 'in_progress' },
                    { title: 'Done', value: 'done' },
                    { title: 'Closed', value: 'closed' },

                  ]"
                  label="Status"
                  min-width="200"
                />
              </div>

              <div class="d-flex align-center flex-wrap ga-2 mt-2">
                <v-select
                  v-if="savedViews.length > 0"
                  class="ma-2"
                  density="compact"
                  hide-details
                  :items="savedViews.map(view => ({ title: view.name, value: view.name }))"
                  label="Saved views"
                  @update:model-value="name => applyView(savedViews.find(view => view.name === name)!)"
                />

                <v-text-field v-model="savedViewName" density="compact" hide-details label="Save current view" />
                <v-btn size="small" variant="tonal" @click="saveView">Save view</v-btn>
              </div>
            </div>
          </v-card-text>
        </Card>

        <Card
          class="margin-top flex-grow-1"
          max-width="fit-content"
        >
          <v-card-text>
            <v-btn-group
              density="compact"
            >
              <v-btn
                color="primary"
                prepend-icon="mdi-view-dashboard"
                :variant="displayType === 'board' ? 'tonal' : 'outlined'"
                @click="displayType = 'board'"
              >
                Board
              </v-btn>

              <v-btn
                color="primary"
                prepend-icon="mdi-format-list-bulleted"
                :variant="displayType === 'list' ? 'tonal' : 'outlined'"
                @click="displayType = 'list'"
              >
                List
              </v-btn>
            </v-btn-group>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>

    <v-row v-if="!loading && !errorMessage && space">
      <v-col>
        <Card v-if="selected.length > 0 && canWrite" class="mb-4"><v-card-text class="d-flex align-center ga-3"><span>{{ selected.length }} selected</span>

          <v-select
            v-model="bulkStatus"
            density="compact"
            hide-details
            :items="['backlog', 'planned', 'in_progress', 'done', 'closed']"
            label="Set status"
          />

          <v-btn color="primary" @click="applyBulkStatus">Apply</v-btn></v-card-text></Card>

        <IssueBoard
          v-if="displayType === 'board'"
          :issues="filteredIssues"
          :space="space!"
          @status-change="statusChanged"
        />

        <IssueTable
          v-else
          v-model:selected="selected"
          :issues="filteredIssues"
          :space="space!"
        />

        <div class="d-flex justify-space-between mt-4"><v-btn :disabled="offset === 0" @click="offset = Math.max(0, offset - 50); loadPage()">Previous</v-btn><span>{{ Math.min(offset + 1, total) }}–{{ Math.min(offset + issues.length, total) }} of {{ total }}</span><v-btn :disabled="offset + issues.length >= total" @click="offset += 50; loadPage()">Next</v-btn></div>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
@media (max-width: 1919px) {
  .margin-top {
    margin-top: 16px !important;
  }
}
</style>
