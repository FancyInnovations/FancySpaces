<script lang="ts" setup>

import {useUserStore} from "@/stores/user.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import {createBlogArticle} from "@/api/blogs/blogs.ts";

const router = useRouter();
const userStore = useUserStore();
const notificationStore = useNotificationStore();

const props = defineProps<{
  spaceID?: string;
}>();

const title = ref("");
const summary = ref("");
const content = ref("");

async function publishArticle() {
  if (!(await userStore.isAuthenticated)) {
    return;
  }

  if (title.value.trim() === "" || content.value.trim() === "" || summary.value.trim() === "") {
    notificationStore.error("Please fill in all fields before publishing.");
    return;
  }

  await createBlogArticle(props.spaceID!, title.value, summary.value, content.value);

  notificationStore.info("Article published successfully!");

  await router.push(props.spaceID ? `/spaces/${props.spaceID}/blog` : `/users/${userStore.user?.name}/blog`);
}

</script>

<template>
  <v-row justify="center">
    <v-col md="8">
      <Card>
        <v-card-title class="mt-2">
          New Blog Article
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
            @click="publishArticle()"
          >
            Publish
          </v-btn>
        </v-card-text>
      </Card>
    </v-col>
  </v-row>
</template>

<style scoped>

</style>
