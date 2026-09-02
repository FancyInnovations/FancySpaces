import type {SpaceMavenRepository, SpaceMavenRepositoryArtifact} from "@/api/maven/types";
import {useUserStore} from "@/stores/user";

export async function getAllMavenRepositories(spaceId: string): Promise<SpaceMavenRepository[]> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/maven-repositories`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch all maven repos: " + await response.text());
  }

  const repos = await response.json();
  repos.forEach((repo: SpaceMavenRepository) => {
    repo.created_at = new Date(repo.created_at);
  });

  return repos as SpaceMavenRepository[];
}

export async function getMavenRepository(spaceId: string, repoName: string): Promise<SpaceMavenRepository> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/maven-repositories/${repoName}`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch all maven repos: " + await response.text());
  }

  const repo = await response.json();
  repo.created_at = new Date(repo.created_at);

  return repo as SpaceMavenRepository;
}

export async function getAllMavenArtifacts(spaceId: string, repoName: string): Promise<SpaceMavenRepositoryArtifact[]> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/maven-repositories/${repoName}/artifacts`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch all maven repos: " + await response.text());
  }

  const artifacts = await response.json();
  if (!Array.isArray(artifacts)) {
    return [];
  }

  artifacts.forEach((artifact: SpaceMavenRepositoryArtifact) => {
    artifact.versions.forEach((version) => {
      version.published_at = new Date(version.published_at);
    });
  });

  return artifacts as SpaceMavenRepositoryArtifact[];
}

export async function getMavenArtifacts(spaceId: string, repoName: string, groupArtifactID: string): Promise<SpaceMavenRepositoryArtifact> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/maven-repositories/${repoName}/artifacts/${encodeURIComponent(groupArtifactID)}`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch all maven repos: " + await response.text());
  }

  const artifact = await response.json();
  artifact.versions.forEach((version: any) => {
    version.published_at = new Date(version.published_at);
  });

  return artifact as SpaceMavenRepositoryArtifact;
}

export async function deleteMavenArtifact(spaceId: string, repoName: string, groupArtifactID: string): Promise<void> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/maven-repositories/${repoName}/artifacts/${encodeURIComponent(groupArtifactID)}`,
    {
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to delete maven artifact: " + await response.text());
  }
}

export async function deleteMavenArtifactVersion(spaceId: string, repoName: string, groupArtifactID: string, version: string): Promise<void> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/maven-repositories/${repoName}/artifacts/${encodeURIComponent(groupArtifactID)}/versions/${encodeURIComponent(version)}`,
    {
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to delete maven artifact version: " + await response.text());
  }
}
