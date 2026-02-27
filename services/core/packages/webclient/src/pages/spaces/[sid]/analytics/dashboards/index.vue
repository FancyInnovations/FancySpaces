<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import SpaceHeader from "@/components/SpaceHeader.vue";
import {useUserStore} from "@/stores/user.ts";
import AnalyticsSidebar from "@/components/analytics/AnalyticsSidebar.vue";
import type {Dashboard} from "@/api/analytics/dashboards/types.ts";
import {getDashboards} from "@/api/analytics/dashboards/dashboards.ts";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const isLoggedIn = ref(false);
const isMember = ref(false);

const space = ref<Space>();
const dashboards = ref<Dashboard[]>();

onMounted(async () => {
  isLoggedIn.value = await userStore.isAuthenticated;

  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.analytics_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  isMember.value = (await userStore.isAuthenticated) && (space.value.creator == userStore.user?.id || space.value.members.some(member => member.user_id === userStore.user?.id));

  dashboards.value = await getDashboards(space.value.id);

  useHead({
    title: `${space.value.title} analytics portal - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || `Explore the ${space.value.title} project space on FancySpaces.`
      }
    ]
  });
});
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
        <SpaceHeader :space="space"></SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <p class="text-h4 text-center">
          Dashboard overview
        </p>
      </v-col>
    </v-row>

    <v-row align="stretch" class="mb-4" justify="center">
      <v-col
        v-for="d in dashboards"
        :key="d.dashboard_id"
        md="3"
      >
        <Card
          :appendIcon="d.public ? 'mdi-lock-open-variant-outline' : 'mdi-lock-outline'"
          :subtitle="'Created at ' + new Date(d.created_at).toLocaleDateString()"
          :title="d.name"
          :to="`/spaces/${space?.id}/analytics/dashboards/${d.dashboard_id}`"
          height="100%"
        >
          <v-card-text>{{ d.summary }}</v-card-text>
        </Card>
      </v-col>

      <v-col v-if="isMember" md="3">
        <Card
          :to="`/spaces/${space?.id}/analytics/dashboards/new`"
          class="d-flex align-center justify-center"
          height="100%"
          min-height="120px"
        >
          <v-icon size="48">mdi-plus</v-icon>
        </Card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
