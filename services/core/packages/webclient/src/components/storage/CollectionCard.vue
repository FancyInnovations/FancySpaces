<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {mapEngineKeyToName, type SpaceDatabaseCollection} from "@/api/storage/types.ts";
import {kvCount, kvSize} from "@/api/storage/kv/kv.ts";
import Card from "@/components/common/Card.vue";

const props = defineProps<{
  space?: Space,
  coll: SpaceDatabaseCollection
  withoutActions?: boolean
}>();

const count = ref<number>(-1);
const size = ref<number>(-1);

onMounted(async () => {
  if (props.coll.engine === "kv") {
    count.value = await kvCount(props.coll.database, props.coll.name);
    size.value = await kvSize(props.coll.database, props.coll.name);
  }
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
  <Card>
    <v-card-title class="mt-2">
      {{ coll.name }}
    </v-card-title>

    <v-card-text>
      <p><strong>Database:</strong> {{ coll.database }}</p>
      <p><strong>Collection:</strong> {{ coll.name }}</p>
      <p><strong>Engine:</strong> {{ mapEngineKeyToName(coll.engine) }}</p>
      <p><strong>Created at:</strong> {{ coll.created_at.toLocaleString() }}</p>
      <p v-if="count > 0"><strong>Count:</strong> {{ count }}</p>
      <p v-if="size > 0"><strong>Size:</strong> {{ formatSize(size) }}</p>
    </v-card-text>

    <v-card-actions v-if="!withoutActions">
      <v-btn
        :to="`/spaces/${space?.slug}/storage/${coll.database}/${coll.name}`"
        color="primary"
        variant="text"
      >
        View Data
      </v-btn>

      <v-btn
        :to="`/spaces/${space?.slug}/storage/${coll.database}/${coll.name}/settings`"
        color="primary"
        variant="text"
      >
        Settings
      </v-btn>
    </v-card-actions>
  </Card>
</template>

<style scoped>

</style>
