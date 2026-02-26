<script lang="ts" setup>

import {useUserStore} from "@/stores/user.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import {deleteBlogArticle, updateBlogArticle} from "@/api/blogs/blogs.ts";
import type {BlogArticle} from "@/api/blogs/types.ts";
import {useRouter} from "vue-router";
import {useConfirmationStore} from "@/stores/confirmation.ts";

const router = useRouter();
const userStore = useUserStore();
const notificationStore = useNotificationStore();
const confirmationStore = useConfirmationStore();

const props = defineProps<{
  article: BlogArticle;
  content: string;
}>();

const title = ref("");
const summary = ref("");
const content = ref("");

watch(() => props.article, (newArticle) => {
  if (newArticle) {
    title.value = newArticle.title;
    summary.value = newArticle.summary;
  }
}, { immediate: true });

watch(() => props.content, (newContent) => {
  if (newContent) {
    content.value = newContent;
  }
}, { immediate: true });

async function editArticle() {
  if (!(await userStore.isAuthenticated)) {
    return;
  }

  if (title.value.trim() === "" || content.value.trim() === "" || summary.value.trim() === "") {
    notificationStore.error("Please fill in all fields before publishing.");
    return;
  }

  await updateBlogArticle(
    props.article.id,
    title.value !== props.article.title ? title.value : "",
    summary.value !== props.article.summary ? summary.value : "",
    content.value !== props.content ? content.value : ""
  );

  notificationStore.info("Article updated successfully!");

  await router.push(props.article.space_id ? `/spaces/${props.article.space_id}/blog` : `/users/${userStore.user?.name}/blog`);
}

async function deleteArticleReq() {
  if (!(await userStore.isAuthenticated)) {
    return;
  }

  confirmationStore.confirmation = {
    shown: true,
    persistent: true,
    title: "Delete Blog Article",
    text: "Are you sure you want to delete this blog article? This action cannot be undone.",
    yesText: "Delete",
    onConfirm: async () => {
      await deleteBlogArticle(props.article.id);

      notificationStore.info("Article deleted successfully!");

      await router.push(props.article.space_id ? `/spaces/${props.article.space_id}/blog` : `/users/${userStore.user?.name}/blog`);
    }
  };
}

</script>

<template>
  <v-row justify="center">
    <v-col md="8">
      <Card>
        <v-card-title class="mt-2">
          Edit Blog Article
        </v-card-title>

        <v-card-text>
          <v-text-field
            v-model="title"
            color="primary"
            label="Title"
          />

          <v-textarea
            v-model="summary"
            color="primary"
            label="Summary"
            rows="2"
          />

          <v-textarea
            v-model="content"
            color="primary"
            label="Content (Markdown)"
            rows="20"
          />

          <v-btn
            :disabled="title.trim() === '' || summary.trim() === '' || content.trim() === ''"
            class="my-2"
            color="primary"
            @click="editArticle()"
          >
            Edit
          </v-btn>

          <v-btn
            class="my-2 ml-4"
            color="error"
            variant="outlined"
            @click="deleteArticleReq()"
          >
            Delete
          </v-btn>
        </v-card-text>
      </Card>
    </v-col>
  </v-row>
</template>

<style scoped>

</style>
