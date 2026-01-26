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
      <div class="d-flex">
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
              Comments
            </v-card-title>

            <v-card-text>
              <p v-if="props.comments.length === 0">
                No comments yet.
              </p>
              <div v-else>
                <div
                  v-for="comment in props.comments"
                  :key="comment.id"
                  class="mb-3"
                >
                  <strong>{{ comment.author }}:</strong>
                  <p>{{ comment.content }}</p>
                </div>
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
.issue-description {
  max-height: 500px;
  overflow-y: auto;
}
</style>
