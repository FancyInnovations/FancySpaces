<script lang="ts" setup>

const props = defineProps<{
  version: string,
  spaceID: string,
  density?: null | 'default' | 'comfortable' | 'compact'
}>();

</script>

<template>
  <v-menu
    close-delay="100"
    location="top"
    offset="8"
    open-delay="100"
    open-on-hover
  >
    <template #activator="{ props: menuProps }">
      <v-chip
        :density="props.density || 'default'"
        :to="`/spaces/${props.spaceID}/versions/${encodeURIComponent(props.version)}`"
        class="ma-1"
        color="primary"
        prepend-icon="mdi-tag-outline"
        rounded
        v-bind="menuProps"
        variant="tonal"
      >
        v{{ props.version }}
      </v-chip>
    </template>

    <v-card min-width="220">
      <v-card-text>
        <strong>Space:</strong> {{ props.spaceID }}<br />
        <strong>Version:</strong> {{ props.version }}
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<style scoped>

</style>
