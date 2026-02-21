<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import type {SpaceMavenRepository, SpaceMavenRepositoryArtifact} from "@/api/maven/types.ts";
import {getAllMavenArtifacts, getMavenRepository} from "@/api/maven/maven.ts";
import SpaceHeader from "@/components/SpaceHeader.vue";
import {useUserStore} from "@/stores/user.ts";
import Card from "@/components/common/Card.vue";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();

const isLoggedIn = ref(false);

const space = ref<Space>();

const repo = ref<SpaceMavenRepository>();
const artifacts = ref<SpaceMavenRepositoryArtifact[]>([]);

const howToUseTab = ref("build.gradle.kts");

const tableHeaders = [
  { title: 'Group ID', value: 'group' },
  { title: 'Artifact ID', value: 'id' },
  { title: 'Versions', key: 'versions', value: (art: SpaceMavenRepositoryArtifact) => art.versions.length || 'N/A' },
  { title: 'Latest version', key: 'latest-version', value: (art: SpaceMavenRepositoryArtifact) => {
      const latest = art.versions.sort((a, b) => {
        return new Date(b.published_at).getTime() - new Date(a.published_at).getTime();
      })[0];
      return latest ? latest.version : 'N/A';
    }
  },
  { title: 'Last update', key: 'last-update', value: (art: SpaceMavenRepositoryArtifact) => {
      const latest = art.versions.sort((a, b) => {
        return new Date(b.published_at).getTime() - new Date(a.published_at).getTime();
      })[0];
      return latest ? new Date(latest.published_at).toLocaleDateString() : 'N/A';
    }
  },
];

onMounted(async () => {
  isLoggedIn.value = await userStore.isAuthenticated;

  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.maven_repository_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  const mavenRepoName = (route.params as any).mvrid as string;
  repo.value = await getMavenRepository(space.value.id, mavenRepoName);

  artifacts.value = await getAllMavenArtifacts(space.value.id, repo.value.name);

  useHead({
    title: `${space.value.title} Maven Repo ${repo.value.name} - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || `Explore the ${space.value.title} project space on FancySpaces.`
      }
    ]
  });
});

function onRowClick(event: any, { item }: any) {
  router.push(`/spaces/${space.value?.slug}/maven-repos/${repo.value?.name}/${item.group}:${item.id}`);
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
        <SpaceHeader :space="space">
          <template #metadata>
            <p class="text-body-2 mx-4">-</p>
            <p class="text-body-2">{{ artifacts?.length }} artifacts</p>
          </template>

          <template #quick-actions>
            <v-btn
              v-if="isLoggedIn"
              :to="`/spaces/${space?.slug}/maven-repos/new`"
              color="primary"
              disabled
              size="large"
              variant="tonal"
            >
              New Repo
            </v-btn>
          </template>
        </SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col class="mb-4">
        <Card
          class="bg-transparent"
          color="#150D1950"
          min-width="600"
        >
          <v-card-text>
            <v-breadcrumbs
              :items="[
                { title: 'Maven Repositories', to: `/spaces/${space?.slug}/maven-repos` },
                { title: repo?.name || '', to: `/spaces/${space?.slug}/maven-repos/${repo?.name}` },
              ]"
              class="pa-0"
              color="primary"
            />
          </v-card-text>
        </Card>
      </v-col>
    </v-row>

    <v-row>
      <v-col md="7">
        <Card
          class="bg-transparent"
          color="#150D1950"
          min-width="600"
        >
          <v-card-title class="mt-2">
            Artifacts in {{ repo?.name }}
          </v-card-title>

          <v-card-text>
            <v-data-table
              :headers="tableHeaders"
              :items="artifacts"
              class="bg-transparent"
              hover
              @click:row="onRowClick"
            >
            </v-data-table>
          </v-card-text>
        </Card>
      </v-col>

      <v-col md="5">
        <Card
          class="bg-transparent mb-4"
          elevation="6"
        >
          <v-card-title class="mt-2">How to use</v-card-title>

          <v-card-text>
            <v-tabs
              v-model="howToUseTab"
              background-color="#150D1950"
              color="primary"
              grow
            >
              <v-tab value="build.gradle.kts">build.gradle.kts</v-tab>
              <v-tab value="build.gradle">build.gradle</v-tab>
              <v-tab value="pom.xml">pom.xml</v-tab>
            </v-tabs>

            <v-tabs-window v-model="howToUseTab" class="mt-4">
              <v-tabs-window-item value="build.gradle.kts">
                <pre><code>repositories {
    maven (url = "https://maven.fancyspaces.net/{{ space?.slug }}/{{ repo?.name }}")
}</code></pre>
              </v-tabs-window-item>
              <v-tabs-window-item value="build.gradle">
                <pre><code>repositories {
    maven {
        url "https://maven.fancyspaces.net/{{ space?.slug }}/{{ repo?.name }}"
    }
}</code></pre>
              </v-tabs-window-item>
              <v-tabs-window-item value="pom.xml">
                <pre><code>&lt;repositories&gt;
    &lt;repository&gt;
        &lt;id&gt;fancyspaces-{{ space?.slug }}-{{ repo?.name }}&lt;/id&gt;
        &lt;url&gt;https://maven.fancyspaces.net/{{ space?.slug }}/{{ repo?.name }}&lt;/url&gt;
    &lt;/repository&gt;
&lt;/repositories&gt;</code></pre>
              </v-tabs-window-item>
            </v-tabs-window>
          </v-card-text>
        </Card>

        <Card
          v-if="isLoggedIn"
          class="bg-transparent"
          elevation="6"
        >
          <v-card-text>
            <v-btn
              :to="`/spaces/${space?.slug}/maven-repo/${repo?.name}/edit`"
              block
              class="mb-2"
              color="primary"
              variant="tonal"
            >
              Edit Repo
            </v-btn>

            <v-btn
              :to="`/spaces/${space?.slug}/maven-repo`"
              block
              color="error"
              variant="tonal"
            >
              Delete Repo
            </v-btn>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
pre {
  overflow-x: auto;
}
</style>
