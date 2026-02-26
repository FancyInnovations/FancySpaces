<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";
import SpaceHeader from "@/components/SpaceHeader.vue";
import BlogNew from "@/components/blog/BlogNew.vue";

const router = useRouter();
const route = useRoute();

const space = ref<Space>();

onMounted(async () => {
  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.blog_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

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
        <SpaceHeader :space="space"></SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>

    <BlogNew
      :spaceID="space?.id"
      class="mt-8"
    />

  </v-container>
</template>

<style scoped>

</style>
