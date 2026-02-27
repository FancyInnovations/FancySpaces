import {useUserStore} from "@/stores/user.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import type {RecordQueryResult} from "@/api/analytics/metricrecords/types.ts";
import {ANALYTICS_CORE_API_BASE_URL} from "@/api/analytics/url.ts";

export async function getLatestRecordsByCount(projectID: string, metricId: string, count: number = 100): Promise<RecordQueryResult> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${projectID}/metrics/${metricId}/records?amount=` + count,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
        }
    );

    if (!response.ok || response.status !== 200) {
        useNotificationStore().error("Failed to fetch latest records by count: " + await response.text());
        throw new Error("Failed to fetch latest records by count: " + await response.text());
    }

    const data: RecordQueryResult = await response.json();

    for (let i = 0; i < data.records.length; i++) {
        data.records[i]!.timestamp = new Date(data.records[i]!.timestamp);
    }

    return data;
}

export async function getLatestRecordsPerMinute(projectID: string, metricId: string): Promise<RecordQueryResult> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${projectID}/metrics/${metricId}/records/per-minute`,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
        }
    );

    if (!response.ok || response.status !== 200) {
        useNotificationStore().error("Failed to fetch latest records per minute: " + await response.text());
        throw new Error("Failed to fetch latest records per minute: " + await response.text());
    }

    const data: RecordQueryResult = await response.json();

    for (let i = 0; i < data.records.length; i++) {
        data.records[i]!.timestamp = new Date(data.records[i]!.timestamp);
    }

    return data;
}

export async function getLatestRecordsPerHour(projectID: string, metricId: string): Promise<RecordQueryResult> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${projectID}/metrics/${metricId}/records/per-hour`,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
        }
    );

    if (!response.ok || response.status !== 200) {
        useNotificationStore().error("Failed to fetch latest records per hour: " + await response.text());
        throw new Error("Failed to fetch latest records per hour: " + await response.text());
    }

    const data: RecordQueryResult = await response.json();

    for (let i = 0; i < data.records.length; i++) {
        data.records[i]!.timestamp = new Date(data.records[i]!.timestamp);
    }

    return data;
}

export async function getLatestRecordsPerDay(projectID: string, metricId: string): Promise<RecordQueryResult> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/${projectID}/metrics/${metricId}/records/per-day`,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
        }
    );

    if (!response.ok || response.status !== 200) {
        useNotificationStore().error("Failed to fetch latest records per day: " + await response.text());
        throw new Error("Failed to fetch latest records per day: " + await response.text());
    }

    const data: RecordQueryResult = await response.json();

    for (let i = 0; i < data.records.length; i++) {
        data.records[i]!.timestamp = new Date(data.records[i]!.timestamp);
    }

    return data;
}

