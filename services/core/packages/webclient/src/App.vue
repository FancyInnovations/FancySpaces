<template>
  <v-app>
    <AppHeader/>

    <router-view />

    <AppFooter/>

    <IssueDialog/>
  </v-app>
</template>

<script lang="ts" setup>
//
import AppHeader from "@/components/AppHeader.vue";
</script>
