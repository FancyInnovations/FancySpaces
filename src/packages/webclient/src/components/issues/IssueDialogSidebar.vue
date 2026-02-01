<script lang="ts" setup>

import type {Issue, IssueComment} from "@/api/issues/types.ts";
import IssueIDChip from "@/components/issues/IssueIDChip.vue";
import IssueStatusChip from "@/components/issues/IssueStatusChip.vue";
import IssuePriorityChip from "@/components/issues/IssuePriorityChip.vue";
import IssueTypeChip from "@/components/issues/IssueTypeChip.vue";
import IssueExternalSourceChip from "@/components/issues/IssueExternalSourceChip.vue";
import VersionChip from "@/components/VersionChip.vue";

const props = defineProps<{
  issue?: Issue,
  comments: IssueComment[]
}>();

const formattedCreatedAt = ref('');
const formattedUpdatedAt = ref('');
const formattedResolvedAt = ref('');

function formatDate(date: Date): string {
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const diffInHours = diff / (1000 * 60 * 60);
  if (diffInHours < 12) {
    if (diffInHours < 1) {
      const diffInMinutes = diff / (1000 * 60);
      if (diffInMinutes < 1) {
        return `${Math.floor(diff / 1000)}s ago`;
      } else {
        return `${Math.floor(diffInMinutes)}min ago`;
      }
    } else {
      return `${Math.floor(diffInHours)}h ago`;
    }
  } else {
    const createdAtDate = date.getDate();
    const nowDate = now.getDate();
    const createdAtMonth = date.getMonth();
    const nowMonth = now.getMonth();
    const createdAtYear = date.getFullYear();
    const nowYear = now.getFullYear();

    if (createdAtYear === nowYear && createdAtMonth === nowMonth && createdAtDate === nowDate - 1) {
      return 'Yesterday';
    } else if (createdAtYear === nowYear && createdAtMonth === nowMonth && createdAtDate === nowDate) {
      return 'Today';
    } else {
      return date.toLocaleDateString();
    }
  }
}

onMounted(() => {
  formattedCreatedAt.value = formatDate(new Date(props.issue?.created_at || ''));
  formattedUpdatedAt.value = formatDate(new Date(props.issue?.updated_at || ''));
  formattedResolvedAt.value = props.issue?.resolved_at
    ? formatDate(new Date(props.issue?.resolved_at))
    : 'Unresolved';

  setInterval(() => {
    formattedCreatedAt.value = formatDate(new Date(props.issue?.created_at || ''));
    formattedUpdatedAt.value = formatDate(new Date(props.issue?.updated_at || ''));
    formattedResolvedAt.value = props.issue?.resolved_at
      ? formatDate(new Date(props.issue?.resolved_at))
      : 'Unresolved';
  }, 1000);
});

</script>

<template>
  <v-list
    v-if="props.issue"
    class="sidebar__background"
    elevation="12"
    location="right"
    min-width="250px"
    rounded="xl"
  >
    <v-list-item>
      <v-list-item-title>
        ID:
        <IssueIDChip
          :issue="props.issue"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item v-if="props.issue?.parent_issue">
      <v-list-item-title>
        Parent:
        <IssueIDChip
          :issueName="props.issue?.parent_issue"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item v-if="props.issue?.external_source">
      <v-list-item-title>
        Source:
        <IssueExternalSourceChip
          :issue="props.issue!"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item>
      <v-list-item-title>
        Type:
        <IssueTypeChip
          :issue="props.issue!"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item>
      <v-list-item-title>
        Status:
        <IssueStatusChip
          :issue="props.issue!"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item>
      <v-list-item-title>
        Priority:
        <IssuePriorityChip
          :issue="props.issue!"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item v-if="props.issue?.affected_versions && props.issue?.affected_versions.length > 0">
      <v-list-item-title>
        Affected versions:
      </v-list-item-title>
    </v-list-item>

    <div
      v-if="props.issue?.affected_versions && props.issue?.affected_versions.length > 0"
      class="ml-4"
    >
      <VersionChip
        v-for="version in props.issue.affected_versions"
        :key="version"
        :spaceID="issue!.space"
        :version="version"
        density="compact"
      />
    </div>

    <v-list-item v-if="props.issue?.fix_version">
      <v-list-item-title>
        Fix version:
        <VersionChip
          :spaceID="issue!.space"
          :version="props.issue.fix_version"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item
      v-if="props.issue?.resolved_at"
      :title="'Resolved at: ' + formattedResolvedAt"
    />

    <v-list-item>
      <v-list-item-title>
        Reporter:
        <UserChip
          :user="props.issue!.reporter"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item v-if="props.issue?.assignee">
      <v-list-item-title>
        Assignee:
        <UserChip
          :user="props.issue!.assignee"
          class="ml-2"
          density="compact"
        />
      </v-list-item-title>
    </v-list-item>

    <v-list-item
      :title="'Created at: ' + formattedCreatedAt"
    />

    <v-list-item
      :title="'Last updated: ' + formattedUpdatedAt"
    />

    <v-list-item
      :title="'Comments: ' + props.comments?.length"
    />
  </v-list>
</template>

<style scoped>
.sidebar__background {
  background-color: transparent !important;
  border: 1px solid rgba(255, 255, 255, 0.1);
  scrollbar-width: none;
}
</style>
