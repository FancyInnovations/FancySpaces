export interface Issue {
  id: string;
  space: string;
  title: string;
  description: string;
  type: 'epic' | 'bug' | 'task' | 'story' | 'idea';
  status: 'backlog' | 'planned' | 'in_progress' | 'done' | 'closed';
  priority: 'low' | 'medium' | 'high' | 'critical';
  assignee?: string;
  reporter: string;
  created_at: Date;
  updated_at: Date;
  external_source: 'github' | 'discord_forum_post' | 'discord_ticket_bot' | null;
  fix_version?: string;
  affected_versions?: string[];
  resolved_at?: Date;
  parent_issue?: string;
  extra_fields?: Record<string, any>;
}

export interface IssueComment {
  id: string;
  issue: string;
  author: string;
  content: string;
  created_at: Date;
  updated_at: Date;
}
