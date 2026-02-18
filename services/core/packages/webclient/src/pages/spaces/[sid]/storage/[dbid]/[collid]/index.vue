<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import {type SpaceDatabase, type SpaceDatabaseCollection} from "@/api/storage/types.ts";
import KVCollectionDataPage from "@/components/storage/KVCollectionDataPage.vue";
import SpaceHeader from "@/components/SpaceHeader.vue";

const router = useRouter();
const route = useRoute();

const isLoggedIn = computed(() => {
  return localStorage.getItem("fs_api_key") !== null;
});

const space = ref<Space>();
const database = ref<SpaceDatabase>();
const collection = ref<SpaceDatabaseCollection>();

onMounted(async () => {
  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.maven_repository_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  const databaseName = (route.params as any).dbid as string;
  database.value = {
    name: databaseName,
    created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
  };

  const collectionName = (route.params as any).collid as string;
  collection.value = {
    database: databaseName,
    name: collectionName,
    created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
    engine: "kv"
  };

  useHead({
    title: `${space.value.title} storage - FancySpaces`,
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
          <template #quick-actions>
            <v-btn
              v-if="isLoggedIn"
              :to="`/spaces/${space?.slug}/storage/new`"
              color="primary"
              disabled
              size="large"
              variant="tonal"
            >
              New collection
            </v-btn>
          </template>
        </SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>
    <KVCollectionDataPage
      v-if="space && database && collection && collection.engine === 'kv'"
      :collection="collection"
      :database="database"
      :space="space"
    />
  </v-container>
</template>

<style scoped>
.grey-border-color {
  border-color: rgba(0, 0, 0, 0.8);
}

</style>
