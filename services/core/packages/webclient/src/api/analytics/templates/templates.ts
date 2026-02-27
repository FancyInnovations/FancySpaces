import {useNotificationStore} from "@/stores/notifications.ts";
import type {Template} from "@/api/analytics/templates/types.ts";
import {useUserStore} from "@/stores/user.ts";
import {ANALYTICS_CORE_API_BASE_URL} from "@/api/analytics/url.ts";

export async function getTemplates(): Promise<Template[]> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/templates`,
        {
            method: "GET",
            headers: {
                "Accept": "application/json",
            }
        },
    );

    if (!response.ok) {
        useNotificationStore().error("Failed to fetch templates: " + await response.text());
        throw new Error("Failed to fetch templates: " + await response.text());
    }

    return await response.json() as Template[];
}

export async function applyTemplate(projectId: string, templateId: string): Promise<void> {
    const response = await fetch(
        `${ANALYTICS_CORE_API_BASE_URL}/templates/${templateId}/apply`,
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + useUserStore().token!,
            },
            body: JSON.stringify(
                {
                    project_id: projectId
                }
            ),
        }
    );

    if (!response.ok) {
        useNotificationStore().error("Failed to apply template: " + await response.text());
        throw new Error("Failed to apply template: " + await response.text());
    }
}

