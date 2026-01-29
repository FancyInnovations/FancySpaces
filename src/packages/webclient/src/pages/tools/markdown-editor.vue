<script lang="ts" setup>

const rawMarkdown = ref<string>("# Welcome to the Markdown Editor!\n\nType your **Markdown** text in the left pane, and see the _preview_ on the right.\n\n- Easy to use\n- Real-time preview\n- Supports standard Markdown syntax\n\nHappy writing!");

</script>

<template>
  <v-container width="70%">
    <v-row class="mt-4">
      <v-col md="6" >
        <h1 class="text-h4 text-center">Markdown Editor</h1>
      </v-col>
      <v-col md="6" >
        <h1 class="text-h4 text-center">Live Preview</h1>
      </v-col>
    </v-row>

    <v-row >
      <v-col md="6">
        <v-textarea
          v-model="rawMarkdown"
          color="primary"
          hide-details
          no-resize
          outlined
          rounded="xl"
          rows="25"
        />
      </v-col>

      <v-col md="6">
        <v-card
          class="mb-4 card__border"
          color="#29231550"
          elevation="12"
          height="100%"
          rounded="xl"
        >
          <v-card-text>
            <MarkdownRenderer
              :markdown="rawMarkdown"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>

</template>

<style scoped>

</style>
