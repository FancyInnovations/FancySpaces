<script lang="ts" setup>

import {createSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import {useNotificationStore} from "@/stores/notifications.ts";
import {useUserStore} from "@/stores/user.ts";
import {mapCategoryToDisplayname} from "@/api/spaces/types.ts";

const router = useRouter();
const route = useRoute();
const notificationStore = useNotificationStore();
const userStore = useUserStore();

const isLoggedIn = ref(false);

const name = ref('');
const slug = ref('');
const categories = ref([]);

const possibleCategories = [
  "minecraft_plugin",
  "minecraft_server",
  "minecraft_mod",
  "hytale_plugin",
  "web_app",
  "mobile_app",
  "other"
];

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

const isEverythingValid = computed(() => {
  return nameRule(name.value) === true && slugRule(slug.value) === true;
});

onMounted(async () => {
  isLoggedIn.value = await userStore.isAuthenticated;

  useHead({
    title: `FancySpaces`,
    meta: [
      {
        name: 'description',
        content: 'Create a new space on FancySpaces.'
      }
    ]
  });
});

async function createNewSpace() {
  await createSpace(slug.value, name.value, "", "", categories.value, "");

  notificationStore.info("Space created successfully!");

  await router.push(`/spaces/${slug}`);
}

</script>

<template>
  <v-container width="100%">
    <v-row>
      <v-col>
        <h1 class="text-center">Create new space</h1>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="4">
        <v-text-field
          v-model="name"
          :rules="[nameRule]"
          color="primary"
          label="Space Name"
          required
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="4">
        <v-text-field
          v-model="slug"
          :rules="[slugRule]"
          color="primary"
          label="Space Slug"
          required
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="4">
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

    <v-row class="mt-8" justify="center">
      <v-col md="4">
        <p class="text-grey">You can customize your space later to add a description, logo and more!</p>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="4">
        <v-btn
          :disabled="!isEverythingValid"
          color="primary"
          size="large"
          @click="createNewSpace()"
        >
          Create Space
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
