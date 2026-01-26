export interface Issue {
  id: string;
  space: string;
  title: string;
  description: string;
  type: 'epic' | 'bug' | 'task' | 'story' | 'idea';
  status: 'todo' | 'in_progress' | 'done' | 'closed';
  priority: 'low' | 'medium' | 'high' | 'critical';
  assignee: string;
  reporter: string;
  created_at: Date;
  updated_at: Date;
  external_source: 'github' | 'discord' | null;
}

export interface IssueComment {
  id: string;
  issue: string;
  author: string;
  content: string;
  created_at: Date;
  updated_at: Date;
}
