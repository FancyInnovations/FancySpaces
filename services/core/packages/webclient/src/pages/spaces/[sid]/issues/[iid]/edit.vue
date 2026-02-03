<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import {deleteIssue, getIssue, updateIssue} from "@/api/issues/issues.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import type {Issue} from "@/api/issues/types.ts";

const router = useRouter();
const route = useRoute();

const space = ref<Space>();
const issue = ref<Issue>();

onMounted(async () => {
  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.issue_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  const issueID = (route.params as any).iid as string;
  issue.value = await getIssue(spaceID, issueID);

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

async function editIssueReq() {
  if (!space.value || !issue.value) return;

  await updateIssue(space.value!.id, issue.value.id, issue.value);

  await router.push(`/spaces/${space.value?.slug}/issues/${issue.value.id}`);
}

async function deleteIssueReq() {
  if (!space.value || !issue.value) return;

  await deleteIssue(space.value!.id, issue.value.id);

  await router.push(`/spaces/${space.value?.slug}/issues`);
}

</script>

<template>
  <v-container v-if="issue" width="90%">
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
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <h1 class="text-center">Edit Issue #{{issue?.id}}</h1>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-text-field
          v-model="issue!.title"
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
          v-model="issue!.description"
          color="primary"
          hide-details
          label="Description"
          rows="8"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-text-field
          v-model="issue!.parent_issue"
          color="primary"
          hide-details
          label="Parent Issue"
          required
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
          <v-select
            v-model="issue!.status"
            :items="[
                    { title: 'Backlog', value: 'backlog' },
                    { title: 'Planned', value: 'planned' },
                    { title: 'In Progress', value: 'in_progress' },
                    { title: 'Done', value: 'done' },
                    { title: 'Closed', value: 'closed' },
                  ]"
            color="primary"
            hide-details
            label="Status"
            required
          />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="3">
        <v-select
          v-model="issue!.type"
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
          v-model="issue!.priority"
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
        <v-text-field
          v-model="issue!.assignee"
          color="primary"
          hide-details
          label="Issue Assignee"
          required
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-text-field
          v-model="issue!.fix_version"
          color="primary"
          hide-details
          label="Fix Version"
          required
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-btn
          class="mr-4"
          color="primary"
          variant="tonal"
          @click="editIssueReq"
        >
          Edit Issue
        </v-btn>

        <v-btn
          color="error"
          variant="tonal"
          @click="deleteIssueReq"
        >
          Delete Issue
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
