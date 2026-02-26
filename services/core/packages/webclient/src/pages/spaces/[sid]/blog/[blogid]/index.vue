<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import SpaceHeader from "@/components/SpaceHeader.vue";
import type {BlogArticle} from "@/api/blogs/types.ts";
import BlogArticleView from "@/components/blog/BlogArticleView.vue";
import {getBlogArticle, getBlogArticleContent} from "@/api/blogs/blogs.ts";
import {useUserStore} from "@/stores/user.ts";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const space = ref<Space>();
const article = ref<BlogArticle>();
const content = ref<string>('');

const isMember = computed(() => {
  if (!space.value) return false;
  if (!userStore.user) return false;

  const userID =  userStore.user?.id;
  return space.value.creator == userID || space.value.members.some(member => member.user_id === userID);
});

onMounted(async () => {
  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.blog_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  const blogID = (route.params as any).blogid as string;
  article.value = await getBlogArticle(blogID);
  content.value = await getBlogArticleContent(blogID);

  useHead({
    title: `${space.value.title} Blog - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || `Explore the ${space.value.title} project space on FancySpaces.`
      }
    ]
  });
});
</script>

<template>
  <v-container width="90%">
    <v-row>
      <v-col class="flex-grow-0 pa-0">
        <SpaceSidebar
          :space="space"
        />
      </v-col>

      <v-col>
        <SpaceHeader :space="space">
          <template #quick-actions>
            <v-btn
              :to="`/spaces/${space?.slug}/blog`"
              color="primary"
              exact
              size="large"
              variant="tonal"
            >
              Back to Blog
            </v-btn>

            <v-btn
              v-if="isMember"
              :to="`/spaces/${space?.slug}/blog/${article?.id}/edit`"
              class="mt-2"
              color="primary"
              size="large"
              variant="tonal"
            >
              Edit Article
            </v-btn>
          </template>
        </SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>

    <BlogArticleView
      v-if="article"
      :article="article"
      :content="content"
      class="mt-8"
    />
  </v-container>
</template>

<style scoped>

</style>
