<script lang="ts" setup>

import {mapCategoryToDisplayname, type Space} from "@/api/spaces/types.ts";
import {getLatestVersion} from "@/api/versions/versions.ts";
import type {SpaceVersion} from "@/api/versions/types.ts";
import {getDownloadCountForSpace} from "@/api/spaces/spaces.ts";
import Card from "@/components/common/Card.vue";
import type {User} from "@/api/auth/types.ts";
import {getPublicUser} from "@/api/auth/users.ts";

const props = defineProps<{
  space?: Space
  withBadge?: boolean
  withAuthor?: boolean
  soon?: boolean
}>();

const latestVersion = ref<SpaceVersion>();
const downloadCount = ref<number>(0);

const creator = ref<User>();
watch(() => props.space, async (newSpace) => {
  if (newSpace) {
    creator.value = await getPublicUser(newSpace.creator)
  }
}, {immediate: true});

onMounted(async () => {
    if (props.space) {
      latestVersion.value = await getLatestVersion(props.space.id);
      downloadCount.value = await getDownloadCountForSpace(props.space.id);
    }
})

</script>

<template>
  <Card
    :disabled="soon"
    min-width="600"
  >

    <v-card-text>
      <div class="d-flex justify-space-between flex-1-100-1">
        <div class="d-flex flex-column justify-center">
          <RouterLink :to="`/spaces/${space?.slug}`">
            <v-img
              :href="`/spaces/${space?.slug}`"
              :src="space?.icon_url || '/na-logo.png'"
              alt="Space Icon"
              height="100"
              max-height="100"
              max-width="100"
              min-height="100"
              min-width="100"
              width="100"
            />
          </RouterLink>
        </div>

        <div class="mx-4 d-flex flex-column justify-space-between flex-grow-1">
          <div>
            <RouterLink
              :to="`/spaces/${space?.slug}`"
            >
              <v-badge
                v-if="withBadge"
                :content="mapCategoryToDisplayname(space?.categories[0])"
                location="right center"
                offset-x="-20"
              >
                <h1>{{ space?.title }}</h1>
              </v-badge>

              <h1 v-else>{{ space?.title }}</h1>
            </RouterLink>

            <p class="text-body-1 mt-2">{{ space?.summary }}</p>
          </div>

          <div class="d-flex justify-space-between mt-2 text-grey-lighten-1">
            <p v-if="withAuthor" class="text-body-2">By {{ space?.creator }}</p>
            <p v-if="creator" class="text-body-2 link--hover"><RouterLink :to="'/users/'+creator.name">By: {{ creator?.name }}</RouterLink></p>
            <p class="text-body-2">Created {{ space?.created_at.toLocaleDateString() }}</p>
            <p class="text-body-2">Updated {{ latestVersion?.published_at.toLocaleDateString() || space?.created_at.toLocaleDateString() }}</p>
            <p class="text-body-2">{{ downloadCount }} downloads</p>
          </div>
        </div>

        <div class="d-flex flex-column justify-center">
          <v-btn
            v-if="!soon && space?.release_settings.enabled"
            :to="`/spaces/${space?.slug}/versions`"
            color="primary"
            icon="mdi-download"
            variant="tonal"
          />
        </div>
      </div>
    </v-card-text>
  </Card>
</template>

<style scoped>

</style>
