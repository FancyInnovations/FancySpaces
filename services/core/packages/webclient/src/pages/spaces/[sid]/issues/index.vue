<script lang="ts" setup>

  import type { Issue } from '@/api/issues/types'
  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { getIssues, updateIssue } from '@/api/issues/issues'
  import { getSpace } from '@/api/spaces/spaces'
  import Card from '@/components/common/Card.vue'
  import SpaceHeader from '@/components/SpaceHeader.vue'
  import SpaceSidebar from '@/components/SpaceSidebar.vue'
  import { useUserStore } from '@/stores/user'
  import { useNotificationStore } from '@/stores/notifications'

  const router = useRouter()
  const route = useRoute()
  const userStore = useUserStore()
  const notificationStore = useNotificationStore()

  const space = ref<Space>()
  const issues = ref<Issue[]>([])
  const loading = ref(true)
  const errorMessage = ref('')

  const openIssues = computed(() => {
    return issues.value.filter(issue => issue.status !== 'closed')
  })

  const closedIssues = computed(() => {
    return issues.value.filter(issue => issue.status === 'done' || issue.status === 'closed')
  })

  const filteredIssues = computed(() => {
    return issues.value.filter(issue => {
      const matchesSearch = searchQuery.value
        ? issue.title.toLowerCase().includes(searchQuery.value.toLowerCase()) || issue.id.toLowerCase().includes(searchQuery.value.toLowerCase())
        : true

      const matchesType = typeFilter.value ? issue.type === typeFilter.value : true
      const matchesPriority = priorityFilter.value ? issue.priority === priorityFilter.value : true
      const matchesStatus = statusFilter.value ? issue.status === statusFilter.value : true

      return matchesSearch && matchesType && matchesPriority && matchesStatus
    })
  })

  const displayType = ref<'board' | 'list'>('board')
  const searchQuery = ref('')
  const typeFilter = ref()
  const priorityFilter = ref()
  const statusFilter = ref()

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

      issues.value = (await getIssues(space.value.id, { limit: 200 })).items
      const savedDisplayType = localStorage.getItem(`issues_display_type_${space.value.id}`)
      if (savedDisplayType === 'board' || savedDisplayType === 'list') displayType.value = savedDisplayType

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

  // Watch for changes in displayType and save to localStorage
  watch(displayType, newType => {
    if (space.value) {
      localStorage.setItem(`issues_display_type_${space.value.id}`, newType)
    }
  })

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
        <IssueBoard
          v-if="displayType === 'board'"
          :issues="filteredIssues"
          :space="space!"
          @status-change="statusChanged"
        />

        <IssueTable
          v-else
          :issues="filteredIssues"
          :space="space!"
        />
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
