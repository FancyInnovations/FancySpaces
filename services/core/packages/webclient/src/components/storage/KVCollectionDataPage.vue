<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import type {KVValue, SpaceDatabase, SpaceDatabaseCollection} from "@/api/storage/types.ts";
import {kvGetAll} from "@/api/storage/kv/kv.ts";

const props = defineProps<{
  space: Space
  database: SpaceDatabase
  collection: SpaceDatabaseCollection
}>();

const values = ref<KVValue[]>([]);

const tableHeaders = [
  { title: 'Key', value: 'key' },
  { title: 'Type', value: 'type' },
  { title: 'Value', value: 'value' },
];

onMounted(async () => {
  const data = await kvGetAll(props.database.name, props.collection.name);
  for (const key in data) {
    values.value.push({
      key: key,
      value: data[key],
      type: typeof data[key]
    });
  }
});

</script>

<template>
  <v-row>
    <v-col>
      <v-card
        class="card__border"
        color="#19120D33"
        elevation="12"
        rounded="xl"
      >
        <v-card-title class="mt-2">All key-value pairs</v-card-title>

        <v-card-text>
          <v-data-table
            :headers="tableHeaders"
            :items="values"
            class="bg-transparent"
          >
            <template #item.value="{ item }">
              <pre>{{ JSON.stringify(item.value, null, 2) }}</pre>
            </template>
          </v-data-table>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col md="3">
      <CollectionCard
        :coll="collection"
        :space="space"
        without-actions
      />
    </v-col>
  </v-row>
</template>

<style scoped>

</style>
