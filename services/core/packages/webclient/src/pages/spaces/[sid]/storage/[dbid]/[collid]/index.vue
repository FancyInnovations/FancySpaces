<script lang="ts" setup>

  import type { Space } from '@/api/spaces/types'
  import type { SpaceDatabase, SpaceDatabaseCollection } from '@/api/storage/types'
  import { useHead } from '@vueuse/head'
  import { getSpace } from '@/api/spaces/spaces'
  import SpaceHeader from '@/components/SpaceHeader.vue'
  import SpaceSidebar from '@/components/SpaceSidebar.vue'
  import KVCollectionDataPage from '@/components/storage/KVCollectionDataPage.vue'
  import { useUserStore } from '@/stores/user'

  const router = useRouter()
  const route = useRoute()
  const userStore = useUserStore()

  const isLoggedIn = ref(false)

  const space = ref<Space>()
  const database = ref<SpaceDatabase>()
  const collection = ref<SpaceDatabaseCollection>()

  onMounted(async () => {
    isLoggedIn.value = await userStore.isAuthenticated

    const spaceID = (route.params as any).sid as string
    space.value = await getSpace(spaceID)

    if (!space.value.maven_repository_settings.enabled) {
      router.push(`/spaces/${space.value.slug}`)
      return
    }

    const databaseName = (route.params as any).dbid as string
    database.value = {
      name: databaseName,
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
    }

    const collectionName = (route.params as any).collid as string
    collection.value = {
      database: databaseName,
      name: collectionName,
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: 'kv',
    }

    useHead({
      title: `${space.value.title} storage - FancySpaces`,
      meta: [
        {
          name: 'description',
          content: space.value.summary || `Explore the ${space.value.title} project space on FancySpaces.`,
        },
      ],
    })
  })
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
              v-if="isLoggedIn"
              color="primary"
              disabled
              size="large"
              :to="`/spaces/${space?.slug}/storage/new`"
              variant="tonal"
            >
              New collection
            </v-btn>
          </template>
        </SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        >
      </v-col>
    </v-row>

    <KVCollectionDataPage
      v-if="space && database && collection && collection.engine === 'kv'"
      :collection="collection"
      :database="database"
      :space="space"
    />
  </v-container>
</template>

<style scoped>

</style>
