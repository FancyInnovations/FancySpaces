<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import type {SpaceMavenRepository, SpaceMavenRepositoryArtifact} from "@/api/maven/types.ts";
import {getAllMavenArtifacts, getMavenRepository} from "@/api/maven/maven.ts";

const route = useRoute();
const router = useRouter();

const isLoggedIn = computed(() => {
  return localStorage.getItem("fs_api_key") !== null;
});

const space = ref<Space>();

const repo = ref<SpaceMavenRepository>();
const artifacts = ref<SpaceMavenRepositoryArtifact[]>([]);

const latestVersion = computed(() => {
  return artifacts.value
    .flatMap(art => art.versions)
    .sort((a, b) => {
      return new Date(b.published_at).getTime() - new Date(a.published_at).getTime();
    })[0];
});

const tableHeaders = [
  { title: 'Group ID', value: 'group' },
  { title: 'Artifact ID', value: 'id' },
  { title: 'Versions', key: 'versions', value: (art: SpaceMavenRepositoryArtifact) => art.versions.length || 'N/A' },
  { title: 'Latest version', key: 'latest-version', value: (art: SpaceMavenRepositoryArtifact) => latestVersion.value?.version || 'N/A' },
];

onMounted(async () => {
  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.maven_repository_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  const mavenRepoName = (route.params as any).mvrid as string;
  repo.value = await getMavenRepository(space.value.id, mavenRepoName);

  artifacts.value = await getAllMavenArtifacts(space.value.id, repo.value.name);

  useHead({
    title: `${space.value.title} Maven Repo ${repo.value.name} - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || `Explore the ${space.value.title} project space on FancySpaces.`
      }
    ]
  });
});

function onRowClick(event: any, { item }: any) {
  router.push(`/spaces/${space.value?.slug}/maven-repos/${repo.value?.name}/${item.group}:${item.id}`);
}

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
          class="grey-border-color"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col md="9">
        <v-card
          class="card__border bg-transparent"
          color="#150D1950"
          elevation="12"
          min-width="600"
          rounded="xl"
        >
          <v-card-title class="mt-2">
            Artifacts in {{ repo?.name }} ({{ artifacts.length }})
          </v-card-title>

          <v-card-subtitle>
            Created at: {{ repo?.created_at.toLocaleDateString() }}
          </v-card-subtitle>

          <v-card-text>
            <v-data-table
              :headers="tableHeaders"
              :items="artifacts"
              class="bg-transparent"
              hover
              @click:row="onRowClick"
            >
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col>
        <v-card
          v-if="isLoggedIn"
          class="card__border bg-transparent"
          color="#19120D33"
          elevation="6"
          rounded="xl"
        >
          <v-card-text>
            <v-btn
              :to="`/spaces/${space?.slug}/maven-repo/${repo?.name}/edit`"
              block
              class="mb-2"
              color="primary"
              variant="tonal"
            >
              Edit Repo
            </v-btn>

            <v-btn
              :to="`/spaces/${space?.slug}/maven-repo`"
              block
              color="error"
              variant="tonal"
            >
              Delete Repo
            </v-btn>
          </v-card-text>
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
