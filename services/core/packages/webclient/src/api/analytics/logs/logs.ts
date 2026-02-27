import type {LogsQueryResult} from "@/api/analytics/logs/types.ts";
import {useUserStore} from "@/stores/user.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import {ANALYTICS_CORE_API_BASE_URL} from "@/api/analytics/url.ts";

export async function getLatestLogsByTime(projectID: string, service: string = "", hours: number = 1): Promise<LogsQueryResult> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/` + projectID + "/logs?time=" + hours,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
        }
    );

    if (!response.ok || response.status !== 200) {
        useNotificationStore().error("Failed to fetch latest logs by time: " + await response.text());
        throw new Error("Failed to fetch latest logs by time: " + await response.text());
    }

    const data: LogsQueryResult = await response.json();

    for (let i = 0; i < data.records.length; i++) {
        data.records[i]!.id = i.toString();
        data.records[i]!.timestamp = new Date(data.records[i]!.timestamp);
    }

    return data;
}

export async function getLatestLogsByCount(projectID: string, service: string = "", count: number = 100): Promise<LogsQueryResult> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/` + projectID + "/logs?amount=" + count,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
        }
    );

    if (!response.ok || response.status !== 200) {
        useNotificationStore().error("Failed to fetch latest logs by count: " + await response.text());
        throw new Error("Failed to fetch latest logs by count: " + await response.text());
    }

    const data: LogsQueryResult = await response.json();

    for (let i = 0; i < data.records.length; i++) {
        data.records[i]!.id = i.toString();
        data.records[i]!.timestamp = new Date(data.records[i]!.timestamp);
    }

    return data;
}
