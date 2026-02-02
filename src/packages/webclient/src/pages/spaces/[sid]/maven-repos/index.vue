<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import type {SpaceMavenRepository} from "@/api/maven/types.ts";
import {getAllMavenRepositories} from "@/api/maven/maven.ts";

const router = useRouter();

const isLoggedIn = computed(() => {
  return localStorage.getItem("fs_api_key") !== null;
});

const space = ref<Space>();
const repos = ref<SpaceMavenRepository[]>();

onMounted(async () => {
  const spaceID = (useRoute().params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.maven_repository_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  repos.value = await getAllMavenRepositories(space.value.id);

  useHead({
    title: `${space.value.title} maven repositories - FancySpaces`,
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
      <v-col>
        <v-card
          color="error"
          rounded="xl"
          variant="tonal"
        >
          <v-card-text>
            <v-icon
              class="mr-2"
            >
              mdi-alert-circle
            </v-icon>
            The maven repository feature is currently in beta and may not function as expected. We appreciate your patience as we work to improve it!
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col class="flex-grow-0 pa-0">
        <SpaceSidebar
          :space="space"
        />
      </v-col>

      <v-col>
        <div class="d-flex justify-space-between">
          <div class="d-flex flex-column justify-center">
            <v-img
              :href="`/spaces/${space?.slug}`"
              :src="space?.icon_url || '/logo.png'"
              alt="Space Icon"
              height="100"
              max-height="100"
              max-width="100"
              min-height="100"
              min-width="100"
              width="100"
            />
          </div>

          <div class="mx-4 d-flex flex-column justify-space-between flex-grow-1">
            <div>
              <h1>{{ space?.title }}</h1>
              <p class="text-body-1 mt-2">{{ space?.summary }}</p>
            </div>
          </div>

          <div class="d-flex flex-column justify-center">
            <v-btn
              v-if="isLoggedIn"
              :to="`/spaces/${space?.slug}/maven-repos/new`"
              color="primary"
              disabled
              size="large"
              variant="tonal"
            >
              New Repo
            </v-btn>
          </div>
        </div>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col
        v-for="repo in repos"
        :key="repo.name"
        md="3"
      >
        <v-card
          class="card__border"
          color="#19120D33"
          elevation="12"
          rounded="xl"
        >
          <v-card-title class="mt-2">
            Repository: {{ repo.name }}
          </v-card-title>

          <v-card-text>
            <p><strong>Public:</strong> {{ repo.public ? 'Yes' : 'No' }}</p>
            <p><strong>Created at:</strong> {{ repo.created_at.toLocaleString() }}</p>
          </v-card-text>

          <v-card-actions>
            <v-btn
              :to="`/spaces/${space?.slug}/maven-repos/${repo.name}`"
              color="primary"
              variant="text"
            >
              View Repository
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.grey-border-color {
  border-color: rgba(0, 0, 0, 0.8);
}

</style>
