<script lang="ts" setup>

  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { createIssue } from '@/api/issues/issues'
  import { getSpace } from '@/api/spaces/spaces'
  import SpaceHeader from '@/components/SpaceHeader.vue'
  import SpaceSidebar from '@/components/SpaceSidebar.vue'
  import { useNotificationStore } from '@/stores/notifications'
  import { useUserStore } from '@/stores/user'

  const router = useRouter()
  const route = useRoute()
  const notifications = useNotificationStore()
  const userStore = useUserStore()

  const space = ref<Space>()

  const title = ref('')
  const description = ref('')
  const type = ref('task')
  const priority = ref('medium')
  const submitting = ref(false)
  const loadError = ref('')

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
      useHead({
        title: `${space.value.title} - FancySpaces`,
        meta: [{ name: 'description', content: space.value.summary || 'Create a new issue in this space on FancySpaces.' }],
      })
    } catch (error) {
      loadError.value = error instanceof Error ? error.message : 'Failed to load space.'
    }
  })

  async function createNewIssue () {
    if (!space.value || !canWrite.value || submitting.value) return

    if (!title.value.trim()) {
      notifications.error('Title is required.')
      return
    }

    if (!type.value) {
      notifications.error('Type is required.')
      return
    }

    if (!priority.value) {
      notifications.error('Priority is required.')
      return
    }

    submitting.value = true
    try {
      const issue = await createIssue(space.value.id, {
        title: title.value,
        description: description.value,
        type: type.value as any,
        priority: priority.value as any,
      })

      title.value = ''
      description.value = ''
      type.value = 'task'
      priority.value = 'medium'

      await router.push(`/spaces/${space.value.slug}/issues/${issue.id}`)
    } catch (error) {
      notifications.error(error instanceof Error ? error.message : 'Failed to create issue.')
    } finally {
      submitting.value = false
    }
  }

</script>

<template>
  <v-container v-if="!loadError" width="90%">
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

    <v-row justify="center">
      <v-col md="8">
        <Card>
          <v-card-title class="mt-2">
            New Issue
          </v-card-title>

          <v-card-text>
            <v-text-field
              v-model="title"
              class="mb-4"
              color="primary"
              hide-details
              label="Title"
              required
            />

            <v-textarea
              v-model="description"
              class="mb-4"
              color="primary"
              hide-details
              label="Description"
              rows="8"
            />

            <div class="d-flex mb-4">
              <v-select
                v-model="type"
                class="mr-2"
                color="primary"
                hide-details
                :items="[
                  { title: 'Epic', value: 'epic' },
                  { title: 'Bug', value: 'bug' },
                  { title: 'Task', value: 'task' },
                  { title: 'Story', value: 'story' },
                  { title: 'Idea', value: 'idea' },

                ]"
                label="Type"
                required
              />

              <v-select
                v-model="priority"
                class="ml-2"
                color="primary"
                hide-details
                :items="[
                  { title: 'Low', value: 'low' },
                  { title: 'Medium', value: 'medium' },
                  { title: 'High', value: 'high' },
                  { title: 'Critical', value: 'critical' },
                ]"
                label="Priority"
                required
              />
            </div>

            <v-btn
              class="mt-4"
              color="primary"
              :disabled="!canWrite || submitting"
              :loading="submitting"
              @click="createNewIssue"
            >
              Create Issue
            </v-btn>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>
  <v-container v-else class="text-center">{{ loadError }}</v-container>
</template>

<style scoped>

</style>
