<script lang="ts" setup>

  const props = defineProps<{
    user: string
    density?: null | 'default' | 'comfortable' | 'compact'
  }>()

</script>

<template>
  <v-chip
    color="secondary"
    :density="props.density || 'default'"
    prepend-icon="mdi-account"
    rounded
    :to="`/users/${props.user}`"
    variant="tonal"
  >
    {{ props.user }}
  </v-chip>
</template>

<style scoped>

</style>
