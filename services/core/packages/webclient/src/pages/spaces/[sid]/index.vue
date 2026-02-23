<script lang="ts" setup>

import {mapCategoryToDisplayname, mapLinkToDisplayname, type Space} from "@/api/spaces/types.ts";
import {getDownloadCountForSpace, getSpace} from "@/api/spaces/spaces.ts";
import SpaceHeader from "@/components/SpaceHeader.vue";
import type {SpaceVersion} from "@/api/versions/types.ts";
import {getLatestVersion} from "@/api/versions/versions.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import Card from "@/components/common/Card.vue";
import {getPublicUser} from "@/api/auth/users.ts";
import type {User} from "@/api/auth/types.ts";

const route = useRoute();

const space = ref<Space>();
const creator = ref<User>();
const members = ref<User[]>([]);
const latestVersion = ref<SpaceVersion>();
const downloadCount = ref<number>(0);

onMounted(async () => {
  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  try {
    latestVersion.value = await getLatestVersion(space.value.id);
  } catch (e) {
    console.error("Error fetching latest version:", e);
  }

  downloadCount.value = await getDownloadCountForSpace(space.value.id);

  creator.value = await getPublicUser(space.value.creator);
   for (const member of space.value.members) {
     const user = await getPublicUser(member.user_id);
     members.value.push(user);
   }

  useHead({
    title: `${space.value.title} - FancySpaces`,
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
  <v-container v-if="space" width="90%">
    <v-row>
      <v-col class="flex-grow-0 pa-0">
        <SpaceSidebar
          :space="space"
        />
      </v-col>

      <v-col>
        <SpaceHeader :space="space"/>

        <hr
          class="mt-4 grey-border-color"
        />
      </v-col>
    </v-row>

    <v-row class="mobile-space-sidebar-buttons">
      <v-col>
        <v-btn
          :to="`/spaces/${space?.id}`"
          class="mr-4"
          color="secondary"
          prepend-icon="mdi-information-slab-circle-outline"
          size="large"
          variant="tonal"
        >
          Information
        </v-btn>

        <v-btn
          :to="`/spaces/${space?.id}/versions`"
          color="secondary"
          prepend-icon="mdi-file-download-outline"
          size="large"
          variant="tonal"
        >
          Versions
        </v-btn>
      </v-col>
    </v-row>

    <v-row>
      <v-col md="8">
        <Card class="mb-4">
          <v-card-text>
            <MarkdownRenderer
              :markdown="space?.description || ''"
            />
          </v-card-text>
        </Card>
      </v-col>

      <v-col md="4">
        <Card
          class="mb-4"
          min-width="300"
        >
          <v-card-title class="mt-2">Details</v-card-title>

          <v-card-text>
            <p class="text-body-1"><strong>ID:</strong> {{ space?.id }}</p>
            <p class="text-body-1"><strong>Slug:</strong> {{ space?.slug }}</p>
            <p class="text-body-1"><strong>Status:</strong> {{ space?.status }}</p>
            <p class="text-body-1"><strong>Updated at:</strong> {{ latestVersion?.published_at.toLocaleDateString() || space?.created_at.toLocaleDateString() }}</p>
            <p class="text-body-1"><strong>Created at:</strong> {{ space?.created_at.toLocaleDateString() }}</p>
          </v-card-text>
        </Card>

        <Card
          class="mb-4"
          min-width="300"
        >
          <v-card-title class="mt-2">Categories</v-card-title>

          <v-card-text>
            <div v-for="category in space?.categories" :key="category" class="d-inline-block">
              <v-chip
                class="ma-1"
                color="tertiary"
                rounded
                variant="tonal"
              >
                {{ mapCategoryToDisplayname(category) }}
              </v-chip>
            </div>
          </v-card-text>
        </Card>

        <Card
          class="mb-4"
          elevation="12"
        >
          <v-card-title class="mt-2">Links</v-card-title>

          <v-card-text>
            <div v-for="(link) in space?.links" :key="link.name" class="mb-2">
              <a
                :href="link.url"
                target="_blank"
              >
                <div class="d-flex align-center">
                  <v-icon
                    class="me-2"
                    icon="mdi-link-variant"
                  />
                  <p class="text-body-1 link--hover">{{ mapLinkToDisplayname(link.name) }}</p>
                </div>
              </a>
            </div>
          </v-card-text>
        </Card>

        <Card
          class="mb-4"
          min-width="300"
        >
          <v-card-title class="mt-2">Members</v-card-title>

          <v-card-text>
            <RouterLink :to="'/users/'+creator?.name">
              <div class="d-flex align-center">
                <v-icon
                  class="me-2"
                  icon="mdi-account-star-outline"
                />
                <p class="text-body-1 link--hover">{{ creator?.name }} (Creator)</p>
              </div>
            </RouterLink>

            <div v-for="member in members" :key="member.id">
              <RouterLink :to="'/users/'+member.name">
                <div class="d-flex align-center">
                  <v-icon
                    class="me-2"
                    icon="mdi-account-outline"
                  />
                  <p class="text-body-1 link--hover">{{ member.name }}</p>
                </div>
              </RouterLink>
            </div>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>
  <v-container v-else width="60%">
    <v-row>
      <v-col>
        <p class="text-h4 text-center mt-8">Space not found or still loading...</p>
      </v-col>
    </v-row>
  </v-container>

</template>

<style scoped>

</style>

<style>
.mobile-space-sidebar-buttons {
  display: none;
}

@media (max-width: 960px) {
  .mobile-space-sidebar-buttons {
    display: block;
    margin-top: 16px;
    margin-bottom: 16px;
  }
}
</style>
