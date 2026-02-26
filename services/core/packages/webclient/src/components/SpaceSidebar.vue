<script lang="ts" setup>

import {mapCategoryToDisplayname, mapLinkToDisplayname, mapLinkToIcon, type Space} from "@/api/spaces/types.ts";
import {useUserStore} from "@/stores/user.ts";

const userStore = useUserStore();

const props = defineProps<{
  space?: Space
}>();

const isMember = computed(() => {
  if (!props.space) return false;
  if (!userStore.isAuthenticated) return false;

  const userID =  userStore.user?.id;
  return props.space.creator == userID || props.space.members.some(member => member.user_id === userID);
});

</script>

<template>
  <v-navigation-drawer
    class="sidebar__mobile sidebar__background ma-4"
    elevation="12"
    rounded="xl"
  >
    <v-list>
      <v-list-item>
        <v-list-item-title class="text-h6 font-weight-bold">{{ space?.title }}</v-list-item-title>
        <v-list-item-subtitle>{{ mapCategoryToDisplayname(space?.categories[0]) }}</v-list-item-subtitle>
      </v-list-item>

      <v-divider class="mt-2"/>

      <v-list-item
        :to="`/spaces/${space?.slug}`"
        exact
        link
        prepend-icon="mdi-information-slab-circle-outline"
        title="Information"
      />

      <v-list-item
        v-if="space?.blog_settings.enabled"
        :to="`/spaces/${space?.slug}/blog`"
        link
        prepend-icon="mdi-notebook-outline"
        title="Blog"
      />

<!--      <v-list-item-->
<!--        :to="`/spaces/${space?.slug}/docs`"-->
<!--        link-->
<!--        prepend-icon="mdi-book-open-variant-outline"-->
<!--        title="Documentation"-->
<!--      />-->

<!--      <v-list-item-->
<!--        :to="`/spaces/${space?.slug}/source`"-->
<!--        link-->
<!--        prepend-icon="mdi-source-branch"-->
<!--        title="Source Code"-->
<!--      />-->

      <v-list-item
        v-if="space?.release_settings.enabled"
        :to="`/spaces/${space?.slug}/versions`"
        link
        prepend-icon="mdi-file-download-outline"
        title="Downloads"
      />

      <v-list-item
        v-if="space?.maven_repository_settings.enabled"
        :to="`/spaces/${space?.slug}/maven-repos`"
        link
        prepend-icon="mdi-database-outline"
        title="Maven Repository"
      />

      <v-list-item
        v-if="space?.maven_repository_settings.enabled"
        :to="`/spaces/${space?.slug}/javadoc`"
        link
        prepend-icon="mdi-book-outline"
        title="JavaDoc"
      />

      <v-list-item
        v-if="space?.issue_settings.enabled"
        :to="`/spaces/${space?.slug}/issues`"
        link
        prepend-icon="mdi-format-list-checks"
        title="Issues"
      />

<!--      <v-list-item-->
<!--        :to="`/spaces/${space?.slug}/support-tickets`"-->
<!--        link-->
<!--        prepend-icon="mdi-bug-outline"-->
<!--        title="Support Tickets"-->
<!--      />-->

<!--      <v-list-item-->
<!--        :to="`/spaces/${space?.slug}/roadmap`"-->
<!--        link-->
<!--        prepend-icon="mdi-road-variant"-->
<!--        title="Roadmap"-->
<!--      />-->

      <v-list-item
        v-if="space?.storage_settings.enabled && isMember"
        :to="`/spaces/${space?.slug}/storage`"
        link
        prepend-icon="mdi-library-shelves"
        title="Storage"
      />

      <v-list-item
        v-if="space?.analytics_settings.enabled"
        :to="`/analytics/${space?.slug}`"
        link
        prepend-icon="mdi-chart-box-outline"
        title="Analytics"
      />

      <v-list-item
        v-if="space?.secrets_settings.enabled && isMember"
        :to="`/spaces/${space?.slug}/secrets`"
        link
        prepend-icon="mdi-shield-key-outline"
        title="Secrets"
      />

      <v-divider v-if="space && space.links && space.links.length > 0"/>
      <v-list-subheader v-if="space && space.links && space.links.length > 0">External Links</v-list-subheader>

      <v-list-item
        v-for="(link) in space?.links" :key="link.name"
        :href="link.url"
        :prepend-icon="mapLinkToIcon(link.name)"
        :title="mapLinkToDisplayname(link.name)"
        link
        target="_blank"
      />

      <v-divider v-if="isMember" />
      <v-list-subheader v-if="isMember">Settings</v-list-subheader>

      <v-list-item
        v-if="isMember"
        :to="`/spaces/${space?.slug}/settings`"
        link
        prepend-icon="mdi-cog-outline"
        title="Settings"
      />

      <v-list-item
        v-if="isMember"
        :to="`/spaces/${space?.slug}/features`"
        link
        prepend-icon="mdi-star-outline"
        title="Features"
      />

      <v-list-item
        v-if="isMember"
        :to="`/spaces/${space?.slug}/members`"
        link
        prepend-icon="mdi-account-multiple-outline"
        title="Members"
      />

    </v-list>
  </v-navigation-drawer>
</template>

<style scoped>
.sidebar__background {
  max-height: calc(100vh - 96px);
  background-color: #19120D33 !important;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

@media (max-width: 960px) {
  .sidebar__mobile {
    display: none;
  }
}
</style>
