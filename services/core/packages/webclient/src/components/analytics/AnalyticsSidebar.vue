<script lang="ts" setup>

import {type Space} from "@/api/spaces/types.ts";
import {useUserStore} from "@/stores/user.ts";
import type {Dashboard} from "@/api/analytics/dashboards/types.ts";

const userStore = useUserStore();

const props = defineProps<{
  space?: Space;
  dashboards?: Dashboard[];
}>();

const isMember = computed(() => {
  if (!props.space) return false;
  if (!userStore.isAuthenticated) return false;

  const userID =  userStore.user?.id;
  return props.space.creator == userID || props.space.members.some(member => member.user_id === userID);
});

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

      <v-divider class="mt-2"/>

      <v-list-item
        v-if="isMember"
        :to="`/spaces/${space?.slug}/analytics/metrics`"
        link
        prepend-icon="mdi-chart-timeline-variant"
        title="Metrics"
      />
      <v-list-item
        v-if="isMember"
        :to="`/spaces/${space?.slug}/analytics/events`"
        link
        prepend-icon="mdi-radar"
        title="Events"
      />
      <v-list-item
        v-if="isMember"
        :to="`/spaces/${space?.slug}/analytics/logs`"
        link
        prepend-icon="mdi-script-text"
        title="Logs"
      />

      <v-list-item
        v-if="isMember"
        :to="`/spaces/${space?.slug}/analytics/exceptions`"
        disabled
        link
        prepend-icon="mdi-bug"
        title="Exceptions"
      >
        <template v-slot:append>
          <v-badge
            color="error"
            content="SOON"
            inline
          />
        </template>
      </v-list-item>

      <v-list-item
        v-if="isMember"
        :to="`/spaces/${space?.slug}/analytics/alerts`"
        disabled
        link
        prepend-icon="mdi-bell-ring"
        title="Alerts"
      >
        <template v-slot:append>
          <v-badge
            color="error"
            content="SOON"
            inline
          />
        </template>
      </v-list-item>

      <v-divider class="mx-2"/>
      <v-list-subheader>Dashboards</v-list-subheader>

      <v-list-item
        :to="`/spaces/${space?.slug}/analytics/dashboards/`"
        class="mb-4"
        exact
        link
        prepend-icon="mdi-view-list"
        title="Dashboard Overview"
      />

      <template
        v-for="dashboard in dashboards"
        v-if="dashboards"
        :key="dashboard.dashboard_id"
      >
        <v-list-item
          v-if="dashboard.public || isMember"
          :title="dashboard.name"
          :to="`/spaces/${space?.slug}/analytics/dashboards/${dashboard.dashboard_id}`"
          link
          prepend-icon="mdi-view-dashboard-variant"
        />

      </template>

      <v-divider class="mx-2"/>
      <v-list-subheader>Space</v-list-subheader>

      <v-list-item
        v-if="isMember"
        :to="`/spaces/${space?.slug}`"
        exact
        link
        prepend-icon="mdi-home"
        title="Back to Space"
      />

      <v-list-item
        v-if="isMember"
        :to="`/spaces/${space?.slug}/settings`"
        link
        prepend-icon="mdi-cog-outline"
        title="Settings"
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
