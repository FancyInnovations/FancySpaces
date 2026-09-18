<script lang="ts" setup>
  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { getSpace } from '@/api/spaces/spaces'
  import SpaceHeader from '@/components/SpaceHeader.vue'
  import SpaceSidebar from '@/components/SpaceSidebar.vue'
  import VersionForm from '@/components/versions/VersionForm.vue'
  import { useUserStore } from '@/stores/user'

  const route = useRoute()
  const router = useRouter()
  const userStore = useUserStore()
  const space = ref<Space>()

  onMounted(async () => {
    try {
      const requestedID = (route.params as any).sid as string
      space.value = await getSpace(requestedID)
      const userID = userStore.user?.id
      if (!userStore.isAuthenticated || !userID || !(space.value.creator === userID || space.value.members.some(member => member.user_id === userID))) {
        await router.push(`/spaces/${space.value.slug}`)
        return
      }
      useHead({ title: `New version - ${space.value.title} - FancySpaces` })
    } catch {
      await router.push('/')
    }
  })
</script>

<template>
  <v-container width="90%">
    <v-row>
      <v-col class="flex-grow-0 pa-0">
        <SpaceSidebar :space="space" />
      </v-col>

      <v-col>
        <SpaceHeader :space="space">
          <template #quick-actions>
            <v-btn color="primary" size="large" :to="`/spaces/${space?.slug}/versions`" variant="tonal">View Versions
            </v-btn>
          </template>
        </SpaceHeader>

        <hr class="grey-border-color mt-4">
      </v-col>
    </v-row>

    <VersionForm v-if="space" class="mt-8" :space-i-d="space.id" :space-slug="space.slug" />
    <v-progress-circular v-else class="d-block mx-auto mt-12" color="primary" indeterminate />
  </v-container>
</template>
