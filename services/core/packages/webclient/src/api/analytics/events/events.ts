import type { EventQueryResult } from '@/api/analytics/events/types'
import { ANALYTICS_CORE_API_BASE_URL } from '@/api/analytics/url'
import { useNotificationStore } from '@/stores/notifications'
import { useUserStore } from '@/stores/user'

export async function getLatestEventsByTime (projectID: string, name = '', hours = 1): Promise<EventQueryResult> {
  const response = await fetch(
    `${ANALYTICS_CORE_API_BASE_URL}/projects/` + projectID + `/events?name=${name}&time=${hours}`,
    {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + useUserStore().token!,
      },
    },
  )

  if (!response.ok || response.status !== 200) {
    useNotificationStore().error('Failed to fetch latest events by time: ' + await response.text())
    throw new Error('Failed to fetch latest events by time: ' + await response.text())
  }

  const data: EventQueryResult = await response.json()

  for (let i = 0; i < data.events.length; i++) {
    data.events[i]!.id = i.toString()
    data.events[i]!.timestamp = new Date(data.events[i]!.timestamp)
  }

  return data
}

export async function getLatestEventsByCount (projectID: string, name = '', count = 100): Promise<EventQueryResult> {
  const response = await fetch(
    `${ANALYTICS_CORE_API_BASE_URL}/projects/` + projectID + `/events?name=${name}&amount=${count}`,
    {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + useUserStore().token!,
      },
    },
  )

  if (!response.ok || response.status !== 200) {
    useNotificationStore().error('Failed to fetch latest events by count: ' + await response.text())
    throw new Error('Failed to fetch latest events by count: ' + await response.text())
  }

  const data: EventQueryResult = await response.json()

  for (let i = 0; i < data.events.length; i++) {
    data.events[i]!.id = i.toString()
    data.events[i]!.timestamp = new Date(data.events[i]!.timestamp)
  }

  return data
}
