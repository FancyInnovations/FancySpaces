<script lang="ts" setup>

import type {BlogArticle} from "@/api/blogs/types.ts";

const props = defineProps<{
  article: BlogArticle;
  content: string;
}>();

</script>

<template>
  <v-row justify="center">
    <v-col md="8">
      <Card>
        <v-card-title class="mt-2">{{ article.title }}</v-card-title>
        <v-card-subtitle>{{ article.published_at.toLocaleString() }}</v-card-subtitle>

        <v-card-text>
          {{ article.summary }}
        </v-card-text>
      </Card>
    </v-col>
  </v-row>

  <v-row justify="center">
    <v-col md="8">
      <Card>
        <v-card-text class="py-0">
          <MarkdownRenderer
            :markdown="content"
          />
        </v-card-text>
      </Card>
    </v-col>
  </v-row>

</template>

<style scoped>

</style>
