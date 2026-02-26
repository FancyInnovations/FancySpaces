<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import SpaceHeader from "@/components/SpaceHeader.vue";
import {getBlogArticlesForSpace} from "@/api/blogs/blogs.ts";
import type {BlogArticle} from "@/api/blogs/types.ts";
import BlogList from "@/components/blog/BlogList.vue";
import {useUserStore} from "@/stores/user.ts";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const space = ref<Space>();
const articles = ref<BlogArticle[]>([]);

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

  articles.value = await getBlogArticlesForSpace(space.value.id);

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
              v-if="isMember"
              :to="`/spaces/${space?.slug}/blog/new`"
              color="primary"
              size="large"
              variant="tonal"
            >
              New Article
            </v-btn>
          </template>
        </SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="8">
        <h1 class="text-h4 text-center font-weight-bold mb-2">{{ space?.title }} Blog</h1>
        <p class="text-subtitle-1 text-center text-grey">{{ articles.length }} {{ articles.length === 1 ? 'article' : 'articles' }}</p>
      </v-col>
    </v-row>

    <BlogList
      :articles="articles"
      :space-slug="space?.slug"
      class="mt-8"
    />

  </v-container>
</template>

<style scoped>

</style>
