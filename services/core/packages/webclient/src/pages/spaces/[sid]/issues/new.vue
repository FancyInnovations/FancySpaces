<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import {createIssue} from "@/api/issues/issues.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import SpaceHeader from "@/components/SpaceHeader.vue";

const router = useRouter();
const route = useRoute();

const space = ref<Space>();

const title = ref('');
const description = ref('');
const type = ref('task');
const priority = ref('medium');

onMounted(async () => {
  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.issue_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

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
      <v-col class="flex-grow-0 pa-0">
        <SpaceSidebar
          :space="space"
        />
      </v-col>

      <v-col>
        <SpaceHeader :space="space">
          <template #quick-actions>
            <v-btn
              :to="`/spaces/${space?.slug}/issues`"
              color="primary"
              size="large"
              variant="tonal"
            >
              View Issues
            </v-btn>
          </template>
        </SpaceHeader>

        <hr
          class="grey-border-color mt-4"
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

</style>
