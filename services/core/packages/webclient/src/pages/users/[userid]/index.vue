<script lang="ts" setup>
  import type { User } from '@/api/auth/types'
  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { getPublicUser } from '@/api/auth/users'
  import { getDownloadCountForSpace, getSpacesOfCreator } from '@/api/spaces/spaces'
  import Card from '@/components/common/Card.vue'
  import UserHeader from '@/components/users/UserHeader.vue'
  import { useUserStore } from '@/stores/user'

  const route = useRoute()
  const userStore = useUserStore()

  const user = ref<User>()
  const spaces = ref<Space[]>([])
  const sortedSpaces = computed(() => {
    return spaces.value.sort((a, b) => b.created_at.getTime() - a.created_at.getTime())
  })

  const totalDownloads = ref(0)

  onMounted(async () => {
    const userID = (route.params as any).userid as string // username
    user.value = await getPublicUser(userID)

    spaces.value = await getSpacesOfCreator(user.value.id)

    for (const sp of spaces.value) {
      totalDownloads.value += await getDownloadCountForSpace(sp.id)
    }

    useHead({
      title: `${user.value.name} - FancySpaces`,
      meta: [
        {
          name: 'description',
          content: `Check out ${user.value.name}'s profile on FancySpaces, showcasing their spaces.`,
        },
      ],
    })
  })

</script>

<template>
  <v-container v-if="user" width="60%">
    <v-row>
      <v-col>
        <UserHeader :user="user">
          <template #quick-actions>
            <v-btn
              v-if="user?.id === userStore.user?.id"
              class="sidebar__mobile"
              color="primary"
              size="large"
              :to="`/spaces/new`"
              variant="tonal"
            >
              New Space
            </v-btn>
          </template>
        </UserHeader>

        <hr
          class="mt-4 grey-border-color"
        >
      </v-col>
    </v-row>

    <v-row>
      <v-col md="8">
        <template v-for="space in sortedSpaces" :key="space.id">
          <SpaceCard
            class="mb-4"
            :space="space"
            :with-badge="true"
          />
        </template>
      </v-col>

      <v-col md="4">
        <Card
          class="mb-4"
          min-width="200"
        >
          <v-card-title class="mt-2">Details</v-card-title>

          <v-card-text>
            <p class="text-body-1"><strong>ID:</strong> {{ user?.id }}</p>
            <p class="text-body-1"><strong>Username:</strong> {{ user?.name }}</p>
            <p class="text-body-1"><strong>Roles:</strong> {{ user?.roles.join(", ") }}</p>
            <p class="text-body-1"><strong>Joined at:</strong> {{ user?.created_at.toLocaleDateString() }}</p>
            <p class="text-body-1"><strong>Spaces:</strong> {{ spaces.length }}</p>
            <p class="text-body-1"><strong>Total Downloads:</strong> {{ totalDownloads }}</p>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>

  <v-container v-else width="60%">
    <v-row>
      <v-col>
        <p class="text-h4 text-center mt-8">User not found or still loading...</p>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>

<style>

</style>
