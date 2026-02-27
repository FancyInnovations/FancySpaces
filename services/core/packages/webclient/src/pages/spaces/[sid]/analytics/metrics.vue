<script lang="ts" setup>

import type {Space} from "@/api/spaces/types.ts";
import {getSpace} from "@/api/spaces/spaces.ts";
import {useHead} from "@vueuse/head";
import SpaceHeader from "@/components/SpaceHeader.vue";
import {useUserStore} from "@/stores/user.ts";
import AnalyticsSidebar from "@/components/analytics/AnalyticsSidebar.vue";
import Card from "@/components/common/Card.vue";
import type {Dashboard} from "@/api/analytics/dashboards/types.ts";
import {getDashboards} from "@/api/analytics/dashboards/dashboards.ts";
import type {Metric} from "@/api/analytics/metrics/types.ts";
import {createMetric, getMetrics} from "@/api/analytics/metrics/metrics.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import type {RecordQueryResult} from "@/api/analytics/metricrecords/types.ts";
import {getLatestRecordsByCount} from "@/api/analytics/metricrecords/metricrecords.ts";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
const notifications = useNotificationStore();

const isLoggedIn = ref(false);
const isMember = ref(false);

const space = ref<Space>();
const dashboards = ref<Dashboard[]>();
const metrics = ref<Metric[]>([]);

const selectedMetricIDToQuery = ref<string|null>(null);
const queryResult = ref<RecordQueryResult | null>(null);

const showNewMetricDialog = ref(false);
const aggregationIntervalTicks = {
  0: '10s',
  1: '1min',
  2: '5mins',
  3: '15mins',
  4: '30mins',
  5: '1h',
};
const newMetricName = ref<string>('');
const newMetricMultiSender = ref<boolean>(false);
const newMetricAggregationInterval = ref<number>(2);
const newMetricExtraAggregation = ref<boolean>(false);

onMounted(async () => {
  isLoggedIn.value = await userStore.isAuthenticated;

  const spaceID = (route.params as any).sid as string;
  space.value = await getSpace(spaceID);

  if (!space.value.analytics_settings.enabled) {
    router.push(`/spaces/${space.value.slug}`);
    return;
  }

  isMember.value = (await userStore.isAuthenticated) && (space.value.creator == userStore.user?.id || space.value.members.some(member => member.user_id === userStore.user?.id));

  if (!isMember.value) {
    router.push(`/spaces/${space.value.slug}/analytics`);
    return;
  }

  dashboards.value = await getDashboards(space.value?.id);
  await refreshMetrics();

  useHead({
    title: `${space.value.title} analytics portal - FancySpaces`,
    meta: [
      {
        name: 'description',
        content: space.value.summary || `Explore the ${space.value.title} project space on FancySpaces.`
      }
    ]
  });
});

async function refreshMetrics(sendNotification = false) {
  metrics.value = await getMetrics(space.value!.id);

  if (sendNotification) {
    notifications.info("Metrics refreshed");
  }
}

async function createMetricReq() {
  let aggregationInterval = 300;
  switch (newMetricAggregationInterval.value) {
    case 0:
      aggregationInterval = 10;
      break;
    case 1:
      aggregationInterval = 60;
      break;
    case 2:
      aggregationInterval = 60 * 5;
      break;
    case 3:
      aggregationInterval = 60 * 15;
      break;
    case 4:
      aggregationInterval = 60 * 30;
      break;
    case 5:
      aggregationInterval = 60 * 60;
      break;
  }

  const newMetric: Metric = {
    project_id: space.value!.id,
    metric_id: "will-be-generated",
    name: newMetricName.value,
    multi_sender: newMetricMultiSender.value,
    aggregation_interval: aggregationInterval,
    apply_extra_aggregation: newMetricExtraAggregation.value,
    pull_metric: false,
  };

  try {
    await createMetric(newMetric);
  } catch (e) {
    console.error("Failed to create metric:", e);
    return;
  }

  await refreshMetrics();

  newMetricName.value = '';
  newMetricMultiSender.value = false;
  newMetricAggregationInterval.value = 2;
  newMetricExtraAggregation.value = false;
  showNewMetricDialog.value = false;

  notifications.info("Metric created successfully");
}

async function queryMetrics() {
  if (!selectedMetricIDToQuery.value) {
    notifications.error("Please select a metric to query");
    return;
  }

  queryResult.value = await getLatestRecordsByCount(space.value!.id, selectedMetricIDToQuery.value, 10);
}
</script>

<template>
  <v-container width="90%">
    <v-row>
      <v-col class="flex-grow-0 pa-0">
        <AnalyticsSidebar
          :dashboards="dashboards"
          :space="space"
        />
      </v-col>

      <v-col>
        <SpaceHeader :space="space"></SpaceHeader>

        <hr
          class="grey-border-color mt-4"
        />
      </v-col>
    </v-row>


    <v-row justify="center">
      <v-col md="8">
        <Card>
          <v-card-title class="mt-2">
            Metric definitions
          </v-card-title>

          <v-card-text>
            <v-table density="compact" hover>
              <thead>
              <tr>
                <th>Metric Name</th>
                <th>Multi sender</th>
                <th>Pull metric</th>
                <th>Aggregation interval</th>
                <th>Actions</th>
              </tr>
              </thead>
              <tbody>
              <tr v-for="m in metrics" :key="m.metric_id">
                <td>{{ m.name }}</td>
                <td>{{ m.multi_sender ? 'Yes' : 'No' }}</td>
                <td>{{ m.pull_metric ? 'Yes' : 'No' }}</td>
                <td>{{ m.aggregation_interval }} seconds</td>
                <td style="width: 140px">
                  <v-btn
                    icon="mdi-pencil"
                    variant="text"
                  />
                  <v-btn
                    class="ml-2"
                    color="red"
                    icon="mdi-delete"
                    variant="text"
                  />
                </td>
              </tr>
              </tbody>
            </v-table>

            <v-btn
              class="mt-4"
              color="primary"
              prepend-icon="mdi-plus"
              variant="tonal"
              @click="showNewMetricDialog = true"
            >
              New Metric
            </v-btn>
             <v-btn
              class="ml-4 mt-4"
              color="primary"
              prepend-icon="mdi-refresh"
              variant="tonal"
              @click="refreshMetrics(true)"
            >
              Refresh
            </v-btn>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>

    <v-row class="mt-8" justify="center">
      <v-col md="8">
        <Card>
          <v-card-title class="mt-2">
            Metric query
          </v-card-title>

          <v-card-text>
            <div class="d-flex align-center">
              <v-select
                v-model="selectedMetricIDToQuery"
                :items="metrics"
                color="primary"
                density="compact"
                hide-details
                item-title="name"
                item-value="metric_id"
                label="Metric"
                max-width="300"
                placeholder="Select Metric"
              />

              <v-btn
                class="ml-4"
                  variant="tonal"
                @click="queryMetrics()"
              >
                Query
              </v-btn>
            </div>
          </v-card-text>
        </Card>
      </v-col>
    </v-row>

    <v-row v-if="queryResult" justify="center">
      <v-col md="8">
        <Card>
          <v-card-title class="mt-2">
            Metric query result
          </v-card-title>

          <v-card-subtitle>
            <span>Rows returned: {{ queryResult.records ? queryResult.records.length : 0 }}</span>
            <span class="ml-4">Query time: {{ queryResult.time }} ms</span>
          </v-card-subtitle>

          <v-card-text v-if="queryResult!.records">
            <v-table density="compact" hover>
              <thead>
              <tr>
                <th>Timestamp</th>
                <th>Label</th>
                <th>Value</th>
              </tr>
              </thead>
              <tbody>
              <tr
                v-for="record in queryResult!.records"
                :key="record.timestamp.toISOString()"
              >
                <td>{{ record.timestamp.toLocaleString() }}</td>
                <td>{{ record.label || "&ltempty>" }}</td>
                <td>{{ record.value }}</td>
              </tr>
              </tbody>
            </v-table>

<!--            <LineChart-->
<!--              :duration-ms="queryResult.time/1000000"-->
<!--              :values="queryResult!.records"-->
<!--              title="Preview latest records"-->
<!--            />-->
          </v-card-text>
          <v-card-text v-else>
            No records found for this metric.
          </v-card-text>
        </Card>
      </v-col>
    </v-row>
  </v-container>

  <Dialog :shown="showNewMetricDialog" persistent>
    <v-card
      elevation="8"
      min-width="500"
      rounded="xl"
    >
      <v-card-title class="mx-2 mt-2">
        New Metric
      </v-card-title>

      <v-card-text>
        <v-form>
          <v-text-field
            v-model="newMetricName"
            color="primary"
            label="Metric Name"
            required
          />

          <v-checkbox
            v-model="newMetricMultiSender"
            color="primary"
            hide-details
            label="Multi sender"
          />

          <v-label class="mt-4">
            Aggregation interval
          </v-label>
          <v-slider
            v-model="newMetricAggregationInterval"
            :max="5"
            :ticks="aggregationIntervalTicks"
            color="primary"
            show-ticks="always"
            step="1"
            tick-size="4"
          />

          <v-checkbox
            v-model="newMetricExtraAggregation"
            color="primary"
            hide-details
            label="Apply extra aggregation"
          />
        </v-form>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>

        <v-btn
          @click="createMetricReq"
        >
          Create
        </v-btn>
        <v-btn
          @click="showNewMetricDialog = false"
        >
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </Dialog>
</template>

<style scoped>

</style>
