<script lang="ts" setup>

const props = defineProps<{
  shown: boolean;
  persistent?: boolean;
  width?: string;
  height?: string;
}>();

</script>

<template>
  <v-dialog
    v-model="props.shown"
    :height="props.height ?? 'max-content'"
    :persistent="props.persistent ?? false"
    :width="props.width ?? 'max-content'"
    opacity="0.7"
    transition="dialog-transition"
  >
    <div class="dialog__background">
      <slot></slot>
    </div>
  </v-dialog>
</template>


<style scoped>
.dialog__background {
  height: 100%;
}

.dialog__background > * {
  background-color: rgba(25, 21, 13, 0.4) !important;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}
</style>
