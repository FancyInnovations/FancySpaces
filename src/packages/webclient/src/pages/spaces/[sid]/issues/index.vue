<script lang="ts" setup>

import {type Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import SpaceSidebar from "@/components/SpaceSidebar.vue";
import {useHead} from "@vueuse/head";

const router = useRouter();

const space = ref<Space>();


onMounted(async () => {
  const spaceID = (useRoute().params as any).sid as string;
  space.value = await getSpace(spaceID);

  useHead({
    title: `${space.value.title} issues - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || 'View issues for this space on FancySpaces.'
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
        <div class="d-flex justify-space-between">
          <div class="d-flex flex-column justify-center">
            <v-img
              :href="`/spaces/${space?.slug}`"
              :src="space?.icon_url || '/logo.png'"
              alt="Space Icon"
              height="100"
              max-height="100"
              max-width="100"
              min-height="100"
              min-width="100"
              width="100"
            />
          </div>

          <div class="mx-4 d-flex flex-column justify-space-between flex-grow-1">
            <div>
              <h1>{{ space?.title }}</h1>
              <p class="text-body-1 mt-2">{{ space?.summary }}</p>
            </div>
          </div>

          <div class="d-flex flex-column justify-center">
            <v-btn
              :to="`/spaces/${space?.slug}/issues/new`"
              color="primary"
              size="large"
              variant="tonal"
            >
              New Issue
            </v-btn>
          </div>
        </div>

        <hr
          class="mt-4 grey-border-color"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col v-for="category in ['Planned', 'In Progress', 'Completed']" :key="category" cols="12" md="4">
        <v-card
          class="card__border"
          color="#29152550"
          elevation="12"
          rounded="xl"
        >
          <v-card-title class="my-2 ml-2">{{ category }}</v-card-title>

          <v-card-text>
            <v-card
              v-for="i in 6"
              :key="i"
              class="mb-4"
              color="#3E1A6D20"
              elevation="6"
              rounded="xl"
            >
              <v-card-text class="d-flex">
                <div class="d-flex flex-column justify-space-between">
                  <p class="text-body-1 mb-2">This is a placeholder issue card. Issue functionality is coming soon!</p>
                  <div class="d-inline-block">
                    <v-chip
                      class="mr-2"
                      color="primary"
                      prepend-icon="mdi-sign-text"
                      rounded
                      variant="tonal"
                    >
                      #7G5B1
                    </v-chip>

                    <v-chip
                      class="mr-2"
                      color="red"
                      prepend-icon="mdi-bug-outline"
                      rounded
                      variant="tonal"
                    >
                      Bug
                    </v-chip>

                    <v-chip
                      class="mr-2"
                      color="orange"
                      prepend-icon="mdi-alert-circle-outline"
                      rounded
                      variant="tonal"
                    >
                      Medium
                    </v-chip>
                  </div>
                </div>
                <div class="flex-grow-1 d-flex flex-column justify-space-between">
                  <v-btn
                    class="mb-1"
                    icon="mdi-transfer-right"
                    variant="text"
                  />

                  <v-btn
                    class="mt-1"
                    icon="mdi-transfer-left"
                    variant="text"
                  />
                </div>
              </v-card-text>
            </v-card>
          </v-card-text>
        </v-card>
      </v-col>

    </v-row>
  </v-container>

  <IssueDialog
    :comments="[
      {
        id: 'CMT123',
        issue: '7G5B1',
        author: 'user789',
        content: 'I have encountered this bug as well. It seems to occur when performing [specific action].',
        created_at: new Date(),
        updated_at: new Date()
      },
      {
        id: 'CMT124',
        issue: '7G5B1',
        author: 'user321',
        content: 'A temporary workaround is to [workaround details], but a permanent fix is needed.',
        created_at: new Date(),
        updated_at: new Date()
      },
      {
        id: 'CMT125',
        issue: '7G5B1',
        author: 'user654',
        content: 'The development team is actively investigating this issue and will provide updates as they become available.',
        created_at: new Date(2025, 0, 26, 10, 0, 0, 0),
        updated_at: new Date(2025, 0, 26, 10, 0, 0, 0)
      }
    ]"
    :issue="{
    id: '7G5B1',
    space: 'example-space',
    title: 'Sample Bug Report',
    description: '## 🐞 Bug Description\n'+
'A defect was identified that causes unexpected behavior in the application. Further investigation is required to determine the root cause and scope of impact.\n'+
'\n'+
'## 🔁 Steps to Reproduce\n'+
'1. Navigate to `[page / feature]`\n'+
'2. Perform `[action]`\n'+
'3. Observe the result\n'+
'\n'+
'## ✅ Expected Result\n'+
'The system should `[expected behavior]`.\n'+
'\n'+
'## ❌ Actual Result\n'+
'The system instead `[actual behavior]`.\n'+
'\n'+
'## 🌍 Environment\n'+
'- App version: `[version]`\n'+
'- Environment: `[dev / staging / prod]`\n'+
'- Browser / Device: `[if applicable]`\n'+
'\n'+
'## 📎 Additional Notes\n'+
'- Frequency: `[always / intermittent / once]`\n'+
'- Severity: `[low / medium / high / critical]`\n'+
'- Attachments: `[logs / screenshots / videos if any]`',
    type: 'story',
    priority: 'medium',
    status: 'todo',
    assignee: 'user123',
    reporter: 'user456',
    created_at: new Date(2026, 0, 26, 21, 0, 0, 0),
    updated_at: new Date(),
    external_source: 'github'
    }"
  />
</template>

<style scoped>
.grey-border-color {
  border-color: rgba(0, 0, 0, 0.8);
}
</style>
