<script lang="ts" setup>

</script>

<template>
  <v-footer
    class="mt-16"
    color="tertiary-container"
    height="fit-content"
  >
    <v-container width="70%">
      <v-row justify="space-evenly">
        <v-col class="mr-16" md="auto">
          <h1 class="text-h3 mb-2">FancySpaces</h1>
          <p class="mb-4">Download platform for products by FancyInnovations.</p>
          <p>&copy; 2025-{{ (new Date()).getFullYear() }} <span class="d-none d-sm-inline-block">FancyInnovations</span>
          </p>
        </v-col>

        <v-col md="auto">
          <h3 class="mb-2">Links</h3>
          <a class="hide-link" href="https://github.com/fancyinnovations/fancyspaces" target="_blank"><p>
            GitHub</p></a>
          <a class="hide-link" href="https://discord.gg/ZUgYCEJUEx" target="_blank"><p>Discord</p></a>
        </v-col>

        <v-col md="auto">
          <h3 class="mb-2">Resources</h3>
          <a class="hide-link" href="https://discord.gg/ZUgYCEJUEx" target="_blank"><p>Help Center</p></a>
          <a class="hide-link" href="https://discord.gg/ZUgYCEJUEx" target="_blank"><p>Report issues</p></a>
          <a class="hide-link" href="https://github.com/fancyinnovations/fancyspaces" target="_blank"><p>Documentation</p></a>
        </v-col>

        <v-col md="auto">
          <h3 class="mb-2">Other products</h3>
          <a class="hide-link" href="https://fancyanalytics.net" target="_blank"><p>FancyAnalytics</p></a>
        </v-col>

        <v-col md="auto">
          <h3 class="mb-2">Legal</h3>
          <a class="hide-link" href="/coming-soon" target="_blank"><p>Terms of Use</p></a>
          <a class="hide-link" href="/coming-soon" target="_blank"><p>Privacy Policy</p></a>
          <a class="hide-link" href="/coming-soon" target="_blank"><p>Copyright Policy</p></a>
        </v-col>
      </v-row>

      <v-row class="mt-8" justify="center">
        <v-col md="auto">
          <p class="text-caption text-disabled">
            NOT AN OFFICIAL MINECRAFT OR HYTALE SERVICE.
            NOT APPROVED BY OR ASSOCIATED WITH MOJANG, MICROSOFT OR HYPIXEL STUDIOS.
          </p>
        </v-col>
      </v-row>
    </v-container>
  </v-footer>
</template>

<style scoped>

</style>
