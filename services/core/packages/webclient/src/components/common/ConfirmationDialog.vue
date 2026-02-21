<script lang="ts" setup>

const props = defineProps<{
    shown: boolean;
    persistent?: boolean;
    title: string;
    text: string;
    yesText?: string;
}>();

const emit = defineEmits(['clickedYes', 'clickedClose'])

</script>

<template>
    <Dialog
        :persistent="persistent"
        :shown="shown"
    >
        <v-card
            elevation="8"
            rounded="xl"
        >
            <v-card-title class="mt-2">
                {{ props.title }}
            </v-card-title>

            <v-card-text>
                <p>{{ props.text }}</p>
            </v-card-text>

            <v-card-actions>
                <v-spacer></v-spacer>

                <v-btn @click="emit('clickedYes')">{{ props.yesText ?? 'Yes' }}</v-btn>
                <v-btn @click="emit('clickedClose')">Close</v-btn>
            </v-card-actions>
        </v-card>
    </Dialog>
</template>

