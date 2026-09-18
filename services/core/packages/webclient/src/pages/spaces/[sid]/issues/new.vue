<script lang="ts" setup>

  import type { Issue } from '@/api/issues/types'
  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { createIssue, getIssues } from '@/api/issues/issues'
  import { getSpace } from '@/api/spaces/spaces'
  import IssueForm from '@/components/issues/IssueForm.vue'
  import SpaceHeader from '@/components/SpaceHeader.vue'
  import SpaceSidebar from '@/components/SpaceSidebar.vue'
  import { useNotificationStore } from '@/stores/notifications'
  import { useUserStore } from '@/stores/user'

  const router = useRouter()
  const route = useRoute()
  const notifications = useNotificationStore()
  const userStore = useUserStore()

  const space = ref<Space>()

  const draft = ref<Partial<Issue>>({ title: '', description: '', type: 'task', priority: 'medium', labels: [] })
  const submitting = ref(false)
  const loadError = ref('')
  const templates = ref<Array<{ name: string, issue: Partial<Issue> }>>([])
  const templateName = ref('')

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
      templates.value = JSON.parse(localStorage.getItem(`issue_templates_${space.value.id}`) || '[]')
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

    if (!draft.value.title?.trim()) {
      notifications.error('Title is required.')
      return
    }

    if (!draft.value.type) {
      notifications.error('Type is required.')
      return
    }

    if (!draft.value.priority) {
      notifications.error('Priority is required.')
      return
    }

    submitting.value = true
    try {
      const similar = await getIssues(space.value.id, { q: draft.value.title, limit: 5 })
      if (similar.items.some(item => item.title.trim().toLowerCase() === draft.value.title!.trim().toLowerCase())) {
        notifications.error('A similarly named issue already exists. Review it before creating a duplicate.')
        return
      }
      const issue = await createIssue(space.value.id, draft.value)
      draft.value = { title: '', description: '', type: 'task', priority: 'medium', labels: [] }

      await router.push(`/spaces/${space.value.slug}/issues/${issue.id}`)
    } catch (error) {
      notifications.error(error instanceof Error ? error.message : 'Failed to create issue.')
    } finally {
      submitting.value = false
    }
  }

  function applyTemplate (name: string | null) {
    if (!name) return
    const template = templates.value.find(item => item.name === name)
    if (template) draft.value = { ...template.issue, labels: [...(template.issue.labels || [])] }
  }

  function saveTemplate () {
    if (!space.value || !templateName.value.trim()) return
    const issue = { ...draft.value, title: '', description: '', labels: [...(draft.value.labels || [])] }
    templates.value = [...templates.value.filter(item => item.name !== templateName.value.trim()), { name: templateName.value.trim(), issue }]
    localStorage.setItem(`issue_templates_${space.value.id}`, JSON.stringify(templates.value))
    templateName.value = ''
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
            <div class="d-flex ga-2 mb-4 flex-wrap">
              <v-select
                v-if="templates.length > 0"
                density="compact"
                hide-details
                :items="templates.map(template => template.name)"
                label="Start from template"
                @update:model-value="applyTemplate"
              />

              <v-text-field v-model="templateName" density="compact" hide-details label="Save as template" />
              <v-btn size="small" variant="tonal" @click="saveTemplate">Save template</v-btn>
            </div>

            <IssueForm v-model="draft" />

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
