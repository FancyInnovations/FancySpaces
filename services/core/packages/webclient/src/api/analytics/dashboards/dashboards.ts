import {useUserStore} from "@/stores/user.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import type {Chart, Dashboard} from "@/api/analytics/dashboards/types.ts";
import {ANALYTICS_CORE_API_BASE_URL} from "@/api/analytics/url.ts";

export async function getDashboards(projectID: string): Promise<Dashboard[]> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${projectID}/dashboards`,
        {
            method: "GET",
            headers: {
                "Accept": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            }
        },
    );

    if (!response.ok) {
        useNotificationStore().error("Failed to fetch dashboards: " + await response.text());
        throw new Error("Failed to fetch dashboards: " + await response.text());
    }

    return await response.json() as Dashboard[];
}

export async function getDashboard(projectID: string, dashboardID: string): Promise<Dashboard> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${projectID}/dashboards/${dashboardID}`,
        {
            method: "GET",
            headers: {
                "Accept": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            }
        },
    );

    if (!response.ok) {
        useNotificationStore().error("Failed to fetch dashboard: " + await response.text());
        throw new Error("Failed to fetch dashboard: " + await response.text());
    }

    return await response.json() as Dashboard;
}

export async function createDashboard(projectID: string, name: string, summary: string, isPublic: boolean, charts: Chart[]): Promise<void> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${projectID}/dashboards`,
        {
            method: "POST",
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
            body: JSON.stringify({
                name: name,
                summary: summary,
                public: isPublic,
                charts: charts,
            }),
        },
    );

    if (!response.ok) {
        useNotificationStore().error("Failed to create dashboard: " + await response.text());
        throw new Error("Failed to create dashboard: " + await response.text());
    }
}

export async function updateDashboard(projectID: string, dashboardID: string, name: string, summary: string, isPublic: boolean, charts: Chart[]): Promise<void> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${projectID}/dashboards/${dashboardID}`,
        {
            method: "PUT",
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
            body: JSON.stringify({
                name: name,
                summary: summary,
                public: isPublic,
                charts: charts,
            }),
        },
    );

    if (!response.ok) {
        useNotificationStore().error("Failed to update dashboard: " + await response.text());
        throw new Error("Failed to update dashboard: " + await response.text());
    }
}

export async function deleteDashboard(projectID: string, dashboardID: string): Promise<void> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${projectID}/dashboards/${dashboardID}`,
        {
            method: "DELETE",
            headers: {
                "Accept": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
        },
    );

    if (!response.ok) {
        useNotificationStore().error("Failed to delete dashboard: " + await response.text());
        throw new Error("Failed to delete dashboard: " + await response.text());
    }
}

