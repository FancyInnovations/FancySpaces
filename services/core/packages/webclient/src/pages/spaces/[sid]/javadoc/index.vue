<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import type {
  SpaceMavenRepository,
  SpaceMavenRepositoryArtifact,
  SpaceMavenRepositoryArtifactVersion
} from "@/api/maven/types.ts";
import {getAllMavenArtifacts, getAllMavenRepositories} from "@/api/maven/maven.ts";
import SpaceHeader from "@/components/SpaceHeader.vue";

const router = useRouter();

const space = ref<Space>();
const repos = ref<SpaceMavenRepository[]>([]);
const artifacts = ref<SpaceMavenRepositoryArtifact[]>([]);
const versions = computed(() => {
  return selectedArtifact.value?.versions.sort((a, b) => {
    return b.published_at.getTime() - a.published_at.getTime();
  }) || [];
});

const selectedRepo = ref<SpaceMavenRepository>();
const selectedArtifact = ref<SpaceMavenRepositoryArtifact>();
const selectedVersion = ref<SpaceMavenRepositoryArtifactVersion>();

const javadocURL = computed(() => {
  if (!space.value || !selectedRepo.value || !selectedArtifact.value || !selectedVersion.value) {
    return '';
  }

  const baseURL = window.location.origin;
  return `${baseURL}/javadoc/${space.value?.id}/${selectedRepo.value.name}/${selectedArtifact.value.group + ':' + selectedArtifact.value.id}/${selectedVersion.value.version}/index.html`;
});

watch(selectedRepo, async (newRepo) => {
  if (newRepo) {
    artifacts.value = await getAllMavenArtifacts(space.value!.id, newRepo.name);
    selectedArtifact.value = artifacts.value[0];
  } else {
    artifacts.value = [];
    selectedArtifact.value = undefined;
  }
});

watch(selectedArtifact, (newArtifact) => {
  if (newArtifact) {
    selectedVersion.value = versions.value[0];
  } else {
    selectedVersion.value = undefined;
  }
});

onMounted(async () => {
  const spaceID = (useRoute().params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.maven_repository_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  repos.value = await getAllMavenRepositories(space.value.id);

  useHead({
    title: `${space.value.title} javadoc - FancySpaces`,
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
        <SpaceSidebar
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
        <v-card
          class="card__border"
          color="#19120D33"
          elevation="12"
          rounded="xl"
        >
          <v-card-title class="mt-2">Select artifact</v-card-title>

          <v-card-text class="d-flex flex-wrap align-center">
            <v-select
              v-model="selectedRepo"
              :item-value="(item) => item"
              :items="repos"
              color="primary"
              density="compact"
              hide-details
              item-title="name"
              label="Select repository"
            />

            <v-select
              v-model="selectedArtifact"
              :item-title="(item) => item.group + ':' + item.id"
              :item-value="(item) => item"
              :items="artifacts"
              class="ml-4"
              color="primary"
              density="compact"
              hide-details
              label="Select artifact"
            />

            <v-select
              v-model="selectedVersion"
              :item-value="(item) => item"
              :items="versions"
              class="ml-4"
              color="primary"
              density="compact"
              hide-details
              item-title="version"
              label="Select version"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <v-card
          class="card__border"
          color="#19120D33"
          elevation="12"
          rounded="xl"
        >
          <v-card-title class="d-flex justify-space-between align-center">
            <span>JavaDoc</span>

            <v-btn
              v-if="javadocURL"
              :href="javadocURL"
              icon
              target="_blank"
              variant="text"
            >
              <v-icon>mdi-open-in-new</v-icon>
            </v-btn>
          </v-card-title>

          <v-card-text>
            <iframe
              v-if="javadocURL"
              :src="javadocURL"
              height="750px"
              width="100%"
            />
            <span v-else>
              Please select a repository, artifact, and version to view the javadoc.
            </span>
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
