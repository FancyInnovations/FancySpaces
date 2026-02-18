<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import type {SpaceSecret} from "@/api/secrets/types.ts";
import {deleteSecret, getAllSecrets, getSecretDecrypted} from "@/api/secrets/secrets.ts";

const router = useRouter();

const space = ref<Space>();
const secrets = ref<SpaceSecret[]>();

const isLoggedIn = computed(() => {
  return localStorage.getItem("fs_api_key") !== null;
});

const tableHeaders = [
  { title: 'Key', key: 'key', sortable: true },
  { title: 'Description', key: 'description', sortable: false },
  { title: 'Created at', key: 'created_at', sortable: false, value: (s: SpaceSecret) => s.created_at.toLocaleString() },
  { title: 'Updated at', key: 'updated_at', sortable: false, value: (s: SpaceSecret) => s.updated_at.toLocaleString() },
  { title: '', key: 'actions', sortable: false, align: 'end' as any },
]

onMounted(async () => {
  const spaceID = (useRoute().params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!isLoggedIn || !space.value.secrets_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  secrets.value = await getAllSecrets(spaceID);

  useHead({
    title: `${space.value.title} secrets - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || `Explore the ${space.value.title} project space on FancySpaces.`
      }
    ]
  });
});

async function copySecretToClipboard(secret: SpaceSecret) {
  const decryptedValue = await getSecretDecrypted(secret.space_id, secret.key);

  navigator.clipboard.writeText(decryptedValue);
  window.alert("Copied secret to clipboard!");
}

async function deleteSecretReq(secret: SpaceSecret) {
  if (confirm(`Are you sure you want to delete the secret "${secret.key}"? This action cannot be undone.`)) {
    await deleteSecret(secret.space_id, secret.key);
    secrets.value = secrets.value?.filter(s => s.key !== secret.key);
  }
}

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
              :to="`/spaces/${space?.slug}/secrets/new`"
              color="primary"
              size="large"
              variant="tonal"
            >
              New Secret
            </v-btn>
          </div>
        </div>

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
          <v-card-text>
            <v-data-table
              :headers="tableHeaders"
              :items="secrets"
              class="bg-transparent"
            >
              <template v-slot:item.actions="{ item }">
                <div class="actions__width">
                  <v-btn
                    class="mr-4"
                    icon="mdi-link-variant"
                    variant="text"
                    @click="copySecretToClipboard(item)"
                  />

                  <v-btn
                    :to="`/spaces/${space?.slug}/secrets/${item.key}`"
                    class="mr-4"
                    icon="mdi-pencil"
                    variant="text"
                  />

                  <v-btn
                    color="error"
                    icon="mdi-delete"
                    variant="text"
                    @click="deleteSecretReq(item)"
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
