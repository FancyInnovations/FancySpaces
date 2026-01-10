<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getDownloadCountForSpace, getSpace} from "@/api/spaces/spaces.ts";
import {mapPlatformToDisplayname, type SpaceVersion} from "@/api/versions/types.ts";
import {getAllVersions, getDownloadCountForVersion, getLatestVersion} from "@/api/versions/versions.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import SpaceHeader from "@/components/SpaceHeader.vue";
import {useHead} from "@vueuse/head";

const router = useRouter();

const space = ref<Space>();
const latestVersion = ref<SpaceVersion>();
const versions = ref<SpaceVersion[]>();
const downloadCount = ref<number>(0);
const downloadCounts = ref<Record<string, number>>({});

const tableHeaders = [
  { title: 'Version', key: 'name', sortable: false },
  { title: 'Channel', key: 'channel', value: (ver: SpaceVersion) => ver.channel.toUpperCase(), sortable: false },
  { title: 'Platform', key: 'platform', value: (ver: SpaceVersion) => mapPlatformToDisplayname(ver.platform), sortable: false },
  { title: 'Platform versions', key: 'supported_platform_versions', sortable: false, value: (ver: SpaceVersion) => ver.supported_platform_versions.join(", "), class: 'platform-versions__max-width' },
  { title: 'Released at', key: 'published_at', sortable: false, value: (ver: SpaceVersion) => new Date(ver.published_at).toLocaleString() },
  { title: 'Downloads', key: 'downloads', sortable: false, value: (ver: SpaceVersion) => downloadCounts.value[ver.id] || 0 },
  { title: '', key: 'actions', sortable: false, align: 'end' as any },
]

onMounted(async () => {
  const spaceID = (useRoute().params as any).sid as string;
  space.value = await getSpace(spaceID);
  latestVersion.value = await getLatestVersion(space.value.id);
  versions.value = await getAllVersions(space.value.id);

  downloadCount.value = await getDownloadCountForSpace(space.value.id);
  for (const version of versions.value) {
    downloadCounts.value[version.id] = await getDownloadCountForVersion(version.space_id, version.id);
  }

  useHead({
    title: `${space.value.title} versions - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || `Explore the ${space.value.title} project space on FancySpaces.`
      }
    ]
  });
});

function onRowClick(event: any, { item }: any) {
  router.push(`/spaces/${space.value?.slug}/versions/${item.name}`);
}

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text);
  window.alert("Copied to clipboard!");
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
          :download-count="downloadCount"
          :latest-version="latestVersion"
          :space="space"
        />

        <hr
          class="mt-4 grey-border-color"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <v-card
          class="card__border"
          color="#29152550"
          elevation="12"
          rounded="xl"
        >
          <v-card-text>
            <v-data-table
              :headers="tableHeaders"
              :items="versions"
              class="bg-transparent"
              hover
              @click:row="onRowClick"
            >
              <template v-slot:item.actions="{ item }">
                <div class="actions__width">
                  <v-btn
                    v-if="item.files.length != 1"
                    :to="`/spaces/${space?.slug}/versions/${item.name}`"
                    class="mr-4 my-1"
                    icon="mdi-download"
                    variant="text"
                  />
                  <v-btn
                    v-else
                    :href="item.files[0]?.url"
                    class="mr-4 my-1"
                    icon="mdi-download"
                    variant="text"
                  />

                  <v-btn
                    icon="mdi-link-variant"
                    variant="text"
                    @click="copyToClipboard(`https://fancyspaces.net/spaces/${space?.slug}/versions/${item.name}`)"
                  />
                </div>
              </template>

            </v-data-table>
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

table, tr, td, thead, tbody {
  background: transparent;
  border-collapse: collapse;
}

.actions__width {
  min-width: 130px;
}
</style>
