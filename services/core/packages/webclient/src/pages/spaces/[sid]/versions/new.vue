<script lang="ts" setup>
import type {Space} from "@/api/spaces/types";
import {getSpace} from "@/api/spaces/spaces";
import {useUserStore} from "@/stores/user";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import SpaceHeader from "@/components/SpaceHeader.vue";
import VersionForm from "@/components/versions/VersionForm.vue";
import {useHead} from "@vueuse/head";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const space = ref<Space>();

onMounted(async () => {
  try {
    const requestedID = (route.params as any).sid as string;
    space.value = await getSpace(requestedID);
    const userID = userStore.user?.id;
    if (!userStore.isAuthenticated || !userID || !(space.value.creator === userID || space.value.members.some(member => member.user_id === userID))) {
      await router.push(`/spaces/${space.value.slug}`);
      return;
    }
    useHead({ title: `New version - ${space.value.title} - FancySpaces` });
  } catch {
    await router.push("/");
  }
});
</script>

<template>
  <v-container width="90%">
    <v-row>
      <v-col class="flex-grow-0 pa-0"><SpaceSidebar :space="space" /></v-col>
      <v-col>
        <SpaceHeader :space="space">
          <template #quick-actions><v-btn :to="`/spaces/${space?.slug}/versions`" color="primary" size="large" variant="tonal">View Versions</v-btn></template>
        </SpaceHeader>
        <hr class="grey-border-color mt-4" />
      </v-col>
    </v-row>
    <VersionForm v-if="space" :spaceID="space.id" :spaceSlug="space.slug" class="mt-8" />
    <v-progress-circular v-else class="d-block mx-auto mt-12" color="primary" indeterminate />
  </v-container>
</template>
