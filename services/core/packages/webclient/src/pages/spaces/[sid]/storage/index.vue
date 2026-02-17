<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import {type SpaceDatabaseCollection} from "@/api/storage/types.ts";
import CollectionCard from "@/components/storage/CollectionCard.vue";

const router = useRouter();

const isLoggedIn = computed(() => {
  return localStorage.getItem("fs_api_key") !== null;
});

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
  const spaceID = (useRoute().params as any).sid as string;
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
              class="mb-2"
              color="primary"
              size="large"
              variant="tonal"
            >
              New database
            </v-btn>

            <v-btn
              v-if="isLoggedIn"
              :to="`/spaces/${space?.slug}/maven-repos/new`"
              color="primary"
              size="large"
              variant="tonal"
            >
              New collection
            </v-btn>
          </div>
        </div>

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
.grey-border-color {
  border-color: rgba(0, 0, 0, 0.8);
}

</style>
