<script lang="ts" setup>

</script>

<template>
  <v-card
    class="card__border"
    color="#19120D33"
    elevation="12"
    rounded="xl"
  >
    <slot/>
  </v-card>
</template>

<style scoped>
.card__border {
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}
</style>
