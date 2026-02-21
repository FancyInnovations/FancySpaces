<script lang="ts" setup>
import {useHead} from "@vueuse/head";
import Card from "@/components/common/Card.vue";
import UserHeader from "@/components/users/UserHeader.vue";
import {getDownloadCountForSpace, getSpacesOfCreator} from "@/api/spaces/spaces.ts";
import {mapLinkToDisplayname, mapLinkToIcon, type Space} from "@/api/spaces/types.ts";

const route = useRoute();

const user = ref<string>(); // TODO replace with public user info
const userLinks = ref<{ name: string; url: string }[]>([{name: "website", url: "https://github.com/oliverschlueter"}]); // TODO replace with public user info
const spaces = ref<Space[]>([]);
const totalDownloads = ref(0);

onMounted(async () => {
  user.value = (route.params as any).userid as string;

  spaces.value = await getSpacesOfCreator(user.value);

  for (let sp of spaces.value) {
    totalDownloads.value += await getDownloadCountForSpace(sp.id);
  }

  useHead({
    title: `${user.value} - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: `Check out ${user.value}'s profile on FancySpaces, showcasing their spaces.`
      }
    ]
  });
});

</script>

<template>
  <v-container width="60%">
    <v-row>
      <v-col>
        <UserHeader :userID="user"/>

        <hr
          class="mt-4 grey-border-color"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col md="8">
        <template v-for="space in spaces" :key="space.id" >
            <SpaceCard
              :space="space"
              :with-badge="true"
              class="mb-4"
            />
        </template>
      </v-col>

      <v-col md="4">
        <Card
          class="mb-4"
          min-width="200"
        >
          <v-card-title class="mt-2">Details</v-card-title>

          <v-card-text>
            <p class="text-body-1"><strong>ID:</strong> {{ user }}</p>
            <p class="text-body-1"><strong>Joined at:</strong> {{ '2026-01-01' }}</p>
            <p class="text-body-1"><strong>Spaces:</strong> {{ spaces.length }}</p>
            <p class="text-body-1"><strong>Total Downloads:</strong> {{ totalDownloads }}</p>
          </v-card-text>
        </Card>

        <Card
          v-if="userLinks.length > 0"
          class="mb-4"
          elevation="12"
        >
          <v-card-title class="mt-2">Links</v-card-title>

          <v-card-text>
            <div v-for="(link) in userLinks" :key="link.name" class="mb-2">
              <a
                :href="link.url"
                target="_blank"
              >
                <div class="d-flex align-center">
                  <v-icon
                    :icon="mapLinkToIcon(link.name)"
                    class="me-2"
                  />
                  <p class="text-body-1 link--hover">{{ mapLinkToDisplayname(link.name) }}</p>
                </div>
              </a>
            </div>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>

</template>

<style scoped>

</style>

<style>
.link--hover:hover {
  text-decoration: underline;
}

</style>
