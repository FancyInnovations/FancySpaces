<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import {type SpaceDatabaseCollection} from "@/api/storage/types.ts";
import CollectionCard from "@/components/storage/CollectionCard.vue";
import SpaceHeader from "@/components/SpaceHeader.vue";
import {useUserStore} from "@/stores/user.ts";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const isLoggedIn = ref(false);

const space = ref<Space|undefined>();
const collections = ref<SpaceDatabaseCollection[]>();

const collectionsByEngine = computed(() => {
  const map: Record<string, SpaceDatabaseCollection[]> = {};
  collections.value?.forEach(coll => {
    if (!map[coll.engine]) {
      map[coll.engine] = [];
    }
    map[coll.engine]!.push(coll);
  });

  return map;
});

onMounted(async () => {
  isLoggedIn.value = await userStore.isAuthenticated;

  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.maven_repository_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  collections.value = [
    {
      database: "fancyanalytics",
      name: "users",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "document"
    },
    {
      database: "fancyanalytics",
      name: "auth-tokens",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "document"
    },
    {
      database: "fancyanalytics",
      name: "auth-audits",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "document"
    },
    {
      database: "fancyanalytics",
      name: "projects",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "document"
    },
    {
      database: "fancyanalytics",
      name: "dashboards",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "document"
    },
    {
      database: "fancyanalytics",
      name: "metrics",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "document"
    },
    {
      database: "fancyanalytics",
      name: "templates",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "document"
    },
  ];

  collections.value.push(
    {
      database: "system",
      name: "kv_test",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "kv"
    },
    {
      database: "fancyanalytics",
      name: "metadata-cache",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "kv"
    },
    {
      database: "fancyanalytics",
      name: "auth-cache",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "kv"
    },
    {
      database: "fancyanalytics",
      name: "metric-records-cache",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "kv"
    },
    {
      database: "fancyanalytics",
      name: "logs-cache",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "kv"
    },
    {
      database: "fancyanalytics",
        name: "events-cache",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "kv"
    },
  );

  collections.value.push(
    {
      database: "fancyanalytics",
      name: "metric-records",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "analytical"
    },
    {
      database: "fancyanalytics",
      name: "event-records",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "analytical"
    },
    {
      database: "fancyanalytics",
      name: "log-records",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "analytical"
    },
  );

  collections.value.push(
    {
      database: "fancyanalytics",
      name: "ingest-queue",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "broker"
    },
    {
      database: "fancyanalytics",
      name: "internal-event-bus",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "broker"
    },
  );

  collections.value.push(
    {
      database: "fancyanalytics",
      name: "backups",
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 4),
      engine: "object"
    },
  );

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

    <template v-if="collectionsByEngine['document']">
      <v-row justify="center">
        <v-col>
          <h2 class="text-h5">Document stores</h2>
        </v-col>
      </v-row>

      <v-row>
        <v-col
          v-for="coll in collectionsByEngine['document']"
          :key="coll.name"
          md="3"
        >
          <CollectionCard :coll="coll" :space="space" />
        </v-col>
      </v-row>
    </template>

    <template v-if="collectionsByEngine['kv']">
      <v-row class="mt-8">
        <v-col>
          <h2 class="text-h5">Key-value stores</h2>
        </v-col>
      </v-row>

      <v-row>
        <v-col
          v-for="coll in collectionsByEngine['kv']"
          :key="coll.name"
          md="3"
        >
          <CollectionCard :coll="coll" :space="space" />
        </v-col>
      </v-row>
    </template>

    <template v-if="collectionsByEngine['analytical']">
      <v-row class="mt-8">
        <v-col>
          <h2 class="text-h5">Analytical stores</h2>
        </v-col>
      </v-row>

      <v-row>
        <v-col
          v-for="coll in collectionsByEngine['analytical']"
          :key="coll.name"
          md="3"
        >
          <CollectionCard :coll="coll" :space="space" />
        </v-col>
      </v-row>
    </template>

    <template v-if="collectionsByEngine['object']">
      <v-row class="mt-8">
        <v-col>
          <h2 class="text-h5">Object stores</h2>
        </v-col>
      </v-row>

      <v-row>
        <v-col
          v-for="coll in collectionsByEngine['object']"
          :key="coll.name"
          md="3"
        >
          <CollectionCard :coll="coll" :space="space" />
        </v-col>
      </v-row>
    </template>

    <template v-if="collectionsByEngine['broker']">
      <v-row class="mt-8">
        <v-col>
          <h2 class="text-h5">Message brokers</h2>
        </v-col>
      </v-row>

      <v-row>
        <v-col
          v-for="coll in collectionsByEngine['broker']"
          :key="coll.name"
          md="3"
        >
          <CollectionCard :coll="coll" :space="space" />
        </v-col>
      </v-row>
    </template>

  </v-container>
</template>

<style scoped>

</style>
