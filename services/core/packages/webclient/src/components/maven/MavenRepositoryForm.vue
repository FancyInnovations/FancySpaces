<script lang="ts" setup>
  import type { MavenRepositoryMutation, SpaceMavenRepository } from '@/api/maven/types'
  import { createMavenRepository, updateMavenRepository } from '@/api/maven/maven'
  import { useNotificationStore } from '@/stores/notifications'

  const props = defineProps<{
    spaceID: string
    spaceSlug: string
    repository?: SpaceMavenRepository
  }>()

  const router = useRouter()
  const notifications = useNotificationStore()
  const name = ref('')
  const isPublic = ref(true)
  const saving = ref(false)
  const error = ref('')

  watch(() => props.repository, repository => {
    if (!repository) return
    name.value = repository.name
    isPublic.value = repository.public
  }, { immediate: true })

  const isEditing = computed(() => !!props.repository)
  const valid = computed(() => name.value.trim().length > 0)

  async function save () {
    error.value = ''
    if (!valid.value) {
      error.value = 'Repository name is required.'
      return
    }

    saving.value = true
    try {
      const data: MavenRepositoryMutation = { name: name.value.trim(), public: isPublic.value }
      const repository = props.repository
        ? await updateMavenRepository(props.spaceID, props.repository.name, data)
        : await createMavenRepository(props.spaceID, data)

      notifications.info(isEditing.value ? 'Repository updated successfully' : 'Repository created successfully')
      await router.push(`/spaces/${props.spaceSlug}/maven-repos/${encodeURIComponent(repository.name)}`)
    } catch (error_) {
      error.value = error_ instanceof Error ? error_.message : 'Unable to save repository.'
      notifications.error(error.value)
    } finally {
      saving.value = false
    }
  }
</script>

<template>
  <v-row justify="center">
    <v-col md="8">
      <Card>
        <v-card-title class="mt-2">{{ isEditing ? 'Edit Maven Repository' : 'New Maven Repository' }}</v-card-title>

        <v-card-text>
          <v-alert v-if="error" class="mb-4" type="error" variant="tonal">{{ error }}</v-alert>

          <v-text-field
            v-model="name"
            class="mb-4"
            color="primary"
            :disabled="saving || isEditing"
            hint="Use letters, numbers, hyphens, or underscores."
            label="Repository name"
            persistent-hint
            required
          />

          <v-switch
            v-model="isPublic"
            color="primary"
            :disabled="saving"
            label="Public repository"
          />

          <div class="d-flex ga-3 mt-6">
            <v-btn :disabled="saving" :to="`/spaces/${props.spaceSlug}/maven-repos`" variant="text">Cancel</v-btn>

            <v-btn color="primary" :loading="saving" @click="save">
              {{ isEditing ? 'Save Changes' : 'Create Repository' }}
            </v-btn>
          </div>
        </v-card-text>
      </Card>
    </v-col>
  </v-row>
</template>
