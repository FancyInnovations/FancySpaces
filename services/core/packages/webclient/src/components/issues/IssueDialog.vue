<script lang="ts" setup>

  import type { IssueComment } from '@/api/issues/types'
  import type { Space } from '@/api/spaces/types'
  import { addIssueComment, getIssueComments, updateIssue } from '@/api/issues/issues'
  import { getSpace } from '@/api/spaces/spaces'
  import Card from '@/components/common/Card.vue'
  import Dialog from '@/components/common/Dialog.vue'
  import IssueDialogSidebar from '@/components/issues/IssueDialogSidebar.vue'
  import { useIssueDialogStore } from '@/stores/issue-dialog'
  import { useNotificationStore } from '@/stores/notifications'
  import { useUserStore } from '@/stores/user'

  const notificationStore = useNotificationStore()
  const issueDialogStore = useIssueDialogStore()
  const userStore = useUserStore()
  const comments = ref<IssueComment[]>([])
  const commentText = ref('')
  const loadingComments = ref(false)
  const submittingComment = ref(false)
  const issueSpace = ref<Space>()

  const canWrite = computed(() => {
    const userID = userStore.user?.id
    if (!userID || !issueSpace.value) return false
    return issueSpace.value.creator === userID || issueSpace.value.members.some(member => member.user_id === userID && ['member', 'admin'].includes(member.role))
  })

  async function loadComments () {
    if (!issueDialogStore.issue) {
      comments.value = []
      return
    }
    loadingComments.value = true
    try {
      comments.value = await getIssueComments(issueDialogStore.issue.space, issueDialogStore.issue.id)
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to load comments.')
    } finally {
      loadingComments.value = false
    }
  }

  async function loadIssueContext () {
    issueSpace.value = undefined
    if (!issueDialogStore.issue) {
      comments.value = []
      return
    }
    try {
      issueSpace.value = await getSpace(issueDialogStore.issue.space)
    } catch {
      // The issue itself remains viewable if its space metadata is unavailable.
    }
    await loadComments()
  }

  watch(() => issueDialogStore.issue?.id, loadIssueContext, { immediate: true })

  async function statusChanged (newStatus: string) {
    if (!issueDialogStore.issue) return
    try {
      const updated = await updateIssue(issueDialogStore.issue.space, issueDialogStore.issue.id, { status: newStatus as any })
      issueDialogStore.issue = updated
      notificationStore.info('Issue status updated successfully')
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to update issue status.')
    }
  }

  async function addComment () {
    if (!issueDialogStore.issue || !commentText.value.trim()) return
    submittingComment.value = true
    try {
      comments.value.push(await addIssueComment(issueDialogStore.issue.space, issueDialogStore.issue.id, commentText.value))
      commentText.value = ''
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to add comment.')
    } finally {
      submittingComment.value = false
    }
  }

  function copyLink () {
    const issue = issueDialogStore.issue
    if (!issue) return
    navigator.clipboard.writeText(`${window.location.origin}/spaces/${issue.space}/issues/${issue.id}`)
    notificationStore.info('Issue link copied to clipboard!')
  }

  function copyID () {
    if (!issueDialogStore.issue) return
    navigator.clipboard.writeText(issueDialogStore.issue.id)
    notificationStore.info('Issue ID copied to clipboard!')
  }

  onMounted(async () => {
    await userStore.isAuthenticated
  })

</script>

<template>
  <Dialog :shown="issueDialogStore.isOpen" width="64%">
    <div class="rounded-xl">
      <div class="py-2 border-b d-flex align-center px-4">
        <h1 class="ml-2 text-h4 text-secondary">{{ issueDialogStore.issue?.title }}</h1>

        <div class="flex-grow-1 d-flex justify-end align-center">
          <v-select
            v-if="issueDialogStore.issue && canWrite"
            v-model="issueDialogStore.issue.status"
            class="mr-4"
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
            max-width="200"
            variant="solo"
            @update:model-value="statusChanged"
          />

          <v-btn
            class="mr-2"
            color="secondary"
            :href="`/spaces/${issueDialogStore.issue?.space}/issues/${issueDialogStore.issue?.id}`"
            icon="mdi-open-in-new"
            target="_blank"
            variant="text"
          />

          <v-btn color="secondary" icon="mdi-close" variant="text" @click="issueDialogStore.close()" />
        </div>
      </div>

      <div class="issue-dialog-inner d-flex">
        <IssueDialogSidebar class="ma-4" :comments="comments" :issue="issueDialogStore.issue!" />

        <div class="issue-dialog-inner pr-4 flex-grow-1">
          <Card class="mt-4 bg-transparent" color="#150D1950">
            <v-card-title class="mt-2">Description</v-card-title>
            <v-card-text><MarkdownRenderer class="issue-description" :markdown="issueDialogStore.issue?.description" /></v-card-text>
          </Card>

          <Card class="my-4 bg-transparent" color="#150D1950">
            <v-card-title class="mt-2">Comments ({{ comments.length }})</v-card-title>

            <v-card-text>
              <p v-if="loadingComments">Loading comments…</p>
              <p v-else-if="comments.length === 0">No comments yet.</p>

              <template v-else>
                <Card v-for="comment in comments" :key="comment.id" class="bg-transparent mb-3" elevation="6">
                  <v-card-text>
                    <div class="d-flex justify-space-between mb-2"><span class="font-weight-medium">{{ comment.author }}</span><span class="text-caption grey--text">{{ comment.created_at.toLocaleString() }}</span></div>
                    <MarkdownRenderer :markdown="comment.content" />
                  </v-card-text>
                </Card>
              </template>

              <div v-if="canWrite" class="mt-4">
                <v-textarea
                  v-model="commentText"
                  auto-grow
                  color="primary"
                  label="Add a comment"
                  rows="3"
                />

                <v-btn color="primary" :disabled="!commentText.trim() || submittingComment" :loading="submittingComment" @click="addComment">Add Comment</v-btn>
              </div>
            </v-card-text>
          </Card>
        </div>
      </div>

      <div class="d-flex justify-end pa-2 border-t">
        <v-btn class="mr-2" variant="text" @click="copyLink">Copy Link</v-btn>
        <v-btn class="mr-2" variant="text" @click="copyID">Copy ID</v-btn>

        <v-btn
          v-if="canWrite"
          class="mr-2"
          :to="`/spaces/${issueDialogStore.issue?.space}/issues/${issueDialogStore.issue?.id}/edit`"
          variant="text"
          @click="issueDialogStore.close()"
        >Edit</v-btn>
      </div>
    </div>
  </Dialog>
</template>

<style scoped>
.issue-dialog-inner {
  max-height: 75vh;
  overflow-y: auto;
  scrollbar-width: none;
}

.issue-description {
  max-height: 500px;
  overflow-y: auto;
  scrollbar-width: none;
}
</style>
