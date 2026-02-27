<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import SpaceHeader from "@/components/SpaceHeader.vue";
import {useUserStore} from "@/stores/user.ts";
import AnalyticsSidebar from "@/components/analytics/AnalyticsSidebar.vue";
import type {Dashboard} from "@/api/analytics/dashboards/types.ts";
import {createDashboard, getDashboards} from "@/api/analytics/dashboards/dashboards.ts";
import {useNotificationStore} from "@/stores/notifications.ts";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
const notifications = useNotificationStore();

const isLoggedIn = ref(false);

const space = ref<Space>();
const dashboards = ref<Dashboard[]>();

const name = ref('');
const summary = ref('');
const isPublic = ref(false);

const nameRule = (value: string) => {
  if (!value) return 'Name is required';
  return true;
};

const isEverythingValid = computed(() => {
  return nameRule(name.value) === true
});

async function create() {
  if (!space.value) return;

  await createDashboard(space.value?.id, name.value, summary.value, isPublic.value, []);
  notifications.info("Created new dashboard!");
  await router.push(`/spaces/${space.value?.slug}/analytics/dashboards`);
}

onMounted(async () => {
  isLoggedIn.value = await userStore.isAuthenticated;

  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.analytics_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

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

    <v-row justify="center">
      <v-col md="6">
        <Card>
          <v-card-title class="mt-2">
            Create a new dashboard
          </v-card-title>

          <v-card-text>
            <v-text-field
              v-model="name"
              :rules="[nameRule]"
              autofocus
              class="mb-4 mt-4"
              color="primary"
              label="Name *"
            />

            <v-textarea
              v-model="summary"
              class="mb-4"
              color="primary"
              hide-details
              label="Summary"
              no-resize
              rows="4"
            />

            <v-checkbox
              v-model="isPublic"
              class="mb-4"
              color="primary"
              hide-details
              label="Make dashboard public"
            />

            <v-btn
              :disabled="!isEverythingValid"
              class="mt-4"
              color="primary"
              @click="create"
            >
              Create
            </v-btn>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
