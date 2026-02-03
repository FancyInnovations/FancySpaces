<script lang="ts" setup>

const props = defineProps<{
  user: string,
  density?: null | 'default' | 'comfortable' | 'compact'
}>();

</script>

<template>
  <v-chip
    :density="props.density || 'default'"
    color="secondary"
    prepend-icon="mdi-account"
    rounded
    variant="tonal"
  >
    {{ props.user }}
  </v-chip>
</template>

<style scoped>

</style>
