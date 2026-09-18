<script lang="ts" setup>

  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { getDownloadCountForSpace, getDownloadCountForSpacePerVersion, getSpace } from '@/api/spaces/spaces'
  import { mapPlatformToDisplayname, type SpaceVersion } from '@/api/versions/types'
  import { deleteVersion, getAllVersions, getLatestVersion } from '@/api/versions/versions'
  import Card from '@/components/common/Card.vue'
  import SpaceHeader from '@/components/SpaceHeader.vue'
  import SpaceSidebar from '@/components/SpaceSidebar.vue'
  import VersionChannelChip from '@/components/versions/VersionChannelChip.vue'
  import { useConfirmationStore } from '@/stores/confirmation'
  import { useNotificationStore } from '@/stores/notifications'
  import { useUserStore } from '@/stores/user'

  const router = useRouter()
  const route = useRoute()
  const notificationStore = useNotificationStore()
  const confirmationStore = useConfirmationStore()
  const userStore = useUserStore()

  const isMember = computed(() => {
    if (!space.value) return false
    if (!userStore.isAuthenticated) return false

    const userID = userStore.user?.id
    return space.value?.creator == userID || space.value?.members.some(member => member.user_id === userID)
  })

  const space = ref<Space>()
  const latestVersion = ref<SpaceVersion>()
  const versions = ref<SpaceVersion[]>()
  const downloadCount = ref<number>(0)
  const downloadCounts = ref<Record<string, number>>({})
  const loading = ref(true)
  const loadError = ref('')

  const possiblePlatforms = computed(() => {
    const platforms = new Set<string>()
    if (versions.value) for (const ver of versions.value) {
      platforms.add(ver.platform)
    }
    return Array.from(platforms)
  })

  const filterChannel = ref<string[]>([])
  const filterPlatform = ref<string[]>([])
  const filteredVersions = computed(() => {
    return versions.value?.filter(ver => {
      const channelMatch = filterChannel.value.length === 0 || filterChannel.value.includes(ver.channel)
      const platformMatch = filterPlatform.value.length === 0 || filterPlatform.value.includes(ver.platform)
      return channelMatch && platformMatch
    }) || []
  })

  const tableHeaders = [
    { title: 'Version', key: 'name', sortable: false },
    { title: 'Channel', key: 'channel', sortable: false },
    {
      title: 'Platform',
      key: 'platform',
      value: (ver: SpaceVersion) => mapPlatformToDisplayname(ver.platform),
      sortable: false,
    },
    {
      title: 'Platform versions',
      key: 'supported_platform_versions',
      sortable: false,
      value: (ver: SpaceVersion) => ver.supported_platform_versions.join(', '),
      class: 'platform-versions__max-width',
    },
    {
      title: 'Released at',
      key: 'published_at',
      sortable: false,
      value: (ver: SpaceVersion) => new Date(ver.published_at).toLocaleString(),
    },
    {
      title: 'Downloads',
      key: 'downloads',
      sortable: false,
      value: (ver: SpaceVersion) => downloadCounts.value[ver.id] || 0,
    },
    { title: '', key: 'actions', sortable: false, align: 'end' as any },
  ]

  onMounted(async () => {
    try {
      const spaceID = (route.params as any).sid as string
      space.value = await getSpace(spaceID)

      if (!space.value.release_settings.enabled) {
        await router.push(`/spaces/${space.value.slug}`)
        return
      }

      latestVersion.value = await getLatestVersion(space.value.id).catch(() => undefined)
      versions.value = await getAllVersions(space.value.id)
      downloadCount.value = await getDownloadCountForSpace(space.value.id)
      downloadCounts.value = await getDownloadCountForSpacePerVersion(space.value.id)

      useHead({
        title: `${space.value.title} versions - FancySpaces`,
        meta: [{
          name: 'description',
          content: space.value.summary || `Explore the ${space.value.title} project space on FancySpaces.`,
        }],
      })
    } catch (error) {
      loadError.value = error instanceof Error ? error.message : 'Unable to load versions.'
      notificationStore.error(loadError.value)
    } finally {
      loading.value = false
    }
  })

  function onRowClick (event: any, { item }: any) {
    if (event?.target?.closest?.('a, button')) return
    router.push(`/spaces/${space.value?.slug}/versions/${item.name}`)
  }

  function deleteVersionReq (evt: any, v: SpaceVersion) {
    evt.stopPropagation()

    confirmationStore.confirmation = {
      shown: true,
      persistent: true,
      title: 'Delete version',
      text: 'Are you sure you want to delete this version? This action cannot be undone.',
      yesText: 'Delete',
      onConfirm: async () => {
        await deleteVersion(v.space_id, v.id)

        versions.value = versions.value?.filter(ver => ver.id !== v.id)
        notificationStore.info('Version deleted')
      },
    }
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
          <template #metadata>
            <p class="text-body-2 mx-4">-</p>
            <p class="text-body-2">{{ versions?.length }} versions</p>

            <p class="text-body-2 mx-4">-</p>
            <p class="text-body-2">{{ downloadCount }} downloads</p>

            <p v-if="latestVersion" class="text-body-2 mx-4">-</p>
            <p v-if="latestVersion" class="text-body-2">Latest version: {{ latestVersion?.name }}</p>
          </template>

          <template #quick-actions>
            <v-btn
              v-if="latestVersion && latestVersion.files.length != 1"
              class="sidebar__mobile"
              color="primary"
              prepend-icon="mdi-download"
              size="large"
              :to="`/spaces/${space?.slug}/versions/latest`"
              variant="tonal"
            >
              latest
            </v-btn>

            <v-btn
              v-else
              class="sidebar__mobile"
              color="primary"
              :href="latestVersion?.files[0]?.url"
              prepend-icon="mdi-download"
              size="large"
              variant="tonal"
            >
              Latest
            </v-btn>

            <v-btn
              v-if="isMember"
              class="sidebar__mobile mt-4"
              color="primary"
              prepend-icon="mdi-plus"
              size="large"
              :to="`/spaces/${space?.slug}/versions/new`"
              variant="tonal"
            >
              New version
            </v-btn>
          </template>
        </SpaceHeader>

        <hr
          class="mt-4 grey-border-color"
        >
      </v-col>
    </v-row>

    <v-row>
      <v-col md="6">
        <Card>
          <v-card-title class="mt-2">
            Filter
          </v-card-title>

          <v-card-text>
            <v-container class="pa-2">
              <v-row>
                <v-col>
                  <v-select
                    v-model="filterChannel"
                    chips
                    clearable
                    color="primary"
                    density="compact"
                    hide-details
                    :items="['release', 'beta', 'alpha']"
                    label="Channel"
                    multiple
                  />
                </v-col>

                <v-col>
                  <v-select
                    v-model="filterPlatform"
                    chips
                    clearable
                    color="primary"
                    density="compact"
                    hide-details
                    :items="possiblePlatforms"
                    label="Platform"
                    multiple
                  />
                </v-col>
              </v-row>
            </v-container>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <Card>
          <v-card-text>
            <v-data-table
              class="bg-transparent"
              :headers="tableHeaders"
              hover
              :items="filteredVersions"
              :loading="loading"
              @click:row="onRowClick"
            >
              <template #no-data>
                <div class="pa-6 text-medium-emphasis">{{
                  loadError || 'No versions match the selected filters.'
                }}
                </div>
              </template>

              <template #item.name="{ item }">
                <VersionChip
                  :space-i-d="item.space_id"
                  :version="item.name"
                />
              </template>

              <template #item.channel="{ item }">
                <VersionChannelChip
                  :version="item"
                />
              </template>

              <template #item.platform="{ item }">
                <VersionPlatformChip
                  :version="item"
                />
              </template>

              <template #item.actions="{ item }">
                <div class="actions__width">
                  <v-btn
                    v-if="item.files.length != 1"
                    class="mr-4 my-1"
                    icon="mdi-download"
                    :to="`/spaces/${space?.slug}/versions/${item.name}`"
                    variant="text"
                  />

                  <v-btn
                    v-else
                    class="my-1"
                    :href="item.files[0]?.url"
                    icon="mdi-download"
                    variant="text"
                  />

                  <v-btn
                    v-if="isMember"
                    class="ml-4 my-1"
                    icon="mdi-pencil"
                    :to="`/spaces/${space?.slug}/versions/${item.name}/edit`"
                    variant="text"
                  />

                  <v-btn
                    v-if="isMember"
                    class="ml-4 my-1"
                    color="red"
                    icon="mdi-delete"
                    variant="text"
                    @click="deleteVersionReq($event, item)"
                  />
                </div>
              </template>

            </v-data-table>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
table, tr, td, thead, tbody {
  background: transparent;
  border-collapse: collapse;
}

.actions__width {
  min-width: 55px;
}
</style>
