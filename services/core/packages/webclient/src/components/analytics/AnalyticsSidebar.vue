<script lang="ts" setup>

  import type { Dashboard } from '@/api/analytics/dashboards/types'
  import type { Space } from '@/api/spaces/types'
  import { useUserStore } from '@/stores/user'

  const userStore = useUserStore()

  const props = defineProps<{
    space?: Space
    dashboards?: Dashboard[]
  }>()

  const isMember = computed(() => {
    if (!props.space) return false
    if (!userStore.isAuthenticated) return false

    const userID = userStore.user?.id
    return props.space.creator == userID || props.space.members.some(member => member.user_id === userID)
  })

</script>

<template>
  <v-navigation-drawer
    class="sidebar__mobile sidebar__background ma-4"
    elevation="12"
    rounded="xl"
  >
    <v-list>
      <v-list-item>
        <v-list-item-title class="text-h6 font-weight-bold">{{ space?.title }}</v-list-item-title>
        <v-list-item-subtitle>Analytics Portal</v-list-item-subtitle>
      </v-list-item>

      <v-divider class="mt-2" />

      <v-list-item
        v-if="isMember"
        link
        prepend-icon="mdi-chart-timeline-variant"
        title="Metrics"
        :to="`/spaces/${space?.slug}/analytics/metrics`"
      />

      <v-list-item
        v-if="isMember"
        link
        prepend-icon="mdi-radar"
        title="Events"
        :to="`/spaces/${space?.slug}/analytics/events`"
      />

      <v-list-item
        v-if="isMember"
        link
        prepend-icon="mdi-script-text"
        title="Logs"
        :to="`/spaces/${space?.slug}/analytics/logs`"
      />

      <v-list-item
        v-if="isMember"
        disabled
        link
        prepend-icon="mdi-bug"
        title="Exceptions"
        :to="`/spaces/${space?.slug}/analytics/exceptions`"
      >
        <template #append>
          <v-badge
            color="error"
            content="SOON"
            inline
          />
        </template>
      </v-list-item>

      <v-list-item
        v-if="isMember"
        disabled
        link
        prepend-icon="mdi-bell-ring"
        title="Alerts"
        :to="`/spaces/${space?.slug}/analytics/alerts`"
      >
        <template #append>
          <v-badge
            color="error"
            content="SOON"
            inline
          />
        </template>
      </v-list-item>

      <v-divider class="mx-2" />
      <v-list-subheader>Dashboards</v-list-subheader>

      <v-list-item
        class="mb-4"
        exact
        link
        prepend-icon="mdi-view-list"
        title="Dashboard Overview"
        :to="`/spaces/${space?.slug}/analytics/dashboards/`"
      />

      <template
        v-for="dashboard in dashboards"
        v-if="dashboards"
        :key="dashboard.dashboard_id"
      >
        <v-list-item
          v-if="dashboard.public || isMember"
          link
          prepend-icon="mdi-view-dashboard-variant"
          :title="dashboard.name"
          :to="`/spaces/${space?.slug}/analytics/dashboards/${dashboard.dashboard_id}`"
        />

      </template>

      <v-divider class="mx-2" />
      <v-list-subheader>Space</v-list-subheader>

      <v-list-item
        exact
        link
        prepend-icon="mdi-home"
        title="Back to Space"
        :to="`/spaces/${space?.slug}`"
      />

      <v-list-item
        v-if="isMember"
        link
        prepend-icon="mdi-cog-outline"
        title="Settings"
        :to="`/spaces/${space?.slug}/settings`"
      />

    </v-list>
  </v-navigation-drawer>
</template>

<style scoped>
.sidebar__background {
  max-height: calc(100vh - 96px);
  background-color: #19120D33 !important;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

@media (max-width: 960px) {
  .sidebar__mobile {
    display: none;
  }
}
</style>
