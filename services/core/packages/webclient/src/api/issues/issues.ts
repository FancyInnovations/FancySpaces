import type { Issue, IssueActivity, IssueComment, IssueListResponse, IssueQuery } from '@/api/issues/types'
import { useUserStore } from '@/stores/user'

function authHeaders (includeContentType = false): HeadersInit {
  const userStore = useUserStore()
  return {
    ...(includeContentType ? { 'Content-Type': 'application/json' } : {}),
    Accept: 'application/json',
    ...(userStore.token ? { Authorization: `Bearer ${userStore.token}` } : {}),
  }
}

async function readError (response: Response): Promise<Error> {
  const body = await response.text()
  try {
    const parsed = JSON.parse(body)
    const detail = parsed.detail || parsed.message || parsed.title
    return new Error(detail || `Request failed (${response.status})`)
  } catch {
    return new Error(body || `Request failed (${response.status})`)
  }
}

function normalizeIssue (issue: Issue): Issue {
  issue.created_at = new Date(issue.created_at)
  issue.updated_at = new Date(issue.updated_at)
  if (issue.resolved_at) {
    issue.resolved_at = new Date(issue.resolved_at)
  }
  if (issue.archived_at) {
    issue.archived_at = new Date(issue.archived_at)
  }
  return issue
}

function normalizeComment (comment: IssueComment): IssueComment {
  comment.created_at = new Date(comment.created_at)
  comment.updated_at = new Date(comment.updated_at)
  return comment
}

function normalizeActivity (activity: IssueActivity): IssueActivity {
  activity.created_at = new Date(activity.created_at)
  return activity
}

export async function getIssue (spaceId: string, issueId: string): Promise<Issue> {
  const response = await fetch(`/api/v1/spaces/${spaceId}/issues/${issueId}`, {
    method: 'GET',
    headers: authHeaders(),
  })
  if (!response.ok) {
    throw await readError(response)
  }
  return normalizeIssue(await response.json() as Issue)
}

export async function getIssues (spaceId: string, query: IssueQuery = {}): Promise<IssueListResponse> {
  const params = new URLSearchParams()
  for (const [key, value] of Object.entries(query)) {
    if (value !== undefined && value !== '') {
      params.set(key, String(value))
    }
  }
  const suffix = params.toString() ? `?${params.toString()}` : ''
  const response = await fetch(`/api/v1/spaces/${spaceId}/issues${suffix}`, {
    method: 'GET',
    headers: authHeaders(),
  })
  if (!response.ok) {
    throw await readError(response)
  }

  const payload = await response.json()
  // Accept the legacy array response during rolling deployments.
  if (Array.isArray(payload)) {
    return { items: payload.map(normalizeIssue), total: payload.length, offset: 0, limit: payload.length }
  }
  return {
    ...payload,
    items: (payload.items || []).map(normalizeIssue),
  } as IssueListResponse
}

export async function getAllIssues (spaceId: string): Promise<Issue[]> {
  return (await getIssues(spaceId, { limit: 200 })).items
}

export async function createIssue (spaceId: string, issueData: Partial<Issue>): Promise<Issue> {
  const userStore = useUserStore()
  if (!(await userStore.isAuthenticated)) {
    throw new Error('User is not logged in')
  }
  const response = await fetch(`/api/v1/spaces/${spaceId}/issues`, {
    method: 'POST',
    headers: authHeaders(true),
    body: JSON.stringify(issueData),
  })
  if (!response.ok) {
    throw await readError(response)
  }
  return normalizeIssue(await response.json() as Issue)
}

export async function updateIssue (spaceId: string, issueID: string, issueData: Partial<Issue>): Promise<Issue> {
  const userStore = useUserStore()
  if (!(await userStore.isAuthenticated)) {
    throw new Error('User is not logged in')
  }
  const response = await fetch(`/api/v1/spaces/${spaceId}/issues/${issueID}`, {
    method: 'PATCH',
    headers: authHeaders(true),
    body: JSON.stringify(issueData),
  })
  if (!response.ok) {
    throw await readError(response)
  }
  return normalizeIssue(await response.json() as Issue)
}

export async function deleteIssue (spaceId: string, issueID: string): Promise<void> {
  const userStore = useUserStore()
  if (!(await userStore.isAuthenticated)) {
    throw new Error('User is not logged in')
  }
  const response = await fetch(`/api/v1/spaces/${spaceId}/issues/${issueID}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  if (!response.ok) {
    throw await readError(response)
  }
}

export async function getIssueComments (spaceId: string, issueID: string): Promise<IssueComment[]> {
  const response = await fetch(`/api/v1/spaces/${spaceId}/issues/${issueID}/comments`, {
    method: 'GET',
    headers: authHeaders(),
  })
  if (!response.ok) {
    throw await readError(response)
  }
  return (await response.json() as IssueComment[]).map(normalizeComment)
}

export async function addIssueComment (spaceId: string, issueID: string, content: string): Promise<IssueComment> {
  const userStore = useUserStore()
  if (!(await userStore.isAuthenticated)) {
    throw new Error('User is not logged in')
  }
  const response = await fetch(`/api/v1/spaces/${spaceId}/issues/${issueID}/comments`, {
    method: 'POST',
    headers: authHeaders(true),
    body: JSON.stringify({ content }),
  })
  if (!response.ok) {
    throw await readError(response)
  }
  return normalizeComment(await response.json() as IssueComment)
}

export async function updateIssueComment (spaceId: string, issueID: string, commentID: string, content: string): Promise<IssueComment> {
  const response = await fetch(`/api/v1/spaces/${spaceId}/issues/${issueID}/comments/${commentID}`, {
    method: 'PATCH',
    headers: authHeaders(true),
    body: JSON.stringify({ content }),
  })
  if (!response.ok) {
    throw await readError(response)
  }
  return normalizeComment(await response.json() as IssueComment)
}

export async function deleteIssueComment (spaceId: string, issueID: string, commentID: string): Promise<void> {
  const response = await fetch(`/api/v1/spaces/${spaceId}/issues/${issueID}/comments/${commentID}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  if (!response.ok) {
    throw await readError(response)
  }
}

export async function getIssueActivity (spaceId: string, issueID: string): Promise<IssueActivity[]> {
  const response = await fetch(`/api/v1/spaces/${spaceId}/issues/${issueID}/activity`, { headers: authHeaders() })
  if (!response.ok) {
    throw await readError(response)
  }
  return (await response.json() as IssueActivity[]).map(normalizeActivity)
}

export async function bulkUpdateIssues (spaceId: string, ids: string[], changes: Partial<Issue>): Promise<Issue[]> {
  const response = await fetch(`/api/v1/spaces/${spaceId}/issues/bulk`, {
    method: 'PATCH', headers: authHeaders(true), body: JSON.stringify({ ids, changes }),
  })
  if (!response.ok) {
    throw await readError(response)
  }
  return (await response.json() as Issue[]).map(normalizeIssue)
}
