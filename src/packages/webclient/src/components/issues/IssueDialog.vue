<script lang="ts" setup>

import Dialog from "@/components/common/Dialog.vue";
import type {Issue, IssueComment} from "@/api/issues/types.ts";
import IssueDialogSidebar from "@/components/issues/IssueDialogSidebar.vue";

const props = defineProps<{
  issue: Issue,
  comments: IssueComment[]
}>();

const showDialog = ref(true);

</script>

<template>
  <Dialog
    :shown="showDialog"
    width="60%"
  >
    <div class="rounded-xl">
      <div class="issue-dialog d-flex">
        <IssueDialogSidebar
          :comments="comments"
          :issue="issue"
        />

        <div class="pr-4 flex-grow-1">
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
                :markdown="props.issue.description"
                class="px-2 issue-description"
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
              Comments ({{ props.comments.length }})
            </v-card-title>

            <v-card-text>
              <p v-if="props.comments.length === 0">
                No comments yet.
              </p>
              <div v-else class="issue-comments">
                <v-card
                  v-for="comment in props.comments"
                  :key="comment.id"
                  class="card__border bg-transparent mb-3"
                  color="#29152550"
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
          >
            Copy Link
          </v-btn>

        <v-btn
          variant="text"
          @click="showDialog = false"
        >
          Close
        </v-btn>
      </div>
    </div>
  </Dialog>
</template>

<style scoped>
.issue-dialog {
  max-height: 85vh;
  overflow-y: auto;
}

.issue-description {
  max-height: 500px;
  overflow-y: auto;
}
</style>
