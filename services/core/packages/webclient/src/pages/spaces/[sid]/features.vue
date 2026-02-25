<script lang="ts" setup>

import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import {useUserStore} from "@/stores/user.ts";
import {type Space} from "@/api/spaces/types.ts";
import SpaceHeader from "@/components/SpaceHeader.vue";
import SpaceSidebar from "@/components/SpaceSidebar.vue";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const isLoggedIn = ref(false);

const space = ref<Space>();
const isReleasesEnabled = computed(() => space.value?.release_settings.enabled ?? false);
const isIssuesEnabled = computed(() => space.value?.issue_settings.enabled ?? false);
const isMavenRepositoryEnabled = computed(() => space.value?.maven_repository_settings.enabled ?? false);
const isAnalyticsEnabled = computed(() => space.value?.analytics_settings.enabled ?? false);
const isStorageEnabled = computed(() => space.value?.storage_settings.enabled ?? false);
const isSecretsEnabled = computed(() => space.value?.secrets_settings.enabled ?? false);

onMounted(async () => {
  isLoggedIn.value = await userStore.isAuthenticated;

  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!isLoggedIn) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  useHead({
    title: `${space.value.title} features - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: 'Manage the features of this space on FancySpaces.'
      }
    ]
  });
});

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
        <SpaceHeader :space="space"/>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col>
        <div class="d-flex flex-column justify-center">
          <h1 class="text-center mb-2">Space features</h1>
          <p class="text-center text-body-1">Here you can see which features are enabled for this space.</p>
          <p class="text-center text-body-1 mb-8">If you want to enable or disable features, please contact the support.</p>

          <v-checkbox
            v-model="isReleasesEnabled"
            hide-details
            label="Releases and Downloads"
            readonly
            width="fit-content"
          />

          <v-checkbox
            v-model="isIssuesEnabled"
            hide-details
            label="Issues and Bug Tracking"
            readonly
          />

          <v-checkbox
            v-model="isMavenRepositoryEnabled"
            hide-details
            label="Maven Repository and Javadoc"
            readonly
          />

          <v-checkbox
            v-model="isAnalyticsEnabled"
            hide-details
            label="Analytics"
            readonly
          />

          <v-checkbox
            v-model="isStorageEnabled"
            hide-details
            label="Managed Storage"
            readonly
          />

          <v-checkbox
            v-model="isSecretsEnabled"
            hide-details
            label="Secrets Management"
            readonly
          />
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
