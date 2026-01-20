import type {Space} from "@/api/spaces/types.ts";

export async function getAllSpaces(): Promise<Space[]> {
  const response = await fetch(
    `/api/v1/spaces`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch all spaces: " + await response.text());
  }

  const spaces = await response.json();
  spaces.forEach((space: Space) => {
    space.created_at = new Date(space.created_at);
  });

  return spaces as Space[];
}


export async function getSpace(id: string): Promise<Space> {
  const response = await fetch(
    `/api/v1/spaces/${id}`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch space: " + await response.text());
  }

  const space = await response.json();
  space.created_at = new Date(space.created_at);

  return space as Space;
}

export async function getDownloadCountForSpace(spaceId: string): Promise<number> {
  const response = await fetch(
    `/api/v1/spaces/${spaceId}/downloads`,
    {
      method: "GET",
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch download count for space: " + await response.text());
  }

  return (await response.json()).downloads as number;
}

export async function getDownloadCountForSpacePerVersion(spaceId: string): Promise<Record<string, number>> {
  const response = await fetch(
    `/api/v1/spaces/${spaceId}/downloads`,
    {
      method: "GET",
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch download count for version: " + await response.text());
  }

  return (await response.json()).versions as Record<string, number>;
}

