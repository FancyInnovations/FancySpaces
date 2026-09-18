<script lang="ts" setup>

  import type { Issue, IssueComment } from '@/api/issues/types'
  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { addIssueComment, deleteIssue, getIssue, getIssueComments, updateIssue } from '@/api/issues/issues'
  import { getSpace } from '@/api/spaces/spaces'
  import Card from '@/components/common/Card.vue'
  import IssueDialogSidebar from '@/components/issues/IssueDialogSidebar.vue'
  import SpaceHeader from '@/components/SpaceHeader.vue'
  import SpaceSidebar from '@/components/SpaceSidebar.vue'
  import { useConfirmationStore } from '@/stores/confirmation'
  import { useNotificationStore } from '@/stores/notifications'
  import { useUserStore } from '@/stores/user'

  const route = useRoute()
  const router = useRouter()
  const confirmationStore = useConfirmationStore()
  const notificationStore = useNotificationStore()
  const userStore = useUserStore()
  const space = ref<Space>()
  const currentIssue = ref<Issue>()
  const comments = ref<IssueComment[]>([])
  const loading = ref(true)
  const errorMessage = ref('')
  const commentText = ref('')
  const commentSubmitting = ref(false)

  const canWrite = computed(() => {
    const userID = userStore.user?.id
    if (!userID || !space.value) return false
    return space.value.creator === userID || space.value.members.some(member => member.user_id === userID && ['member', 'admin'].includes(member.role))
  })

  async function load () {
    loading.value = true
    errorMessage.value = ''
    try {
      const spaceID = (route.params as any).sid as string
      space.value = await getSpace(spaceID)
      if (!space.value.issue_settings.enabled) {
        await router.push(`/spaces/${space.value.slug}`)
        return
      }
      const issueID = (route.params as any).iid as string
      currentIssue.value = await getIssue(spaceID, issueID)
      comments.value = await getIssueComments(spaceID, issueID)
      useHead({
        title: `${space.value.title} Issue #${currentIssue.value.id} - FancySpaces`,
        meta: [{ name: 'description', content: space.value.summary || `View issue #${currentIssue.value.id} for this space on FancySpaces.` }],
      })
    } catch (error) {
      errorMessage.value = error instanceof Error ? error.message : 'Failed to load issue.'
    } finally {
      loading.value = false
    }
  }

  onMounted(async () => {
    await userStore.isAuthenticated
    await load()
  })

  async function deleteIssueReq () {
    if (!space.value || !currentIssue.value) return
    confirmationStore.confirmation = {
      shown: true,
      persistent: true,
      title: 'Archive Issue',
      text: 'Are you sure you want to archive this issue? It will no longer appear in the issue list.',
      yesText: 'Archive',
      onConfirm: async () => {
        try {
          await deleteIssue(space.value!.id, currentIssue.value!.id)
          notificationStore.info('Issue archived successfully')
          await router.push(`/spaces/${space.value!.slug}/issues`)
        } catch (error) {
          notificationStore.error(error instanceof Error ? error.message : 'Failed to archive issue.')
        }
      },
    }
  }

  async function statusChanged (newStatus: string) {
    if (!currentIssue.value) return
    try {
      const saved = await updateIssue(currentIssue.value.space, currentIssue.value.id, { status: newStatus as Issue['status'] })
      Object.assign(currentIssue.value, saved)
      notificationStore.info('Issue status updated successfully')
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to update issue status.')
    }
  }

  async function addComment () {
    if (!space.value || !currentIssue.value || !commentText.value.trim()) return
    commentSubmitting.value = true
    try {
      comments.value.push(await addIssueComment(space.value.id, currentIssue.value.id, commentText.value))
      commentText.value = ''
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to add comment.')
    } finally {
      commentSubmitting.value = false
    }
  }

</script>

<template>
  <v-container v-if="!loading && currentIssue" width="90%">
    <v-row>
      <v-col class="flex-grow-0 pa-0"><SpaceSidebar :space="space" /></v-col>
      <v-col>
        <SpaceHeader :space="space">
          <template #quick-actions>
            <v-btn class="mb-2" color="primary" size="large" :to="`/spaces/${space?.slug}/issues`" variant="tonal">View Issues</v-btn>
            <v-btn v-if="canWrite" color="primary" size="large" :to="`/spaces/${space?.slug}/issues/new`" variant="tonal">New Issue</v-btn>
          </template>
        </SpaceHeader>
        <hr class="grey-border-color mt-4">
      </v-col>
    </v-row>

    <v-row>
      <v-col md="9">
        <Card class="mb-4"><v-card-title class="my-2">{{ currentIssue.title }}</v-card-title></Card>
        <Card class="mb-4"><v-card-text class="py-0"><MarkdownRenderer :markdown="currentIssue.description" /></v-card-text></Card>
        <Card class="bg-transparent" color="#150D1950">
          <v-card-title class="mt-2">Comments ({{ comments.length }})</v-card-title>
          <v-card-text>
            <p v-if="comments.length === 0">No comments yet.</p>
            <Card v-for="comment in comments" :key="comment.id" class="bg-transparent mb-3" elevation="6">
              <v-card-text>
                <div class="d-flex justify-space-between mb-2">
                  <span class="font-weight-medium">{{ comment.author }}</span>
                  <span class="text-caption grey--text">{{ comment.created_at.toLocaleString() }}</span>
                </div>
                <MarkdownRenderer :markdown="comment.content" />
              </v-card-text>
            </Card>
            <div v-if="canWrite" class="mt-4">
              <v-textarea v-model="commentText" auto-grow color="primary" label="Add a comment" rows="3" />
              <v-btn color="primary" :disabled="!commentText.trim() || commentSubmitting" :loading="commentSubmitting" @click="addComment">Add Comment</v-btn>
            </div>
          </v-card-text>
        </Card>
      </v-col>

      <v-col>
        <IssueDialogSidebar :comments="comments" :issue="currentIssue" />
        <Card v-if="canWrite" class="bg-transparent mt-4" elevation="6">
          <v-card-text>
            <v-select
              v-model="currentIssue.status"
              class="mb-4"
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
              rounded="xl"
              @update:model-value="statusChanged"
            />
            <v-btn block class="mb-2" color="primary" :to="`/spaces/${space?.slug}/issues/${currentIssue.id}/edit`" variant="tonal">Edit Issue</v-btn>
            <v-btn block color="error" variant="tonal" @click="deleteIssueReq">Archive Issue</v-btn>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>

  <v-container v-else-if="loading" class="text-center"><v-progress-circular color="primary" indeterminate /></v-container>
  <v-container v-else class="text-center"><p class="mb-4">{{ errorMessage || 'Issue not found.' }}</p><v-btn color="primary" @click="load">Retry</v-btn></v-container>
</template>
