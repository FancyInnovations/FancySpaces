<script lang="ts" setup>

import {type Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";


const space = ref<Space>();

const displayType = ref<'board' | 'list'>('board');
const typeFilter = ref();
const priorityFilter = ref();
const statusFilter = ref();

onMounted(async () => {
  const spaceID = (useRoute().params as any).sid as string;
  space.value = await getSpace(spaceID);

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
          class="mt-4 grey-border-color"
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
            <div class="d-flex align-center justify-space-between">
              <div class="d-flex align-center">
                <v-select
                  v-model="typeFilter"
                  :items="[
                    { title: 'Epic', value: 'epic' },
                    { title: 'Bug', value: 'bug' },
                    { title: 'Task', value: 'task' },
                    { title: 'Storie', value: 'story' },
                    { title: 'Idea', value: 'idea' },

                  ]"
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
                  class="ml-4"
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
                    { title: 'TODO', value: 'todo' },
                    { title: 'In Progress', value: 'in_progress' },
                    { title: 'Done', value: 'done' },
                    { title: 'Closed', value: 'closed' },

                  ]"
                  class="ml-4"
                  clearable
                  color="primary"
                  density="compact"
                  hide-details
                  label="Status"
                  min-width="200"
                />
              </div>

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
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <IssueBoard
        v-if="displayType === 'board'"
        :priority-filter="priorityFilter"
        :space="space!"
        :status-filter="statusFilter"
        :type-filter="typeFilter"
      />
      <IssueTable
        v-else
        :priority-filter="priorityFilter"
        :space="space!"
        :status-filter="statusFilter"
        :type-filter="typeFilter"
      />
    </v-row>
  </v-container>
</template>

<style scoped>
.grey-border-color {
  border-color: rgba(0, 0, 0, 0.8);
}
</style>
