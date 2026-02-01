<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import type {Issue, IssueComment} from "@/api/issues/types.ts";
import {deleteIssue, getIssue, updateIssue} from "@/api/issues/issues.ts";

const route = useRoute();
const router = useRouter();

const isLoggedIn = computed(() => {
  return localStorage.getItem("fs_api_key") !== null;
});

const space = ref<Space>();

const currentIssue = ref<Issue>();
const comments = ref<IssueComment[]>([]);

onMounted(async () => {
  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  const issueID = (route.params as any).iid as string;
  currentIssue.value = await getIssue(spaceID, issueID);

  useHead({
    title: `${space.value.title} Issue #${currentIssue.value?.id} - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || `View issue #${currentIssue.value?.id} for this space on FancySpaces.`
      }
    ]
  });
});

async function deleteIssueReq() {
  if (!space.value || !currentIssue.value) return;

  await deleteIssue(space.value!.id, currentIssue.value.id);

  await router.push(`/spaces/${space.value?.slug}/issues`);
}

async function statusChanged(newStatus: string) {
  currentIssue.value!.status = newStatus as any;
  await updateIssue(currentIssue.value!.space, currentIssue.value!.id, currentIssue.value!);
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
      <v-col md="9">
        <v-card
          class="card__border mb-4"
          color="#19120D33"
          elevation="12"
          rounded="xl"
        >
          <v-card-title class="my-2">{{ currentIssue?.title }}</v-card-title>
        </v-card>

        <v-card
          class="card__border mb-4"
          color="#19120D33"
          elevation="12"
          rounded="xl"
        >
          <v-card-text class="py-0">
            <MarkdownRenderer
              :markdown="currentIssue?.description"
            />
          </v-card-text>
        </v-card>

        <v-card
          class="card__border bg-transparent"
          color="#150D1950"
          elevation="12"
          min-width="600"
          rounded="xl"
        >
          <v-card-title class="mt-2">
            Comments ({{ comments?.length }})
          </v-card-title>

          <v-card-text>
            <p v-if="comments?.length === 0">
              No comments yet.
            </p>
            <div v-else class="issue-comments">
              <v-card
                v-for="comment in comments"
                :key="comment.id"
                class="card__border bg-transparent mb-3"
                color="#19120D33"
                elevation="6"
                rounded="xl"
              >
                <v-card-text>
                  <div class="d-flex justify-space-between mb-2">
                    <div class="d-flex align-center">
                      <span class="font-weight-medium">{{ comment.author }}</span>
                    </div>
                    <span class="text-caption grey--text">{{ comment.created_at.toLocaleString() }}</span>
                  </div>
                  <MarkdownRenderer
                    :markdown="comment.content"
                  />
                </v-card-text>
              </v-card>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col>
        <IssueDialogSidebar
          :comments="comments"
          :issue="currentIssue"
        />

        <v-card
          v-if="isLoggedIn"
          class="card__border bg-transparent mt-4"
          color="#19120D33"
          elevation="6"
          rounded="xl"
        >
          <v-card-text>
            <v-select
              v-if="currentIssue"
              v-model="currentIssue!.status"
              :items="[
                    { title: 'Backlog', value: 'backlog' },
                    { title: 'Planned', value: 'planned' },
                    { title: 'In Progress', value: 'in_progress' },
                    { title: 'Done', value: 'done' },
                    { title: 'Closed', value: 'closed' },

                  ]"
              class="mb-4"
              color="primary"
              density="compact"
              hide-details
              rounded="xl"
              @update:modelValue="statusChanged"
            />

            <v-btn
              :to="`/spaces/${space?.slug}/issues/${currentIssue?.id}/edit`"
              block
              class="mb-2"
              color="primary"
              variant="tonal"
            >
              Edit Issue
            </v-btn>

            <v-btn
              :to="`/spaces/${space?.slug}/issues`"
              block
              color="error"
              variant="tonal"
              @click="deleteIssueReq"
            >
              Delete Issue
            </v-btn>
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
</style>
