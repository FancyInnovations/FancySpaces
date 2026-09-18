<script lang="ts" setup>
import type {Space} from "@/api/spaces/types";
import type {SpaceVersion} from "@/api/versions/types";
import {getSpace} from "@/api/spaces/spaces";
import {getVersion} from "@/api/versions/versions";
import {useUserStore} from "@/stores/user";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import SpaceHeader from "@/components/SpaceHeader.vue";
import VersionForm from "@/components/versions/VersionForm.vue";
import {useHead} from "@vueuse/head";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const space = ref<Space>();
const version = ref<SpaceVersion>();
const loading = ref(true);

onMounted(async () => {
  try {
    const requestedID = (route.params as any).sid as string;
    space.value = await getSpace(requestedID);
    const userID = userStore.user?.id;
    if (!userStore.isAuthenticated || !userID || !(space.value.creator === userID || space.value.members.some(member => member.user_id === userID))) {
      await router.push(`/spaces/${space.value.slug}`);
      return;
    }
    version.value = await getVersion(space.value.id, (route.params as any).vid as string);
    useHead({ title: `Edit ${version.value.name} - ${space.value.title} - FancySpaces` });
  } catch {
    await router.push(`/spaces/${space.value?.slug || (route.params as any).sid}/versions`);
  } finally {
    loading.value = false;
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
    <v-progress-circular v-if="loading" class="d-block mx-auto mt-12" color="primary" indeterminate />
    <VersionForm v-else-if="space && version" :spaceID="space.id" :spaceSlug="space.slug" :version="version" class="mt-8" />
  </v-container>
</template>
