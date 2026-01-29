<script lang="ts" setup>

import Dialog from "@/components/common/Dialog.vue";
import IssueDialogSidebar from "@/components/issues/IssueDialogSidebar.vue";
import {useIssueDialogStore} from "@/stores/issue-dialog.ts";
import type {IssueComment} from "@/api/issues/types.ts";

const issueDialogStore = useIssueDialogStore();

const comments = computed<IssueComment[]>(() => {
  return [
    {
      id: 'CMT123',
      issue: '7G5B1',
      author: 'user789',
      content: 'I have encountered this bug as well. It seems to occur when performing [specific action].',
      created_at: new Date(),
      updated_at: new Date()
    },
    {
      id: 'CMT124',
      issue: '7G5B1',
      author: 'user321',
      content: 'A temporary workaround is to [workaround details], but a permanent fix is needed.',
      created_at: new Date(),
      updated_at: new Date()
    },
    {
      id: 'CMT125',
      issue: '7G5B1',
      author: 'user654',
      content: 'The development team is actively investigating this issue and will provide updates as they become available.',
      created_at: new Date(2025, 0, 26, 10, 0, 0, 0),
      updated_at: new Date(2025, 0, 26, 10, 0, 0, 0)
    }
  ];
});

function copyLink() {
  const issueLink = `${window.location.origin}/spaces/${issueDialogStore.issue?.space}/issues/${issueDialogStore.issue?.id}`;
  navigator.clipboard.writeText(issueLink);
}

</script>

<template>
  <Dialog
    :shown="issueDialogStore.isOpen"
    width="64%"
  >
    <div class="rounded-xl">
      <div class="py-2 border-b">
        <h1 class="text-center text-h4 text-secondary">{{ issueDialogStore.issue?.title }}</h1>
      </div>

      <div class="issue-dialog-inner d-flex">
        <IssueDialogSidebar
          :comments="comments"
          :issue="issueDialogStore.issue!"
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
            class="mr-4"
            variant="text"
            @click="copyLink"
          >
            Copy Link
          </v-btn>

        <v-btn
          variant="text"
          @click="() => issueDialogStore.close()"
        >
          Close
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
