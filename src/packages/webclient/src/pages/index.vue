<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getAllSpaces, getDownloadCountForSpace, getSpace} from "@/api/spaces/spaces.ts";
import {getAllVersions} from "@/api/versions/versions.ts";

const spaces = ref<Space[]>();

const totalSpaces = ref<number | null>(null);
const totalVersions = ref<number | null>(null);
const totalDownloads = ref<number | null>(null);
const statsLoading = ref(false);

onMounted(async () => {
  spaces.value = [];
  spaces.value.push(await getSpace("fn"));
  spaces.value.push(await getSpace("fc"));
  spaces.value.push(await getSpace("fh"));
  spaces.value.push(await getSpace("fa"));

  void fetchStats();
});

async function fetchStats() {
  statsLoading.value = true;
  try {
    const all = await getAllSpaces();
    totalSpaces.value = all.length;

    // Parallel requests for versions count per space
    const versionsPromises = all.map(s =>
      getAllVersions(s.id)
        .then(vs => vs.length)
        .catch(err => {
          console.error("Failed to fetch versions for", s.id, err);
          return 0;
        })
    );

    // Parallel requests for download counts per space
    const downloadsPromises = all.map(s =>
      getDownloadCountForSpace(s.id)
        .catch(err => {
          console.error("Failed to fetch downloads for", s.id, err);
          return 0;
        })
    );

    const versionsCounts = await Promise.all(versionsPromises);
    const downloadsCounts = await Promise.all(downloadsPromises);

    totalVersions.value = versionsCounts.reduce((a, b) => a + b, 0);
    totalDownloads.value = downloadsCounts.reduce((a, b) => a + b, 0);
  } catch (e) {
    console.error("Failed to fetch stats:", e);
  } finally {
    statsLoading.value = false;
  }
}
</script>

<template>
  <v-container>
    <v-row class="mt-16" justify="center">
      <v-col md="7">
        <v-img
          src="/banner.png"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col>
        <p class="text-h5 text-center text-primary">Download platform for products by FancyInnovations</p>
      </v-col>
    </v-row>

    <v-row class="mt-16" justify="center">
      <v-col md="2">
        <v-card
          class="card__border"
          color="#29152550"
          elevation="12"
          height="100%"
          rounded="xl"
          to="/explore/minecraft-plugins"
        >
          <v-card-title class="mt-2">
            <div class="d-flex align-center">
              <v-icon class="mr-2" color="green lighten-1">mdi-minecraft</v-icon>
              <p>Minecraft plugins</p>
            </div>
          </v-card-title>

          <v-card-text>
            FancyNpcs, FancyHolograms and more fancy Minecraft plugins made by FancyInnovations.
          </v-card-text>
        </v-card>
      </v-col>

      <v-col md="2">
        <v-card
          class="card__border"
          color="#29152550"
          elevation="12"
          height="100%"
          rounded="xl"
          to="/explore/hytale-plugins"
        >
          <v-card-title class="mt-2">
            <v-icon class="mr-2" color="blue lighten-1">mdi-sword-cross</v-icon>
            Hytale plugins
          </v-card-title>

          <v-card-text>
            Our first Hytale plugin FancyCore is a must-have for every Hytale server.
          </v-card-text>
        </v-card>
      </v-col>

      <v-col md="2">
        <v-card
          class="card__border"
          color="#29152550"
          elevation="12"
          height="100%"
          rounded="xl"
          to="/explore/other-projects"
        >
          <v-card-title class="mt-2">
            <v-icon class="mr-2" color="purple lighten-1">mdi-cube-outline</v-icon>
            Other projects
          </v-card-title>

          <v-card-text>
            Explore our other exciting projects and tools like FancyAnalytics or FancyVerteiler.
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row class="mt-10" justify="center">
      <v-col md="6">
        <v-card
          class="py-4"
          color="#29152550"
          elevation="12"
          height="100%"
          rounded="xl"
        >
          <v-row>
            <v-col class="text-center">
              <div class="text-h6">Projects</div>
              <div class="text-h4 mt-2">
                <v-skeleton-loader v-if="statsLoading" type="heading" />
                <template v-else>
                  {{ totalSpaces !== null ? totalSpaces : '—' }}
                </template>
              </div>
            </v-col>

            <v-col class="text-center">
              <div class="text-h6">Versions</div>
              <div class="text-h4 mt-2">
                <v-skeleton-loader v-if="statsLoading" type="heading" />
                <template v-else>
                  {{ totalVersions !== null ? totalVersions : '—' }}
                </template>
              </div>
            </v-col>

            <v-col class="text-center">
              <div class="text-h6">Downloads</div>
              <div class="text-h4 mt-2">
                <v-skeleton-loader v-if="statsLoading" type="heading" />
                <template v-else>
                  {{ totalDownloads !== null ? totalDownloads : '—' }}
                </template>
              </div>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <h1 class="text-h4 mt-16 text-center">Popular projects</h1>
      </v-col>
    </v-row>

    <v-row v-if="spaces && spaces.length > 0" justify="center">
      <v-col md="6">
        <v-carousel
          :show-arrows="false"
          cycle
          height="fit-content"
          interval="5000"
        >
          <v-carousel-item v-for="space in spaces" :key="space.id">
            <SpaceCard
              :space="space"
              :with-badge="true"
              class="pb-10"
            />
          </v-carousel-item>
        </v-carousel>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
