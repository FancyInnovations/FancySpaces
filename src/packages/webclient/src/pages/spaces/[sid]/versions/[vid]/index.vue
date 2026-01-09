<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getDownloadCountForSpace, getSpace} from "@/api/spaces/spaces.ts";
import {mapPlatformToDisplayname, type SpaceVersion} from "@/api/versions/types.ts";
import {getLatestVersion, getVersion} from "@/api/versions/versions.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import SpaceHeader from "@/components/SpaceHeader.vue";
import {useHead} from "@vueuse/head";

const route = useRoute();

const space = ref<Space>();
const spaceDownloadCount = ref<number>(0);
const latestVersion = ref<SpaceVersion>();

const currentVersion = ref<SpaceVersion>();

onMounted(async () => {
  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  const versionID = (route.params as any).vid as string;
  currentVersion.value = await getVersion(space.value.id, versionID);

  spaceDownloadCount.value = await getDownloadCountForSpace(space.value.id);
  latestVersion.value = await getLatestVersion(space.value.id);

  useHead({
    title: `${space.value.title} ${currentVersion.value.name} - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || `Explore the ${space.value.title} project space on FancySpaces.`
      }
    ]
  });
});

function formatSize(sizeInBytes: number): string {
  if (sizeInBytes < 1024) {
    return `${sizeInBytes} B`;
  } else if (sizeInBytes < 1024 * 1024) {
    return `${(sizeInBytes / 1024).toFixed(2)} KB`;
  } else if (sizeInBytes < 1024 * 1024 * 1024) {
    return `${(sizeInBytes / (1024 * 1024)).toFixed(2)} MB`;
  } else {
    return `${(sizeInBytes / (1024 * 1024 * 1024)).toFixed(2)} GB`;
  }
}

</script>

<template>
  <AppHeader/>

  <v-container width="90%">
    <v-row>
      <v-col class="flex-grow-0 pa-0">
        <SpaceSidebar
          :space="space"
        />
      </v-col>

      <v-col>
        <SpaceHeader
          :download-count="spaceDownloadCount"
          :latest-version="latestVersion"
          :space="space"
        />

        <hr
          class="mt-4 grey-border-color"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col md="8">
        <v-card
          class="card__border"
          color="#29152550"
          elevation="12"
          rounded="xl"
        >
          <v-card-title class="mt-2">Changelog</v-card-title>

          <v-card-text>
            <MarkdownRenderer
              :markdown="currentVersion?.changelog"
            />
          </v-card-text>
        </v-card>

        <v-card
          class="mt-4 card__border"
          color="#29152550"
          elevation="12"
          rounded="xl"
        >
          <v-card-title class="mt-2">Files</v-card-title>

          <v-card-text>
            <v-table class="bg-transparent">
              <thead>
                <tr>
                  <th class="text-left">Name</th>
                  <th class="text-left">Size</th>
                  <th class="text-right"></th>
                </tr>
              </thead>

              <tbody>
                <tr
                  v-for="file in currentVersion?.files"
                  :key="file.name"
                >
                  <td>{{ file.name }}</td>
                  <td>{{ formatSize(file.size) }}</td>
                  <td class="text-right">
                    <v-btn
                      :href="file.url"
                      color="primary"
                      icon="mdi-download"
                      small
                      target="_blank"
                      variant="text"
                    />
                  </td>
                </tr>
              </tbody>
            </v-table>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col>
        <v-card
          class="mb-4 card__border"
          color="#29152550"
          elevation="12"
          rounded="xl"
        >
          <v-card-title class="mt-2">Details</v-card-title>

          <v-card-text>
            <p class="text-body-1"><strong>Version:</strong> {{ currentVersion?.name }}</p>
            <p class="text-body-1"><strong>ID:</strong> {{ currentVersion?.id }}</p>
            <p class="text-body-1"><strong>Channel:</strong> {{ currentVersion?.channel.toUpperCase() }}</p>
            <p class="text-body-1"><strong>Platform:</strong> {{ mapPlatformToDisplayname(currentVersion?.platform) }}</p>
            <p class="text-body-1"><strong>Platform versions:</strong> {{ currentVersion?.supported_platform_versions.join(", ") }}</p>
            <p class="text-body-1"><strong>Released at:</strong> {{ currentVersion?.published_at.toLocaleString() }}</p>
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
