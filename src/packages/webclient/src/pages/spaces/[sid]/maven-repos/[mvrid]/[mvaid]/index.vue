<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import type {SpaceMavenRepository, SpaceMavenRepositoryArtifact} from "@/api/maven/types.ts";
import {getMavenArtifacts, getMavenRepository} from "@/api/maven/maven.ts";

const route = useRoute();
const router = useRouter();

const isLoggedIn = computed(() => {
  return localStorage.getItem("fs_api_key") !== null;
});

const space = ref<Space>();

const repo = ref<SpaceMavenRepository>();
const artifact = ref<SpaceMavenRepositoryArtifact>();
const sortedVersions = computed(() => {
  return artifact.value?.versions.slice().sort((a, b) => {
    return new Date(b.published_at).getTime() - new Date(a.published_at).getTime();
  }) || [];
});

const expanded = ref(<{ [key: string]: boolean }>({}));
function toggleExpand(key: any) {
  expanded.value[key] = !expanded.value[key]
}

const howToUseTab = ref("build.gradle.kts");

onMounted(async () => {
  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.maven_repository_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  const mavenRepoName = (route.params as any).mvrid as string;
  repo.value = await getMavenRepository(space.value.id, mavenRepoName);

  const groupArtifactID = (route.params as any).mvaid as string;
  artifact.value = await getMavenArtifacts(space.value.id, repo.value.name, groupArtifactID);

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

function filesWithoutChecksums(ver: string) {
  return artifact.value?.versions
    .find(v => v.version === ver)?.files
    .filter(
      f => !f.name.endsWith('.md5') && !f.name.endsWith('.sha1') && !f.name.endsWith('.sha256') && !f.name.endsWith('.sha512')
    )
    || [];
}

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
  <v-container width="90%">
    <v-row>
      <v-col>
        <v-card
          color="error"
          rounded="xl"
          variant="tonal"
        >
          <v-card-text>
            <v-icon
              class="mr-2"
            >
              mdi-alert-circle
            </v-icon>
            The maven repository feature is currently in beta and may not function as expected. We appreciate your patience as we work to improve it!
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

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
              color="primary"
              disabled
              size="large"
              variant="tonal"
            >
              New Repo
            </v-btn>
          </div>
        </div>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col class="mb-4">
        <v-card
          class="card__border bg-transparent"
          color="#150D1950"
          elevation="12"
          min-width="600"
          rounded="xl"
        >
          <v-card-text>
            <v-breadcrumbs
              :items="[
                { title: 'Maven Repositories', to: `/spaces/${space?.slug}/maven-repos` },
                { title: repo?.name || '', to: `/spaces/${space?.slug}/maven-repos/${repo?.name}` },
                { title: `${artifact?.group}:${artifact?.id}` || '' }
              ]"
              class="pa-0"
              color="primary"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row
    >
      <v-col
        v-for="version in sortedVersions"
        :key="version.version"
        md="12"
      >
        <v-card
          class="card__border bg-transparent"
          color="#150D1950"
          elevation="12"
          rounded="xl"
        >
          <v-card-title
            class="d-flex align-center justify-space-between cursor-pointer"
            @click="toggleExpand(version.version)"
          >
            <span>
              Version {{ version.version }}
            </span>

            <v-btn
              icon
              size="small"
              variant="text"
              @click="toggleExpand(version.version)"
            >
              <v-icon>
                {{ expanded[version.version] ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
              </v-icon>
            </v-btn>
          </v-card-title>

          <v-card-subtitle class="pb-4">
            Published at {{ version.published_at.toLocaleString() }}
          </v-card-subtitle>

          <v-expand-transition>
            <div v-show="expanded[version.version]">
              <v-card-text class="d-flex align-start justify-space-between">

                <v-card
                  class="card__border bg-transparent"
                  color="#150D1950"
                  elevation="12"
                  rounded="xl"
                  width="100%"
                >
                  <v-card-title class="mt-2">Available files</v-card-title>

                  <v-card-text>
                    <v-table class="bg-transparent">
                      <thead>
                      <tr>
                        <th class="text-left">Name</th>
                        <th class="text-left">Size</th>
                        <th class="text-right"></th>
                      </tr>
                      </thead>

                      <tbody>
                      <tr
                        v-for="file in filesWithoutChecksums(version.version)"
                        :key="file.name"
                      >
                        <td>{{ file.name }}</td>
                        <td>{{ formatSize(file.size) }}</td>
                        <td class="text-right">
                          <v-btn
                            :href="file.url"
                            color="primary"
                            icon="mdi-download"
                            size="small"
                            target="_blank"
                            variant="text"
                          />
                        </td>
                      </tr>
                      </tbody>
                    </v-table>
                  </v-card-text>
                </v-card>

                <v-card
                  class="card__border bg-transparent flex-grow-1 ml-4"
                  color="#19120D33"
                  elevation="6"
                  min-width="50%"
                  rounded="xl"
                >
                  <v-card-title class="mt-2">How to use</v-card-title>

                  <v-card-text>
                    <v-tabs
                      v-model="howToUseTab"
                      background-color="#150D1950"
                      class="mt-2"
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
    maven (url = "https://fancyspaces.net/maven/{{ space?.slug }}/{{ repo?.name }}")
}

dependencies {
    compileOnly("{{ artifact?.group }}:{{ artifact?.id }}:{{ version.version }}")
}</code></pre>
                      </v-tabs-window-item>
                      <v-tabs-window-item value="build.gradle">
                <pre><code>repositories {
    maven {
        url "https://fancyspaces.net/maven/{{ space?.slug }}/{{ repo?.name }}"
    }
}

dependencies {
    compileOnly "{{ artifact?.group }}:{{ artifact?.id }}:{{ version.version }}"
}</code></pre>
                      </v-tabs-window-item>
                      <v-tabs-window-item value="pom.xml">
                <pre><code>&lt;repositories&gt;
    &lt;repository&gt;
        &lt;id&gt;fancyspaces-{{ space?.slug }}-{{ repo?.name }}&lt;/id&gt;
        &lt;url&gt;https://fancyspaces.net/maven/{{ space?.slug }}/{{ repo?.name }}&lt;/url&gt;
    &lt;/repository&gt;
&lt;/repositories&gt;

&lt;dependencies&gt;
    &lt;dependency&gt;
        &lt;groupId&gt;{{ artifact?.group }}&lt;/groupId&gt;
        &lt;artifactId&gt;{{ artifact?.id }}&lt;/artifactId&gt;
        &lt;version&gt;{{ version.version }}&lt;/version&gt;
        &lt;scope&gt;provided&lt;/scope&gt;
    &lt;/dependency&gt;
&lt;/dependencies&gt;</code></pre>
                      </v-tabs-window-item>
                    </v-tabs-window>
                  </v-card-text>
                </v-card>
              </v-card-text>
            </div>
          </v-expand-transition>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.grey-border-color {
  border-color: rgba(0, 0, 0, 0.8);
}

pre {
  overflow-x: auto;
}
</style>
