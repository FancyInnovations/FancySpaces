<script lang="ts" setup>
import type {SpaceVersion, SpaceVersionFile} from "@/api/versions/types";
import {createVersion, deleteVersionFile, updateVersion, uploadVersionFile} from "@/api/versions/versions";
import {useConfirmationStore} from "@/stores/confirmation";
import {useNotificationStore} from "@/stores/notifications";

const props = defineProps<{ spaceID: string; spaceSlug: string; version?: SpaceVersion }>();
const emit = defineEmits<{ saved: [version: SpaceVersion] }>();
const router = useRouter();
const notifications = useNotificationStore();
const confirmation = useConfirmationStore();

const name = ref("");
const platform = ref("paper");
const channel = ref("release");
const supportedPlatforms = ref<string[]>([]);
const changelog = ref("");
const selectedFiles = ref<File[]>([]);
const files = ref<SpaceVersionFile[]>([]);
const saving = ref(false);
const error = ref("");

const platforms = [
1  { title: "Bukkit", value: "bukkit" }, { title: "Spigot", value: "spigot" },
  { title: "Paper", value: "paper" }, { title: "Purpur", value: "purpur" },
  { title: "Folia", value: "folia" }, { title: "BungeeCord", value: "bungeecord" },
  { title: "Waterfall", value: "waterfall" }, { title: "Velocity", value: "velocity" },
  { title: "Fabric", value: "fabric" }, { title: "Forge", value: "forge" },
  { title: "Quilt", value: "quilt" }, { title: "LiteLoader", value: "liteloader" },
  { title: "Hytale Plugin", value: "hytale_plugin" }, { title: "Executable", value: "executable" },
];
const channels = [
  { title: "Release", value: "release" }, { title: "Beta", value: "beta" }, { title: "Alpha", value: "alpha" },
];

watch(() => props.version, (version) => {
  if (!version) return;
  name.value = version.name;
  platform.value = version.platform;
  channel.value = version.channel;
  supportedPlatforms.value = [...version.supported_platform_versions];
  changelog.value = version.changelog;
  files.value = [...version.files];
}, { immediate: true });

const valid = computed(() => name.value.trim().length > 0 && !!platform.value && !!channel.value);

function formatSize(size: number) {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`;
  return `${(size / (1024 * 1024)).toFixed(2)} MB`;
}

async function save() {
  error.value = "";
  if (!valid.value) {
    error.value = "Version name, platform, and channel are required.";
    return;
  }
  saving.value = true;
  try {
    const data = {
      name: name.value.trim(), platform: platform.value, channel: channel.value,
      changelog: changelog.value, supported_platform_versions: supportedPlatforms.value.filter(Boolean),
    };
    const version = props.version
      ? await updateVersion(props.spaceID, props.version.id, data)
      : await createVersion(props.spaceID, data);

    for (const file of selectedFiles.value) {
      await uploadVersionFile(props.spaceID, version.id, file);
    }
    selectedFiles.value = [];
    notifications.info(props.version ? "Version updated successfully" : "Version created successfully");
    emit("saved", version);
    await router.push(`/spaces/${props.spaceSlug}/versions/${encodeURIComponent(version.name)}`);
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Unable to save version.";
    notifications.error(error.value);
  } finally {
    saving.value = false;
  }
}

function removeFile(file: SpaceVersionFile) {
  confirmation.confirmation = {
    shown: true, persistent: true, title: "Delete release file",
    text: `Delete ${file.name}? This action cannot be undone.`, yesText: "Delete",
    onConfirm: async () => {
      if (!props.version) return;
      try {
        await deleteVersionFile(props.spaceID, props.version.id, file.name);
        files.value = files.value.filter(item => item.name !== file.name);
        notifications.info("File deleted");
      } catch (e) {
        notifications.error(e instanceof Error ? e.message : "Unable to delete file.");
      }
    },
  };
}
</script>

<template>
  <v-row justify="center">
    <v-col md="8">
      <Card>
        <v-card-title class="mt-2">{{ props.version ? 'Edit Version' : 'New Version' }}</v-card-title>
        <v-card-text>
          <v-alert v-if="error" class="mb-4" type="error" variant="tonal">{{ error }}</v-alert>
          <v-text-field v-model="name" class="mb-4" :disabled="saving" color="primary" label="Version name" required />
          <div class="d-flex mb-4">
            <v-select v-model="platform" :disabled="saving" :items="platforms" class="mr-2" color="primary" label="Platform" required />
            <v-select v-model="channel" :disabled="saving" :items="channels" class="ml-2" color="primary" label="Channel" required />
          </div>
          <v-combobox v-model="supportedPlatforms" :disabled="saving" chips closable-chips color="primary" hint="Press Enter after each platform version" label="Supported platform versions" multiple persistent-hint class="mb-4" />
          <v-textarea v-model="changelog" :disabled="saving" color="primary" label="Changelog" rows="10" />
          <v-file-input v-model="selectedFiles" :disabled="saving" chips clearable color="primary" label="Add or replace release files" multiple show-size />

          <v-list v-if="files.length" class="bg-transparent mt-4">
            <v-list-subheader>Existing files</v-list-subheader>
            <v-list-item v-for="file in files" :key="file.name" :title="file.name" :subtitle="formatSize(file.size)">
              <template #append>
                <v-btn :href="file.url" icon="mdi-download" target="_blank" variant="text" />
                <v-btn v-if="props.version" color="red" icon="mdi-delete" variant="text" @click="removeFile(file)" />
              </template>
            </v-list-item>
          </v-list>

          <div class="d-flex ga-3 mt-6">
            <v-btn :disabled="saving" :to="`/spaces/${props.spaceSlug}/versions`" variant="text">Cancel</v-btn>
            <v-btn :loading="saving" color="primary" @click="save">{{ props.version ? 'Save Changes' : 'Create Version' }}</v-btn>
          </div>
        </v-card-text>
      </Card>
    </v-col>
  </v-row>
</template>
