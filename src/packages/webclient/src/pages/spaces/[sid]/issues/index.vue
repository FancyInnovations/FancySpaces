<script lang="ts" setup>

import {type Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import {getAllIssues} from "@/api/issues/issues.ts";
import type {Issue} from "@/api/issues/types.ts";

const isLoggedIn = computed(() => {
  return localStorage.getItem("fs_api_key") !== null;
});

const space = ref<Space>();
const issues = ref<Issue[]>([]);
const filteredIssues = computed(() => {
  return issues.value.filter(issue => {
    const matchesSearch = searchQuery.value ?
      issue.title.toLowerCase().includes(searchQuery.value.toLowerCase()) :
      true;

    const matchesType = typeFilter.value ? issue.type === typeFilter.value : true;
    const matchesPriority = priorityFilter.value ? issue.priority === priorityFilter.value : true;
    const matchesStatus = statusFilter.value ? issue.status === statusFilter.value : true;

    return matchesSearch && matchesType && matchesPriority && matchesStatus;
  });
});

const displayType = ref<'board' | 'list'>('board');
const searchQuery = ref('');
const typeFilter = ref();
const priorityFilter = ref();
const statusFilter = ref();

onMounted(async () => {
  const spaceID = (useRoute().params as any).sid as string;
  space.value = await getSpace(spaceID);

  issues.value = await getAllIssues(space.value.id);

  useHead({
    title: `${space.value.title} issues - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || 'View issues for this space on FancySpaces.'
      }
    ]
  });
});

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
            The issues feature is currently in beta and may not function as expected. We appreciate your patience as we work to improve it!
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
              :to="`/spaces/${space?.slug}/issues/new`"
              color="primary"
              size="large"
              variant="tonal"
            >
              New Issue
            </v-btn>
          </div>
        </div>

        <hr
          class="grey-border-color"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col
        class="d-flex align-center justify-space-between flex-wrap"
      >
        <v-card
          class="card__border"
          color="#19120D33"
          elevation="12"
          rounded="xl"
        >
          <v-card-text>
            <div class="d-flex align-center justify-space-between">
              <div class="d-flex align-center flex-wrap">
                <v-text-field
                  v-model="searchQuery"
                  class="ma-2"
                  clearable
                  color="primary"
                  density="compact"
                  hide-details
                  label="Search issues"
                  min-width="300"
                  prepend-inner-icon="mdi-magnify"
                />

                <v-select
                  v-model="typeFilter"
                  :items="[
                    { title: 'Epic', value: 'epic' },
                    { title: 'Bug', value: 'bug' },
                    { title: 'Task', value: 'task' },
                    { title: 'Story', value: 'story' },
                    { title: 'Idea', value: 'idea' },

                  ]"
                  class="ma-2"
                  clearable
                  color="primary"
                  density="compact"
                  hide-details
                  label="Type"
                  min-width="200"
                />

                <v-select
                  v-model="priorityFilter"
                  :items="[
                    { title: 'Low', value: 'low' },
                    { title: 'Medium', value: 'medium' },
                    { title: 'High', value: 'high' },
                    { title: 'Critical', value: 'critical' },
                  ]"
                  class="ma-2"
                  clearable
                  color="primary"
                  density="compact"
                  hide-details
                  label="Priority"
                  min-width="200"
                />

                <v-select
                  v-model="statusFilter"
                  :items="[
                    { title: 'Backlog', value: 'backlog' },
                    { title: 'Planned', value: 'planned' },
                    { title: 'In Progress', value: 'in_progress' },
                    { title: 'Done', value: 'done' },
                    { title: 'Closed', value: 'closed' },

                  ]"
                  class="ma-2"
                  clearable
                  color="primary"
                  density="compact"
                  hide-details
                  label="Status"
                  min-width="200"
                />
              </div>
            </div>
          </v-card-text>
        </v-card>

        <v-card
          class="card__border margin-top flex-grow-1"
          color="#19120D33"
          elevation="12"
          max-width="fit-content"
          rounded="xl"
        >
          <v-card-text>
            <v-btn-group
              density="compact"
            >
              <v-btn
                :variant="displayType === 'board' ? 'tonal' : 'outlined'"
                color="primary"
                prepend-icon="mdi-view-dashboard"
                @click="displayType = 'board'"
              >
                Board
              </v-btn>

              <v-btn
                :variant="displayType === 'list' ? 'tonal' : 'outlined'"
                color="primary"
                prepend-icon="mdi-format-list-bulleted"
                @click="displayType = 'list'"
              >
                List
              </v-btn>
            </v-btn-group>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <IssueBoard
          v-if="displayType === 'board'"
          :issues="filteredIssues"
          :space="space!"
        />
        <IssueTable
          v-else
          :issues="filteredIssues"
          :space="space!"
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.grey-border-color {
  border-color: rgba(0, 0, 0, 0.8);
}

@media (max-width: 1919px) {
  .margin-top {
    margin-top: 16px !important;
  }
}
</style>
