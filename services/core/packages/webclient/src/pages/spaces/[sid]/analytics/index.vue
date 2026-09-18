<script lang="ts" setup>

  import type { Dashboard } from '@/api/analytics/dashboards/types'
  import type { Space } from '@/api/spaces/types'
  import { useHead } from '@vueuse/head'
  import { getDashboards } from '@/api/analytics/dashboards/dashboards'
  import { getSpace } from '@/api/spaces/spaces'
  import AnalyticsSidebar from '@/components/analytics/AnalyticsSidebar.vue'
  import Card from '@/components/common/Card.vue'
  import SpaceHeader from '@/components/SpaceHeader.vue'
  import { useUserStore } from '@/stores/user'

  const router = useRouter()
  const route = useRoute()
  const userStore = useUserStore()

  const isLoggedIn = ref(false)
  const isMember = ref(false)

  const space = ref<Space>()
  const dashboards = ref<Dashboard[]>()

  onMounted(async () => {
    isLoggedIn.value = await userStore.isAuthenticated

    const spaceID = (route.params as any).sid as string
    space.value = await getSpace(spaceID)

    if (!space.value.analytics_settings.enabled) {
      router.push(`/spaces/${space.value.slug}`)
      return
    }

    isMember.value = (await userStore.isAuthenticated) && (space.value.creator == userStore.user?.id || space.value.members.some(member => member.user_id === userStore.user?.id))
    if (!isMember.value) {
      router.push(`/spaces/${space.value.slug}/analytics/dashboards`)
      return
    }

    dashboards.value = await getDashboards(space.value.id)

    useHead({
      title: `${space.value.title} analytics portal - FancySpaces`,
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
        <AnalyticsSidebar
          :dashboards="dashboards"
          :space="space"
        />
      </v-col>

      <v-col>
        <SpaceHeader :space="space" />

        <hr
          class="grey-border-color mt-4"
        >
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <p class="text-h4 text-center">
          Explore data of {{ space?.title }}
        </p>
      </v-col>
    </v-row>

    <v-row class="mb-4" justify="center">
      <v-col md="3">
        <Card
          prepend-icon="mdi-chart-timeline-variant"
          text="In the metrics explorer, you can preview all collected metric data for your project."
          title="Metrics"
          :to="`/spaces/${space?.slug}/analytics/metrics`"
        />
      </v-col>

      <v-col md="3">
        <Card
          prepend-icon="mdi-radar"
          text="In the event explorer, you can preview all collected event data for your project."
          title="Events"
          :to="`/spaces/${space?.slug}/analytics/events`"
        />
      </v-col>

      <v-col md="3">
        <Card
          prepend-icon="mdi-script-text"
          text="In the logs explorer, you can preview all collected log data for your project."
          title="Logs"
          :to="`/spaces/${space?.slug}/analytics/logs`"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="3">
        <v-badge
          color="error"
          content="SOON"
          offset-x="20"
        >
          <Card
            disabled
            prepend-icon="mdi-bug"
            text="In the exception overview, you can view and manage all exceptions for your project."
            title="Exceptions"
            :to="`/spaces/${space?.slug}/analytics/exceptions`"
          />
        </v-badge>
      </v-col>

      <v-col md="3">
        <v-badge
          color="error"
          content="SOON"
          offset-x="20"
        >
          <Card
            disabled
            prepend-icon="mdi-bell-ring"
            text="In the alert overview, you can view and manage all alerts for your project."
            title="Alerts"
            :to="`/spaces/${space?.slug}/analytics/alerts`"
          />
        </v-badge>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
