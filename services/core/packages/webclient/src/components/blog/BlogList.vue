<script lang="ts" setup>

import type {BlogArticle} from "@/api/blogs/types.ts";

const props = defineProps<{
  articles: BlogArticle[];
  spaceSlug?: string;
}>();

const sortedArticles = computed(() => {
  return props.articles.sort((a, b) => b.published_at.getTime() - a.published_at.getTime());
});

</script>

<template>
  <p
    v-if="articles.length === 0"
    class="text-center mt-8"
  >
    No blog articles found.
  </p>
  <div v-else>
    <v-row
      v-for="article in sortedArticles"
      :key="article.id"
      justify="center"
    >
      <v-col md="5">
        <Card
          :to="article.space_id ? `/spaces/${spaceSlug ? spaceSlug : article.space_id}/blog/${article.id}` : `/users/${article.author}/blog/${article.id}`"
          class="hoverable"
        >
          <v-card-title class="mt-2">{{ article.title }}</v-card-title>
          <v-card-subtitle>
            Published: {{ article.published_at.toLocaleString() }}
          </v-card-subtitle>
          <v-card-text>{{ article.summary }}</v-card-text>
        </Card>
      </v-col>
    </v-row>
  </div>
</template>

<style scoped>

</style>
