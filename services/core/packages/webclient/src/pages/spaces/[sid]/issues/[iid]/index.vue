<script lang="ts" setup>

  import type { Issue, IssueActivity, IssueComment } from '@/api/issues/types'
  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { addIssueComment, deleteIssue, deleteIssueComment, getIssue, getIssueActivity, getIssueComments, updateIssue, updateIssueComment } from '@/api/issues/issues'
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
  const activity = ref<IssueActivity[]>([])
  const loading = ref(true)
  const errorMessage = ref('')
  const commentText = ref('')
  const commentSubmitting = ref(false)
  const editingComment = ref<string>()
  const editedComment = ref('')
  const relationshipIssue = ref('')
  const relationshipType = ref<'blocks' | 'duplicates' | 'related'>('related')

  const canWrite = computed(() => {
    const userID = userStore.user?.id
    if (!userID || !space.value) return false
    return space.value.creator === userID || space.value.members.some(member => member.user_id === userID && ['member', 'admin'].includes(member.role))
  })
  const canModerateComments = computed(() => {
    const userID = userStore.user?.id
    if (!userID || !space.value) return false
    return space.value.creator === userID || space.value.members.some(member => member.user_id === userID && member.role === 'admin')
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
      const details = await Promise.all([getIssueComments(spaceID, issueID), getIssueActivity(spaceID, issueID)])
      comments.value = details[0]
      activity.value = details[1]
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

  function startEditComment (comment: IssueComment) {
    editingComment.value = comment.id
    editedComment.value = comment.content
  }

  async function saveComment (comment: IssueComment) {
    if (!space.value || !currentIssue.value || !editedComment.value.trim()) return
    try {
      const saved = await updateIssueComment(space.value.id, currentIssue.value.id, comment.id, editedComment.value)
      Object.assign(comment, saved)
      editingComment.value = undefined
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to update comment.')
    }
  }

  async function removeComment (comment: IssueComment) {
    if (!space.value || !currentIssue.value) return
    try {
      await deleteIssueComment(space.value.id, currentIssue.value.id, comment.id)
      comments.value = comments.value.filter(item => item.id !== comment.id)
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to delete comment.')
    }
  }

  async function addRelationship () {
    if (!currentIssue.value || !relationshipIssue.value.trim()) return
    const relationship = { issue: relationshipIssue.value.trim(), type: relationshipType.value }
    if (relationship.issue === currentIssue.value.id || currentIssue.value.relationships?.some(item => item.issue === relationship.issue && item.type === relationship.type)) return
    try {
      const saved = await updateIssue(currentIssue.value.space, currentIssue.value.id, { relationships: [...(currentIssue.value.relationships || []), relationship] })
      Object.assign(currentIssue.value, saved)
      relationshipIssue.value = ''
      activity.value = await getIssueActivity(currentIssue.value.space, currentIssue.value.id)
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to add relationship.')
    }
  }

  async function removeRelationship (issue: string, type: string) {
    if (!currentIssue.value) return
    try {
      const saved = await updateIssue(currentIssue.value.space, currentIssue.value.id, { relationships: (currentIssue.value.relationships || []).filter(item => item.issue !== issue || item.type !== type) })
      Object.assign(currentIssue.value, saved)
    } catch (error) {
      notificationStore.error(error instanceof Error ? error.message : 'Failed to remove relationship.')
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
            <v-btn
              class="mb-2"
              color="primary"
              size="large"
              :to="`/spaces/${space?.slug}/issues`"
              variant="tonal"
            >View Issues</v-btn>

            <v-btn
              v-if="canWrite"
              color="primary"
              size="large"
              :to="`/spaces/${space?.slug}/issues/new`"
              variant="tonal"
            >New Issue</v-btn>
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

                <v-textarea v-if="editingComment === comment.id" v-model="editedComment" auto-grow density="compact" />
                <MarkdownRenderer v-else :markdown="comment.content" />

                <div v-if="canWrite && (comment.author === userStore.user?.id || canModerateComments)" class="mt-2 d-flex ga-2">
                  <v-btn v-if="editingComment === comment.id" color="primary" size="small" @click="saveComment(comment)">Save</v-btn>
                  <v-btn v-else size="small" variant="text" @click="startEditComment(comment)">Edit</v-btn>
                  <v-btn v-if="editingComment === comment.id" size="small" variant="text" @click="editingComment = undefined">Cancel</v-btn>
                  <v-btn color="error" size="small" variant="text" @click="removeComment(comment)">Delete</v-btn>
                </div>
              </v-card-text>
            </Card>

            <div v-if="canWrite" class="mt-4">
              <v-textarea
                v-model="commentText"
                auto-grow
                color="primary"
                label="Add a comment"
                rows="3"
              />

              <v-btn color="primary" :disabled="!commentText.trim() || commentSubmitting" :loading="commentSubmitting" @click="addComment">Add Comment</v-btn>
            </div>
          </v-card-text>
        </Card>
      </v-col>

      <v-col>
        <IssueDialogSidebar :comments="comments" :issue="currentIssue" />

        <Card class="bg-transparent mt-4" elevation="6">
          <v-card-title>Relationships</v-card-title>

          <v-card-text>
            <div v-for="relationship in currentIssue.relationships || []" :key="`${relationship.type}-${relationship.issue}`" class="d-flex justify-space-between align-center mb-2">
              <span>{{ relationship.type }} <router-link :to="`/spaces/${space?.slug}/issues/${relationship.issue}`">#{{ relationship.issue }}</router-link></span>

              <v-btn
                v-if="canWrite"
                icon="mdi-close"
                size="x-small"
                variant="text"
                @click="removeRelationship(relationship.issue, relationship.type)"
              />
            </div>

            <template v-if="canWrite">
              <v-text-field
                v-model="relationshipIssue"
                class="mb-2"
                density="compact"
                hide-details
                label="Issue ID"
              />

              <v-select
                v-model="relationshipType"
                class="mb-2"
                density="compact"
                hide-details
                :items="['blocks', 'duplicates', 'related']"
              />

              <v-btn size="small" @click="addRelationship">Link issue</v-btn>
            </template>
          </v-card-text>
        </Card>

        <Card class="bg-transparent mt-4" elevation="6">
          <v-card-title>Activity</v-card-title>

          <v-card-text>
            <p v-if="activity.length === 0" class="text-caption">No activity yet.</p>
            <p v-for="entry in activity" :key="entry.id" class="text-caption mb-2"><strong>{{ entry.actor }}</strong> {{ entry.kind }}<template v-if="entry.field"> {{ entry.field }} from “{{ entry.old_value || 'none' }}” to “{{ entry.new_value || 'none' }}”</template> · {{ entry.created_at.toLocaleString() }}</p>
          </v-card-text>
        </Card>

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

            <v-btn
              block
              class="mb-2"
              color="primary"
              :to="`/spaces/${space?.slug}/issues/${currentIssue.id}/edit`"
              variant="tonal"
            >Edit Issue</v-btn>

            <v-btn block color="error" variant="tonal" @click="deleteIssueReq">Archive Issue</v-btn>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>

  <v-container v-else-if="loading" class="text-center"><v-progress-circular color="primary" indeterminate /></v-container>
  <v-container v-else class="text-center"><p class="mb-4">{{ errorMessage || 'Issue not found.' }}</p><v-btn color="primary" @click="load">Retry</v-btn></v-container>
</template>
