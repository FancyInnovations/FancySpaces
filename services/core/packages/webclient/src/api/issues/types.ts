export interface Issue {
  id: string
  space: string
  title: string
  description: string
  type: 'epic' | 'bug' | 'task' | 'story' | 'idea'
  status: 'backlog' | 'planned' | 'in_progress' | 'done' | 'closed'
  priority: 'low' | 'medium' | 'high' | 'critical'
  assignee?: string
  reporter: string
  created_at: Date
  updated_at: Date
  external_source?: 'github' | 'discord_forum_post' | 'discord_ticket_bot' | null
  external_id?: string
  external_url?: string
  fix_version?: string
  affected_versions?: string[]
  resolved_at?: Date
  parent_issue?: string
  extra_fields?: Record<string, any>
  archived_at?: Date
}

export interface IssueListResponse {
  items: Issue[]
  total: number
  offset: number
  limit: number
}

export interface IssueQuery {
  q?: string
  type?: Issue['type']
  status?: Issue['status']
  priority?: Issue['priority']
  assignee?: string
  external_source?: NonNullable<Issue['external_source']>
  offset?: number
  limit?: number
}

export interface IssueComment {
  id: string
  space: string
  issue: string
  author: string
  content: string
  created_at: Date
  updated_at: Date
}
