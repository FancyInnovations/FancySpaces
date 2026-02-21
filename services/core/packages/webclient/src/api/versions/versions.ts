import type {SpaceVersion} from "@/api/versions/types.ts";
import {useUserStore} from "@/stores/user.ts";

export async function getVersion(spaceId: string, versionId: string): Promise<SpaceVersion> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/versions/${versionId}`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch version: " + await response.text());
  }

  const ver = await response.json();
  ver.published_at = new Date(ver.published_at);

  return ver as SpaceVersion;
}

export async function getLatestVersion(spaceId: string): Promise<SpaceVersion> {
  return getVersion(spaceId, "latest");
}

export async function getAllVersions(spaceId: string): Promise<SpaceVersion[]> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/versions`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch latest version: " + await response.text());
  }

  const versions = await response.json();
  versions.forEach((ver: SpaceVersion) => {
    ver.published_at = new Date(ver.published_at);
  });

  return versions as SpaceVersion[];
}

export async function getDownloadCountForVersion(spaceId: string, versionId: string): Promise<number> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/versions/${versionId}/downloads`,
    {
      method: "GET",
      headers: {
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch download count for version: " + await response.text());
  }

  return (await response.json()).downloads as number;
}
