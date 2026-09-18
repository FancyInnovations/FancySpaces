<script lang="ts" setup>

  import type { Issue } from '@/api/issues/types'
  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { deleteIssue, getIssue, updateIssue } from '@/api/issues/issues'
  import { getSpace } from '@/api/spaces/spaces'
  import IssueForm from '@/components/issues/IssueForm.vue'
  import SpaceHeader from '@/components/SpaceHeader.vue'
  import SpaceSidebar from '@/components/SpaceSidebar.vue'
  import { useNotificationStore } from '@/stores/notifications'
  import { useUserStore } from '@/stores/user'

  const router = useRouter()
  const route = useRoute()

  const space = ref<Space>()
  const issue = ref<Issue>()
  const loadError = ref('')
  const saving = ref(false)
  const userStore = useUserStore()
  const notificationStore = useNotificationStore()
  const issueForm = computed<Partial<Issue>>({
    get: () => issue.value || {},
    set: value => {
      if (issue.value) Object.assign(issue.value, value)
    },
  })

  const canWrite = computed(() => {
    const userID = userStore.user?.id
    if (!userID || !space.value) return false
    return space.value.creator === userID || space.value.members.some(member => member.user_id === userID && ['member', 'admin'].includes(member.role))
  })

  onMounted(async () => {
    try {
      await userStore.isAuthenticated
      const spaceID = (route.params as any).sid as string
      space.value = await getSpace(spaceID)
      if (!space.value.issue_settings.enabled) {
        await router.push(`/spaces/${space.value.slug}`)
        return
      }
      const issueID = (route.params as any).iid as string
      issue.value = await getIssue(spaceID, issueID)
      useHead({
        title: `${space.value.title} - FancySpaces`,
        meta: [{ name: 'description', content: space.value.summary || 'Edit this issue on FancySpaces.' }],
      })
    } catch (error) {
      loadError.value = error instanceof Error ? error.message : 'Failed to load issue.'
    }
  })

  async function editIssueReq () {
    if (!space.value || !issue.value || !canWrite.value || saving.value) return
    saving.value = true
    try {
      const saved = await updateIssue(space.value.id, issue.value.id, issue.value)
      issue.value = saved
      await router.push(`/spaces/${space.value.slug}/issues/${issue.value.id}`)
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to update issue.')
    } finally {
      saving.value = false
    }
  }

  async function deleteIssueReq () {
    if (!space.value || !issue.value || !canWrite.value) return
    try {
      await deleteIssue(space.value.id, issue.value.id)
      await router.push(`/spaces/${space.value.slug}/issues`)
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to archive issue.')
    }
  }

</script>

<template>
  <v-container v-if="issue && canWrite && !loadError" width="90%">
    <v-row>
      <v-col class="flex-grow-0 pa-0">
        <SpaceSidebar
          :space="space"
        />
      </v-col>

      <v-col>
        <SpaceHeader :space="space">
          <template #quick-actions>
            <v-btn
              color="primary"
              size="large"
              :to="`/spaces/${space?.slug}/issues`"
              variant="tonal"
            >
              View Issues
            </v-btn>
          </template>
        </SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        >
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <h1 class="text-center">Edit Issue #{{ issue?.id }}</h1>
      </v-col>
    </v-row>

    <v-row justify="center"><v-col md="6"><IssueForm v-model="issueForm" /></v-col></v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-btn
          class="mr-4"
          color="primary"
          :disabled="saving"
          :loading="saving"
          variant="tonal"
          @click="editIssueReq"
        >
          Edit Issue
        </v-btn>

        <v-btn
          color="error"
          variant="tonal"
          @click="deleteIssueReq"
        >
          Delete Issue
        </v-btn>
      </v-col>
    </v-row>
  </v-container>

  <v-container v-else class="text-center">
    {{ loadError || 'You do not have permission to edit this issue.' }}
  </v-container>
</template>

<style scoped>

</style>
