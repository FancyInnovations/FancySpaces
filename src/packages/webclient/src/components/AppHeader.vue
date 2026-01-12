<script lang="ts" setup>

</script>

<template>
  <v-app-bar
    class="app-bar__background"
  >
    <div class="w-100 d-flex justify-space-between align-center">
      <RouterLink to="/">
        <v-img
          class="mx-4"
          max-width="48"
          min-width="48"
          src="/logo.png"
        />
      </RouterLink>

      <div>
        <v-menu
          :close-on-content-click="false"
          open-delay="0"
          open-on-hover
        >
          <template v-slot:activator="{ props }">
            <v-btn
              append-icon="mdi-menu-down"
              color="secondary"
              prepend-icon="mdi-compass-outline"
              v-bind="props"
            >
              Explore projects
            </v-btn>
          </template>

          <v-list
            class="bg-secondary-container pa-2"
            rounded="xl"
          >
            <v-list-item
              class="mb-2"
              exact
              prepend-icon="mdi-earth"
              rounded="xl"
              title="All projects"
              to="/explore"
            />

            <v-list-item
              class="mb-2"
              prepend-icon="mdi-minecraft"
              rounded="xl"
              title="Minecraft plugins"
              to="/explore/minecraft-plugins"
            />

            <v-list-item
              class="mb-2"
              prepend-icon="mdi-sword-cross"
              rounded="xl"
              title="Hytale plugins"
              to="/explore/hytale-plugins"
            />

            <v-list-item
              class="mb-2"
              prepend-icon="mdi-cube-outline"
              rounded="xl"
              title="Other projects"
              to="/explore/other-projects"
            />

            <v-list-item
              prepend-icon="mdi-account-group-outline"
              rounded="xl"
              title="By other creators"
              to="/explore/by-other-creators"
            />

          </v-list>
        </v-menu>

        <v-menu
          :close-on-content-click="false"
          location="bottom center"
          open-delay="0"
          open-on-hover
        >
          <template v-slot:activator="{ props }">
            <v-btn
              append-icon="mdi-menu-down"
              class="ml-4"
              color="secondary"
              prepend-icon="mdi-toolbox-outline"
              v-bind="props"
            >
              Tools
            </v-btn>
          </template>

          <v-list
            class="bg-secondary-container pa-2"
            rounded="xl"
          >
            <v-list-item
              class="mb-2"
              exact
              prepend-icon="mdi-file-document-edit-outline"
              rounded="xl"
              subtitle="Live markdown preview"
              title="Markdown Editor"
              to="/tools/markdown-editor"
            />

            <v-list-item
              class="mb-2"
              exact
              prepend-icon="mdi-coffee-outline"
              rounded="xl"
              subtitle="SDK for the FancySpaces API"
              title="Java SDK"
              to="/tools/java-sdk"
            />

            <v-list-item
              class="mb-2"
              exact
              href="/spaces/fancyverteiler"
              prepend-icon="mdi-upload-network-outline"
              rounded="xl"
              subtitle="Deploy your plugins with ease"
              title="FancyVerteiler"
            />

            <v-list-item
              exact
              href="/spaces/fancyanalytics"
              prepend-icon="mdi-chart-box-outline"
              rounded="xl"
              subtitle="Powerful analytics for your projects"
              title="FancyAnalytics"
            />
          </v-list>
        </v-menu>
      </div>

      <v-menu>
        <template v-slot:activator="{ props }">
          <v-btn
            class="mr-4"
            icon="mdi-dots-vertical"
            v-bind="props"
          />
        </template>
        <v-list>
          <v-list-item
            href="https://github.com/fancyinnovations"
            prepend-icon="mdi-github"
            title="GitHub"/>

          <v-list-item
            href="https://fancyinnovations.com/docs/general"
            prepend-icon="mdi-script-text-outline"
            title="Documentation"/>

          <v-list-item
            href="https://discord.gg/ZUgYCEJUEx"
            prepend-icon="mdi-message"
            title="Discord"/>

        </v-list>
      </v-menu>
    </div>
  </v-app-bar>
</template>

<style scoped>
.app-bar__background {
  background-color: rgb(25, 18, 13, 0.2) !important;
  backdrop-filter: blur(15px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}
</style>
