<script lang="ts" setup>

import Dialog from "@/components/common/Dialog.vue";
import IssueDialogSidebar from "@/components/issues/IssueDialogSidebar.vue";
import {useIssueDialogStore} from "@/stores/issue-dialog.ts";
import type {IssueComment} from "@/api/issues/types.ts";
import {updateIssue} from "@/api/issues/issues.ts";

const issueDialogStore = useIssueDialogStore();
const isLoggedIn = computed(() => {
  return localStorage.getItem("fs_api_key") !== null;
});

const comments = computed<IssueComment[]>(() => {
  // return [
  //   {
  //     id: 'CMT123',
  //     issue: '7G5B1',
  //     author: 'user789',
  //     content: 'I have encountered this bug as well. It seems to occur when performing [specific action].',
  //     created_at: new Date(),
  //     updated_at: new Date()
  //   },
  //   {
  //     id: 'CMT124',
  //     issue: '7G5B1',
  //     author: 'user321',
  //     content: 'A temporary workaround is to [workaround details], but a permanent fix is needed.',
  //     created_at: new Date(),
  //     updated_at: new Date()
  //   },
  //   {
  //     id: 'CMT125',
  //     issue: '7G5B1',
  //     author: 'user654',
  //     content: 'The development team is actively investigating this issue and will provide updates as they become available.',
  //     created_at: new Date(2025, 0, 26, 10, 0, 0, 0),
  //     updated_at: new Date(2025, 0, 26, 10, 0, 0, 0)
  //   },
  //   {
  //     id: 'CMT125',
  //     issue: '7G5B1',
  //     author: 'user654',
  //     content: 'The development team is actively investigating this issue and will provide updates as they become available.',
  //     created_at: new Date(2025, 0, 26, 10, 0, 0, 0),
  //     updated_at: new Date(2025, 0, 26, 10, 0, 0, 0)
  //   },
  //   {
  //     id: 'CMT125',
  //     issue: '7G5B1',
  //     author: 'user654',
  //     content: 'The development team is actively investigating this issue and will provide updates as they become available.',
  //     created_at: new Date(2025, 0, 26, 10, 0, 0, 0),
  //     updated_at: new Date(2025, 0, 26, 10, 0, 0, 0)
  //   }
  // ];

  return [];
});

function copyLink() {
  const issueLink = `${window.location.origin}/spaces/${issueDialogStore.issue?.space}/issues/${issueDialogStore.issue?.id}`;
  navigator.clipboard.writeText(issueLink);

  window.alert('Issue link copied to clipboard!');
}

function copyID() {
  const issueID = issueDialogStore.issue?.id;
  if (issueID) {
    navigator.clipboard.writeText(issueID);
    window.alert('Issue ID copied to clipboard!');
  }
}

async function statusChanged(newStatus: string) {
  if (!issueDialogStore.issue) return;

  const newIssue = { ...issueDialogStore.issue! };
  newIssue.status = newStatus as any;
  await updateIssue(newIssue.space, newIssue.id, newIssue);
}

</script>

<template>
  <Dialog
    :shown="issueDialogStore.isOpen"
    width="64%"
  >
    <div class="rounded-xl">

      <div class="py-2 border-b d-flex align-center px-4">
        <h1 class="ml-2 text-h4 text-secondary">{{ issueDialogStore.issue?.title }}</h1>

        <div class="flex-grow-1 d-flex justify-end align-center">
          <v-select
            v-if="issueDialogStore.issue && isLoggedIn"
            v-model="issueDialogStore.issue!.status"
            :items="[
                    { title: 'Backlog', value: 'backlog' },
                    { title: 'Planned', value: 'planned' },
                    { title: 'In Progress', value: 'in_progress' },
                    { title: 'Done', value: 'done' },
                    { title: 'Closed', value: 'closed' },

                  ]"
            class="mr-4"
            color="primary"
            density="compact"
            hide-details
            max-width="200"
            min-width="200"
            variant="solo"
            @update:modelValue="statusChanged"
          />

          <v-btn
            :href="`/spaces/${issueDialogStore.issue?.space}/issues/${issueDialogStore.issue?.id}`"
            class="mr-2"
            color="secondary"
            icon="mdi-open-in-new"
            target="_blank"
            variant="text"
          />

          <v-btn
            color="secondary"
            icon="mdi-close"
            variant="text"
            @click="() => issueDialogStore.close()"
          />
        </div>
      </div>

      <div class="issue-dialog-inner d-flex">
        <IssueDialogSidebar
          :comments="comments"
          :issue="issueDialogStore.issue!"
          class="ma-4"
        />

        <div class="issue-dialog-inner pr-4 flex-grow-1">
          <v-card
            class="card__border mt-4 bg-transparent"
            color="#150D1950"
            elevation="12"
            min-width="600"
            rounded="xl"
          >
            <v-card-title class="mt-2">
              Description
            </v-card-title>

            <v-card-text>
              <MarkdownRenderer
                :markdown="issueDialogStore.issue?.description"
                class="issue-description"
              />
            </v-card-text>
          </v-card>

          <v-card
            class="card__border my-4 bg-transparent"
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
        </div>
      </div>

      <div class="d-flex justify-end pa-2 border-t">
        <v-btn
          class="mr-2"
          variant="text"
          @click="copyLink"
        >
          Copy Link
        </v-btn>

        <v-btn
          class="mr-2"
          variant="text"
          @click="copyID"
        >
          Copy ID
        </v-btn>

        <v-btn
          v-if="isLoggedIn"
          :to="`/spaces/${issueDialogStore.issue?.space}/issues/${issueDialogStore.issue?.id}/edit`"
          class="mr-2"
          variant="text"
          @click="issueDialogStore.close()"
        >
          Edit
        </v-btn>
      </div>
    </div>
  </Dialog>
</template>

<style scoped>
.issue-dialog-inner {
  max-height: 75vh;
  overflow-y: auto;
  scrollbar-width: none;
}

.issue-description {
  max-height: 500px;
  overflow-y: auto;
  scrollbar-width: none;
}
</style>
