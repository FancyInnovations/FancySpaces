import {useNotificationStore} from "@/stores/notifications.ts";
import {useUserStore} from "@/stores/user.ts";
import type {Metric} from "@/api/analytics/metrics/types.ts";
import {ANALYTICS_CORE_API_BASE_URL} from "@/api/analytics/url.ts";

export async function getMetrics(projectId: string): Promise<Metric[]> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${projectId}/metrics`,
        {
            method: "GET",
            headers: {
                "Accept": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            }
        },
    );

    if (!response.ok) {
        useNotificationStore().error("Failed to fetch metrics: " + await response.text());
        throw new Error("Failed to fetch metrics: " + await response.text());
    }

    return await response.json() as Metric[];
}

export async function createMetric(newMetric: Metric): Promise<void> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${newMetric.project_id}/metrics`,
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
            body: JSON.stringify(newMetric),
        }
    );

    if (!response.ok || response.status !== 201) {
        useNotificationStore().error("Failed to create metric: " + await response.text());
        throw new Error("Failed to create metric: " + await response.text());
    }
}

export async function updateMetric(updatedMetric: Metric): Promise<void> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${updatedMetric.project_id}/metrics/${updatedMetric.metric_id}`,
        {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
            body: JSON.stringify(updatedMetric),
        }
    );

    if (!response.ok) {
        useNotificationStore().error("Failed to update metric: " + await response.text());
        throw new Error("Failed to update metric: " + await response.text());
    }
}

export async function deleteMetric(projectId: string, metricId: string): Promise<void> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${projectId}/metrics/${metricId}`,
        {
            method: "DELETE",
            headers: {
                "Authorization": "Bearer " + useUserStore().token!,
            }
        }
    );

    if (!response.ok) {
        useNotificationStore().error("Failed to delete metric: " + await response.text());
        throw new Error("Failed to delete metric: " + await response.text());
    }
}
