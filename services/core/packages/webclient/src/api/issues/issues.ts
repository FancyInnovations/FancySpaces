import type {Issue} from "@/api/issues/types.ts";
import {useUserStore} from "@/stores/user.ts";

export async function getIssue(spaceId: string, issueId: string): Promise<Issue> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/issues/${issueId}`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch issue: " + await response.text());
  }

  const issue = await response.json();
  issue.created_at = new Date(issue.created_at);
  issue.updated_at = new Date(issue.updated_at);
  if (issue.resolved_at) {
    issue.resolved_at = new Date(issue.resolved_at);
  }

  return issue as Issue;
}

export async function getAllIssues(spaceId: string): Promise<Issue[]> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/issues`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to all issues: " + await response.text());
  }

  const issues = await response.json();
  issues.forEach((issue: Issue) => {
    issue.created_at = new Date(issue.created_at);
    issue.updated_at = new Date(issue.updated_at);
    if (issue.resolved_at) {
      issue.resolved_at = new Date(issue.resolved_at);
    }
  });

  return issues as Issue[];
}

export async function createIssue(spaceId: string, issueData: Partial<Issue>): Promise<Issue> {
  const userStore = useUserStore();
  if (!userStore.isAuthenticated) {
    throw new Error("User is not logged in");
  }

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/issues`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      },
      body: JSON.stringify(issueData),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to create issue: " + await response.text());
  }

  const issue = await response.json();
  issue.created_at = new Date(issue.created_at);
  issue.updated_at = new Date(issue.updated_at);
  if (issue.resolved_at) {
    issue.resolved_at = new Date(issue.resolved_at);
  }

  return issue as Issue;
}

export async function updateIssue(spaceId: string, issueID: string, issueData: Partial<Issue>): Promise<Issue> {
  const userStore = useUserStore();
  if (!userStore.isAuthenticated) {
    throw new Error("User is not logged in");
  }

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/issues/${issueID}`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      },
      body: JSON.stringify(issueData),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to update issue: " + await response.text());
  }

  const issue = await response.json();
  issue.created_at = new Date(issue.created_at);
  issue.updated_at = new Date(issue.updated_at);
  if (issue.resolved_at) {
    issue.resolved_at = new Date(issue.resolved_at);
  }

  return issue as Issue;
}

export async function deleteIssue(spaceId: string, issueID: string): Promise<void> {
  const userStore = useUserStore();
  if (!userStore.isAuthenticated) {
    throw new Error("User is not logged in");
  }

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/issues/${issueID}`,
    {
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to delete issue: " + await response.text());
  }
}
