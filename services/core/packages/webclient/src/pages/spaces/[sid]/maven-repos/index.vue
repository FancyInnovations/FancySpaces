<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import type {SpaceMavenRepository} from "@/api/maven/types.ts";
import {getAllMavenRepositories} from "@/api/maven/maven.ts";
import SpaceHeader from "@/components/SpaceHeader.vue";
import {useUserStore} from "@/stores/user.ts";
import Card from "@/components/common/Card.vue";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const isLoggedIn = ref(false);

const space = ref<Space>();
const repos = ref<SpaceMavenRepository[]>();

onMounted(async () => {
  isLoggedIn.value = await userStore.isAuthenticated;

  const spaceID = (route.params as any).sid as string;
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
      <v-col class="flex-grow-0 pa-0">
        <SpaceSidebar
          :space="space"
        />
      </v-col>

      <v-col>
        <SpaceHeader :space="space">
          <template #metadata>
            <p class="text-body-2 mx-4">-</p>
            <p class="text-body-2">{{ repos?.length }} repositories</p>
          </template>

          <template #quick-actions>
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
          </template>
        </SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col
        v-for="repo in repos"
        v-if="repos && repos.length > 0"
        :key="repo.name"
        md="3"
      >
        <Card>
          <v-card-title class="mt-2">
            Repository: {{ repo.name }}
          </v-card-title>

          <v-card-text>
            <p><strong>Public:</strong> {{ repo.public ? 'Yes' : 'No' }}</p>
            <p><strong>Created at:</strong> {{ repo.created_at.toLocaleString() }}</p>
            <template v-if="repo.internal_mirror">
              <p><strong>Internal Mirror:</strong> {{ repo.internal_mirror?.space_id }} / {{ repo.internal_mirror?.repository }}</p>
            </template>
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
        </Card>
      </v-col>
      <v-col
        v-else
        class="text-center"
        cols="12"
      >
        <h2>No repositories found.</h2>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.grey-border-color {
  border-color: rgba(0, 0, 0, 0.8);
}

</style>
