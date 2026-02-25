<script lang="ts" setup>

import {getSpace, updateSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import {useNotificationStore} from "@/stores/notifications.ts";
import {useUserStore} from "@/stores/user.ts";
import {mapCategoryToDisplayname, type Space} from "@/api/spaces/types.ts";
import SpaceHeader from "@/components/SpaceHeader.vue";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import Card from "@/components/common/Card.vue";

const router = useRouter();
const route = useRoute();
const notificationStore = useNotificationStore();
const userStore = useUserStore();

const isLoggedIn = ref(false);

const name = ref('');
const slug = ref('');
const summary = ref('');
const description = ref('');
const iconURL = ref('');
const categories = ref<string[]>([]);
const possibleCategories = [
  "minecraft_plugin",
  "minecraft_server",
  "minecraft_mod",
  "hytale_plugin",
  "web_app",
  "mobile_app",
  "other"
];

const showDeleteDialog = ref(false);

const space = ref<Space>();

const nameRule = (value: string) => {
  if (!value) return 'Space name is required';

  if (value.length < 3) return 'Space name must be at least 3 characters long';

  if (value.length > 100) return 'Space name must not exceed 100 characters';

  return true;
};

const slugRule = (value: string) => {
  if (!value) return 'Space slug is required';

  if (value.length < 3) return 'Space slug must be at least 3 characters long';

  if (value.length > 20) return 'Space slug must not exceed 20 characters';

  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(value)) return 'Space slug can only contain lowercase letters, numbers and hyphens, and must start and end with a letter or number';

  return true;
};

const summaryRule = (value: string) => {
  if (value.length > 300) return 'Summary must not exceed 300 characters';

  return true;
};

const descriptionRule = (value: string) => {
  if (value.length > 10000) return 'Description must not exceed 10000 characters';

  return true;
};

const iconURLRule = (value: string) => {
  if (value.length > 0 && !/^https?:\/\/.+\.(jpg|jpeg|png|gif|svg)$/.test(value)) return 'Icon URL must be a valid URL pointing to an image (jpg, jpeg, png, gif, svg)';

  return true;
};

const isEverythingValid = computed(() => {
  return slugRule(slug.value) === true && nameRule(name.value) === true && summaryRule(summary.value) === true && descriptionRule(description.value) === true && iconURLRule(iconURL.value) === true;
});

onMounted(async () => {
  isLoggedIn.value = await userStore.isAuthenticated;

  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!isLoggedIn) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  name.value = space.value.title;
  slug.value = space.value.slug;
  summary.value = space.value.summary || '';
  description.value = space.value.description || '';
  iconURL.value = space.value.icon_url || '';
  categories.value = space.value.categories || [];

  useHead({
    title: `${space.value.title} settings - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: 'Customize the settings of your space, including name, slug, summary, description and more.'
      }
    ]
  });
});

async function updateSpaceReq() {
  if (!isEverythingValid.value) {
    notificationStore.error("Please fix the errors in the form before updating the space.");
    return;
  }

  await updateSpace(space.value!.id, slug.value, name.value, summary.value, description.value, categories.value, iconURL.value);

  notificationStore.info("Space created successfully!");

  await router.push(`/spaces/${slug.value}`);
}

async function deleteSpaceReq() {

}

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
        <SpaceHeader :space="space"/>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <h1 class="text-center">Space settings</h1>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="3">
        <v-text-field
          v-model="name"
          :rules="[nameRule]"
          color="primary"
          label="Name"
          required
        />
      </v-col>

      <v-col md="3">
        <v-text-field
          v-model="slug"
          :rules="[slugRule]"
          color="primary"
          label="Slug"
          required
        />
      </v-col>
    </v-row>

    <v-row justify="center">
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <p class="text-caption">Categories (you can select multiple):</p>

        <v-chip-group
          v-model="categories"
          column
          multiple
        >
          <v-chip
            v-for="category in possibleCategories"
            :key="category"
            :value="category"
            color="primary"
          >
            {{ mapCategoryToDisplayname(category) }}
          </v-chip>
        </v-chip-group>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-text-field
          v-model="iconURL"
          :rules="[iconURLRule]"
          color="primary"
          label="Logo URL"
          required
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="6">
        <v-text-field
          v-model="summary"
          :rules="[summaryRule]"
          color="primary"
          label="Summary"
          required
        />
      </v-col>
    </v-row>

    <v-row class="mt-8" justify="center">
      <v-col md="5">
        <h2 class="text-h6 text-center">Description (Markdown supported)</h2>
      </v-col>

      <v-col md="5">
        <h2 class="text-h6 text-center">Live Preview</h2>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="5">
        <v-textarea
          v-model="description"
          :rules="[descriptionRule]"
          color="primary"
          hide-details
          no-resize
          outlined
          rounded="xl"
          rows="28"
        />
      </v-col>

      <v-col md="5">
        <Card
          class="mb-4 overflow-y-auto"
          height="100%"
          max-height="700"
        >
          <v-card-text>
            <MarkdownRenderer
              :markdown="description"
            />
          </v-card-text>
        </Card>
      </v-col>
    </v-row>

    <v-row class="mt-8" justify="center">
      <v-col md="10">
        <v-btn
          :disabled="!isEverythingValid"
          color="primary"
          @click="updateSpaceReq()"
        >
          Save Changes
        </v-btn>

        <v-btn
          class="ml-4"
          color="error"
          variant="outlined"
          @click="showDeleteDialog = true"
        >
          Delete Space
        </v-btn>
      </v-col>
    </v-row>
  </v-container>


  <Dialog :shown="showDeleteDialog">
    <v-card
      max-width="500"
    >
      <v-card-title class="text-h6">Delete space</v-card-title>

      <v-card-text>
        To delete this space, please reach out to our support team (via E-Mail or our Discord server) and provide them with the name and slug of the space you want to delete. We will then verify your ownership of the space and proceed with the deletion.
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>

        <v-btn
          text
          @click="showDeleteDialog = false"
        >
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </Dialog>
</template>

<style scoped>

</style>
