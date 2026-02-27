import {useUserStore} from "@/stores/user.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import type {EventQueryResult} from "@/api/analytics/events/types.ts";
import {ANALYTICS_CORE_API_BASE_URL} from "@/api/analytics/url.ts";

export async function getLatestEventsByTime(projectID: string, name: string = "", hours: number = 1): Promise<EventQueryResult> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/` + projectID + `/events?name=${name}&time=${hours}`,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
        }
    );

    if (!response.ok || response.status !== 200) {
        useNotificationStore().error("Failed to fetch latest events by time: " + await response.text());
        throw new Error("Failed to fetch latest events by time: " + await response.text());
    }

    const data: EventQueryResult = await response.json();

    for (let i = 0; i < data.events.length; i++) {
        data.events[i]!.id = i.toString();
        data.events[i]!.timestamp = new Date(data.events[i]!.timestamp);
    }

    return data;
}

export async function getLatestEventsByCount(projectID: string, name: string = "", count: number = 100): Promise<EventQueryResult> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/projects/` + projectID + `/events?name=${name}&amount=${count}`,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
        }
    );

    if (!response.ok || response.status !== 200) {
        useNotificationStore().error("Failed to fetch latest events by count: " + await response.text());
        throw new Error("Failed to fetch latest events by count: " + await response.text());
    }

    const data: EventQueryResult = await response.json();

    for (let i = 0; i < data.events.length; i++) {
        data.events[i]!.id = i.toString();
        data.events[i]!.timestamp = new Date(data.events[i]!.timestamp);
    }

    return data;
}
