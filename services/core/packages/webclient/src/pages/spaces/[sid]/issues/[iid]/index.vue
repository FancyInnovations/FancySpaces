<script lang="ts" setup>

  import type { Issue, IssueComment } from '@/api/issues/types'
  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { deleteIssue, getIssue, updateIssue } from '@/api/issues/issues'
  import { getSpace } from '@/api/spaces/spaces'
  import Card from '@/components/common/Card.vue'
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

  const isLoggedIn = ref(false)

  const space = ref<Space>()

  const currentIssue = ref<Issue>()
  const comments = ref<IssueComment[]>([])

  onMounted(async () => {
    isLoggedIn.value = await userStore.isAuthenticated

    const spaceID = (route.params as any).sid as string
    space.value = await getSpace(spaceID)

    if (!space.value.issue_settings.enabled) {
      router.push(`/spaces/${space.value.slug}`)
      return
    }

    const issueID = (route.params as any).iid as string
    currentIssue.value = await getIssue(spaceID, issueID)

    useHead({
      title: `${space.value.title} Issue #${currentIssue.value?.id} - FancySpaces`,
      meta: [
        {
          name: 'description',
          content: space.value.summary || `View issue #${currentIssue.value?.id} for this space on FancySpaces.`,
        },
      ],
    })
  })

  async function deleteIssueReq () {
    if (!space.value || !currentIssue.value) return

    confirmationStore.confirmation = {
      shown: true,
      persistent: true,
      title: 'Delete Issue',
      text: 'Are you sure you want to delete this issue? This action cannot be undone.',
      yesText: 'Delete',
      onConfirm: async () => {
        if (!space.value || !currentIssue.value) return

        await deleteIssue(space.value!.id, currentIssue.value.id)
        notificationStore.info('Issue deleted successfully')
        await router.push(`/spaces/${space.value?.slug}/issues`)
      },
    }
  }

  async function statusChanged (newStatus: string) {
    currentIssue.value!.status = newStatus as any
    await updateIssue(currentIssue.value!.space, currentIssue.value!.id, currentIssue.value!)
    notificationStore.info('Issue status updated successfully')
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
          <template #quick-actions>
            <v-btn
              class="mb-2"
              color="primary"
              exact
              size="large"
              :to="`/spaces/${space?.slug}/issues`"
              variant="tonal"
            >
              View Issues
            </v-btn>

            <v-btn
              v-if="isLoggedIn"
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

    <v-row>
      <v-col md="9">
        <Card class="mb-4">
          <v-card-title class="my-2">{{ currentIssue?.title }}</v-card-title>
        </Card>

        <Card class="mb-4">
          <v-card-text class="py-0">
            <MarkdownRenderer
              :markdown="currentIssue?.description"
            />
          </v-card-text>
        </Card>

        <Card
          class="bg-transparent"
          color="#150D1950"
          min-width="600"
        >
          <v-card-title class="mt-2">
            Comments ({{ comments?.length }})
          </v-card-title>

          <v-card-text>
            <p v-if="comments?.length === 0">
              No comments yet.
            </p>

            <div v-else class="issue-comments">
              <Card
                v-for="comment in comments"
                :key="comment.id"
                class="bg-transparent mb-3"
                elevation="6"
              >
                <v-card-text>
                  <div class="d-flex justify-space-between mb-2">
                    <div class="d-flex align-center">
                      <span class="font-weight-medium">{{ comment.author }}</span>
                    </div>

                    <span class="text-caption grey--text">{{ comment.created_at.toLocaleString() }}</span>
                  </div>

                  <MarkdownRenderer
                    :markdown="comment.content"
                  />
                </v-card-text>
              </Card>
            </div>
          </v-card-text>
        </Card>
      </v-col>

      <v-col>
        <IssueDialogSidebar
          :comments="comments"
          :issue="currentIssue"
        />

        <Card
          v-if="isLoggedIn"
          class="bg-transparent mt-4"
          elevation="6"
        >
          <v-card-text>
            <v-select
              v-if="currentIssue"
              v-model="currentIssue!.status"
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
              :to="`/spaces/${space?.slug}/issues/${currentIssue?.id}/edit`"
              variant="tonal"
            >
              Edit Issue
            </v-btn>

            <v-btn
              block
              color="error"
              variant="tonal"
              @click="deleteIssueReq"
            >
              Delete Issue
            </v-btn>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
