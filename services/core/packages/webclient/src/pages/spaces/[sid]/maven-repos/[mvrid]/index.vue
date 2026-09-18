<script lang="ts" setup>

  import type { SpaceMavenRepository, SpaceMavenRepositoryArtifact } from '@/api/maven/types'
  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { deleteMavenRepository, getAllMavenArtifacts, getMavenRepository } from '@/api/maven/maven'
  import { getSpace } from '@/api/spaces/spaces'
  import Card from '@/components/common/Card.vue'
  import SpaceHeader from '@/components/SpaceHeader.vue'
  import SpaceSidebar from '@/components/SpaceSidebar.vue'
  import { useConfirmationStore } from '@/stores/confirmation'
  import { useNotificationStore } from '@/stores/notifications'
  import { useUserStore } from '@/stores/user'

  const route = useRoute()
  const router = useRouter()
  const userStore = useUserStore()
  const confirmationStore = useConfirmationStore()
  const notifications = useNotificationStore()

  const isLoggedIn = ref(false)

  const space = ref<Space>()

  const repo = ref<SpaceMavenRepository>()
  const artifacts = ref<SpaceMavenRepositoryArtifact[]>([])
  const canManage = computed(() => {
    if (!space.value || !isLoggedIn.value || !userStore.user) return false
    return space.value.creator === userStore.user.id || space.value.members.some(member => member.user_id === userStore.user?.id && ['member', 'admin'].includes(member.role))
  })

  const howToUseTab = ref('build.gradle.kts')

  function latestVersion (artifact: SpaceMavenRepositoryArtifact) {
    return artifact.versions.reduce<SpaceMavenRepositoryArtifact['versions'][number] | undefined>((latest, version) => {
      if (!latest || new Date(version.published_at).getTime() > new Date(latest.published_at).getTime()) return version
      return latest
    }, undefined)
  }

  const tableHeaders = [
    { title: 'Group ID', value: 'group' },
    { title: 'Artifact ID', value: 'id' },
    { title: 'Versions', key: 'versions', value: (art: SpaceMavenRepositoryArtifact) => art.versions.length || 'N/A' },
    {
      title: 'Latest version', key: 'latest-version', value: (art: SpaceMavenRepositoryArtifact) => {
        const latest = latestVersion(art)
        return latest ? latest.version : 'N/A'
      },
    },
    {
      title: 'Last update', key: 'last-update', value: (art: SpaceMavenRepositoryArtifact) => {
        const latest = latestVersion(art)
        return latest ? new Date(latest.published_at).toLocaleDateString() : 'N/A'
      },
    },
  ]

  onMounted(async () => {
    isLoggedIn.value = await userStore.isAuthenticated

    const spaceID = (route.params as any).sid as string
    space.value = await getSpace(spaceID)

    if (!space.value.maven_repository_settings.enabled) {
      router.push(`/spaces/${space.value.slug}`)
      return
    }

    const mavenRepoName = (route.params as any).mvrid as string
    repo.value = await getMavenRepository(space.value.id, mavenRepoName)

    artifacts.value = await getAllMavenArtifacts(space.value.id, repo.value.name)

    useHead({
      title: `${space.value.title} Maven Repo ${repo.value.name} - FancySpaces`,
      meta: [
        {
          name: 'description',
          content: space.value.summary || `Explore the ${space.value.title} project space on FancySpaces.`,
        },
      ],
    })
  })

  function onRowClick (event: any, { item }: any) {
    router.push(`/spaces/${space.value?.slug}/maven-repos/${repo.value?.name}/${item.group}:${item.id}`)
  }

  function confirmDelete () {
    if (!space.value || !repo.value) return
    confirmationStore.confirmation = {
      shown: true,
      persistent: true,
      title: 'Delete Maven repository',
      text: `Delete ${repo.value.name}? All artifacts in this repository will no longer be available. This action cannot be undone.`,
      yesText: 'Delete',
      onConfirm: async () => {
        try {
          await deleteMavenRepository(space.value!.id, repo.value!.name)
          notifications.info('Repository deleted successfully')
          await router.push(`/spaces/${space.value!.slug}/maven-repos`)
        } catch (error_) {
          notifications.error(error_ instanceof Error ? error_.message : 'Unable to delete repository.')
        }
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
            <p class="text-body-2">{{ artifacts?.length }} artifacts</p>
          </template>

          <template #quick-actions>
            <v-btn
              v-if="canManage"
              color="primary"
              size="large"
              :to="`/spaces/${space?.slug}/maven-repos/new`"
              variant="tonal"
            >
              New Repo
            </v-btn>
          </template>
        </SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        >
      </v-col>
    </v-row>

    <v-row>
      <v-col class="mb-4">
        <Card
          class="bg-transparent"
          color="#150D1950"
          min-width="600"
        >
          <v-card-text>
            <v-breadcrumbs
              class="pa-0"
              color="primary"
              :items="[
                { title: 'Maven Repositories', to: `/spaces/${space?.slug}/maven-repos` },
                { title: repo?.name || '', to: `/spaces/${space?.slug}/maven-repos/${repo?.name}` },
              ]"
            />
          </v-card-text>
        </Card>
      </v-col>
    </v-row>

    <v-row>
      <v-col md="7">
        <Card
          class="bg-transparent"
          color="#150D1950"
          min-width="600"
        >
          <v-card-title class="mt-2">
            Artifacts in {{ repo?.name }}
          </v-card-title>

          <v-card-text>
            <v-data-table
              class="bg-transparent"
              :headers="tableHeaders"
              hover
              :items="artifacts"
              @click:row="onRowClick"
            />
          </v-card-text>
        </Card>
      </v-col>

      <v-col md="5">
        <Card
          class="bg-transparent mb-4"
          elevation="6"
        >
          <v-card-title class="mt-2">How to use</v-card-title>

          <v-card-text>
            <v-tabs
              v-model="howToUseTab"
              background-color="#150D1950"
              color="primary"
              grow
            >
              <v-tab value="build.gradle.kts">build.gradle.kts</v-tab>
              <v-tab value="build.gradle">build.gradle</v-tab>
              <v-tab value="pom.xml">pom.xml</v-tab>
            </v-tabs>

            <v-tabs-window v-model="howToUseTab" class="mt-4">
              <v-tabs-window-item value="build.gradle.kts">
                <pre><code>repositories {
    maven (url = "https://maven.fancyspaces.net/{{ space?.slug }}/{{ repo?.name }}")
}</code></pre>
              </v-tabs-window-item>

              <v-tabs-window-item value="build.gradle">
                <pre><code>repositories {
    maven {
        url "https://maven.fancyspaces.net/{{ space?.slug }}/{{ repo?.name }}"
    }
}</code></pre>
              </v-tabs-window-item>

              <v-tabs-window-item value="pom.xml">
                <pre><code>&lt;repositories&gt;
    &lt;repository&gt;
        &lt;id&gt;fancyspaces-{{ space?.slug }}-{{ repo?.name }}&lt;/id&gt;
        &lt;url&gt;https://maven.fancyspaces.net/{{ space?.slug }}/{{ repo?.name }}&lt;/url&gt;
    &lt;/repository&gt;
&lt;/repositories&gt;</code></pre>
              </v-tabs-window-item>
            </v-tabs-window>
          </v-card-text>
        </Card>

        <Card
          v-if="canManage"
          class="bg-transparent"
          elevation="6"
        >
          <v-card-text>
            <v-btn
              block
              class="mb-2"
              color="primary"
              :to="`/spaces/${space?.slug}/maven-repos/${encodeURIComponent(repo?.name || '')}/edit`"
              variant="tonal"
            >
              Edit Repo
            </v-btn>

            <v-btn
              block
              color="error"
              variant="tonal"
              @click="confirmDelete"
            >
              Delete Repo
            </v-btn>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
pre {
  overflow-x: auto;
}
</style>
