<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import {createIssue} from "@/api/issues/issues.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";

const router = useRouter();

const space = ref<Space>();

const title = ref('');
const description = ref('');
const type = ref('task');
const priority = ref('medium');

onMounted(async () => {
  const spaceID = (useRoute().params as any).sid as string;
  space.value = await getSpace(spaceID);

  useHead({
    title: `${space.value.title} - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || 'Create a new issue in this space on FancySpaces.'
      }
    ]
  });
});

async function createNewIssue() {
  const issue = await createIssue(space.value!.id, {
    title: title.value,
    description: description.value,
    type: type.value as any,
    priority: priority.value as any
  });

  title.value = '';
  description.value = '';
  type.value = 'task';
  priority.value = 'medium';

  await router.push(`/spaces/${space.value?.slug}/issues/${issue.id}`);
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
        </div>

        <hr
          class="grey-border-color"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <h1 class="text-center">Create new issue in {{ space?.title }}</h1>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-text-field
          v-model="title"
          color="primary"
          hide-details
          label="Title"
          required
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-textarea
          v-model="description"
          color="primary"
          hide-details
          label="Description"
          rows="8"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="3">
        <v-select
          v-model="type"
          :items="[
                    { title: 'Epic', value: 'epic' },
                    { title: 'Bug', value: 'bug' },
                    { title: 'Task', value: 'task' },
                    { title: 'Story', value: 'story' },
                    { title: 'Idea', value: 'idea' },

                  ]"
          color="primary"
          hide-details
          label="Type"
          required
        />
      </v-col>

      <v-col md="3">
        <v-select
          v-model="priority"
          :items="[
                    { title: 'Low', value: 'low' },
                    { title: 'Medium', value: 'medium' },
                    { title: 'High', value: 'high' },
                    { title: 'Critical', value: 'critical' },
                  ]"
          color="primary"
          hide-details
          label="Priority"
          required
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-btn
          color="primary"
          @click="createNewIssue"
        >
          Create Issue
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.grey-border-color {
  border-color: rgba(0, 0, 0, 0.8);
}
</style>
