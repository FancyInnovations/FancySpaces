import type { MavenRepositoryMutation, SpaceMavenRepository, SpaceMavenRepositoryArtifact } from '@/api/maven/types'
import { useUserStore } from '@/stores/user'

export async function getAllMavenRepositories (spaceId: string): Promise<SpaceMavenRepository[]> {
  const userStore = useUserStore()

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/maven-repositories`,
    {
      method: 'GET',
      headers: {
        Accept: 'application/json',
        Authorization: `Bearer ${userStore.token}`,
      },
    },
  )

  if (!response.ok) {
    throw new Error('Failed to fetch all maven repos: ' + await response.text())
  }

  const repos = await response.json()
  for (const repo of repos as SpaceMavenRepository[]) {
    repo.created_at = new Date(repo.created_at)
  }

  return repos as SpaceMavenRepository[]
}

export async function getMavenRepository (spaceId: string, repoName: string): Promise<SpaceMavenRepository> {
  const userStore = useUserStore()

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/maven-repositories/${repoName}`,
    {
      method: 'GET',
      headers: {
        Accept: 'application/json',
        Authorization: `Bearer ${userStore.token}`,
      },
    },
  )

  if (!response.ok) {
    throw new Error('Failed to fetch all maven repos: ' + await response.text())
  }

  const repo = await response.json()
  repo.created_at = new Date(repo.created_at)

  return repo as SpaceMavenRepository
}

export async function createMavenRepository (spaceId: string, data: MavenRepositoryMutation): Promise<SpaceMavenRepository> {
  const userStore = useUserStore()
  const response = await fetch(`/api/v1/spaces/${spaceId}/maven-repositories`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      'Authorization': `Bearer ${userStore.token}`,
    },
    body: JSON.stringify(data),
  })

  if (!response.ok) {
    throw new Error('Failed to create maven repository: ' + await response.text())
  }

  const repo = await response.json()
  repo.created_at = new Date(repo.created_at)
  return repo as SpaceMavenRepository
}

export async function updateMavenRepository (spaceId: string, repoName: string, data: MavenRepositoryMutation): Promise<SpaceMavenRepository> {
  const userStore = useUserStore()
  const response = await fetch(`/api/v1/spaces/${spaceId}/maven-repositories/${encodeURIComponent(repoName)}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      'Authorization': `Bearer ${userStore.token}`,
    },
    body: JSON.stringify(data),
  })

  if (!response.ok) {
    throw new Error('Failed to update maven repository: ' + await response.text())
  }

  const repo = await response.json()
  repo.created_at = new Date(repo.created_at)
  return repo as SpaceMavenRepository
}

export async function deleteMavenRepository (spaceId: string, repoName: string): Promise<void> {
  const userStore = useUserStore()
  const response = await fetch(`/api/v1/spaces/${spaceId}/maven-repositories/${encodeURIComponent(repoName)}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${userStore.token}` },
  })

  if (!response.ok) {
    throw new Error('Failed to delete maven repository: ' + await response.text())
  }
}

export async function getAllMavenArtifacts (spaceId: string, repoName: string): Promise<SpaceMavenRepositoryArtifact[]> {
  const userStore = useUserStore()

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/maven-repositories/${repoName}/artifacts`,
    {
      method: 'GET',
      headers: {
        Accept: 'application/json',
        Authorization: `Bearer ${userStore.token}`,
      },
    },
  )

  if (!response.ok) {
    throw new Error('Failed to fetch all maven repos: ' + await response.text())
  }

  const artifacts = await response.json()
  if (!Array.isArray(artifacts)) {
    return []
  }

  for (const artifact of artifacts as SpaceMavenRepositoryArtifact[]) {
    for (const version of artifact.versions) {
      version.published_at = new Date(version.published_at)
    }
  }

  return artifacts as SpaceMavenRepositoryArtifact[]
}

export async function getMavenArtifacts (spaceId: string, repoName: string, groupArtifactID: string): Promise<SpaceMavenRepositoryArtifact> {
  const userStore = useUserStore()

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/maven-repositories/${repoName}/artifacts/${encodeURIComponent(groupArtifactID)}`,
    {
      method: 'GET',
      headers: {
        Accept: 'application/json',
        Authorization: `Bearer ${userStore.token}`,
      },
    },
  )

  if (!response.ok) {
    throw new Error('Failed to fetch all maven repos: ' + await response.text())
  }

  const artifact = await response.json()
  for (const version of artifact.versions as any[]) {
    version.published_at = new Date(version.published_at)
  }

  return artifact as SpaceMavenRepositoryArtifact
}

export async function deleteMavenArtifact (spaceId: string, repoName: string, groupArtifactID: string): Promise<void> {
  const userStore = useUserStore()

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/maven-repositories/${repoName}/artifacts/${encodeURIComponent(groupArtifactID)}`,
    {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${userStore.token}`,
      },
    },
  )

  if (!response.ok) {
    throw new Error('Failed to delete maven artifact: ' + await response.text())
  }
}

export async function deleteMavenArtifactVersion (spaceId: string, repoName: string, groupArtifactID: string, version: string): Promise<void> {
  const userStore = useUserStore()

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/maven-repositories/${repoName}/artifacts/${encodeURIComponent(groupArtifactID)}/versions/${encodeURIComponent(version)}`,
    {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${userStore.token}`,
      },
    },
  )

  if (!response.ok) {
    throw new Error('Failed to delete maven artifact version: ' + await response.text())
  }
}
