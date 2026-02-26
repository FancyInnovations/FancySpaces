import type {Space} from "@/api/spaces/types.ts";
import {useUserStore} from "@/stores/user.ts";

export async function getAllSpaces(): Promise<Space[]> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
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

export async function getSpacesOfCreator(creatorID: string): Promise<Space[]> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces?creator=${encodeURIComponent(creatorID)}`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
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
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${id}`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
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
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/downloads`,
    {
      method: "GET",
      headers: {
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch download count for space: " + await response.text());
  }

  return (await response.json()).downloads as number;
}

export async function getDownloadCountForSpacePerVersion(spaceId: string): Promise<Record<string, number>> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/downloads`,
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

  return (await response.json()).versions as Record<string, number>;
}

export async function createSpace(slug: string, title: string, summary: string, description: string, categories: string[], iconURL: string): Promise<Space> {
  const userStore = useUserStore();
  if (!(await userStore.isAuthenticated)) {
    throw new Error("User is not logged in");
  }

  const response = await fetch(
    `/api/v1/spaces`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      },
      body: JSON.stringify({
        slug: slug,
        title: title,
        summary: summary,
        description: description,
        categories: categories,
        icon_url: iconURL,
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to create space: " + await response.text());
  }

  const space = await response.json();
  space.created_at = new Date(space.created_at);

  return space as Space;
}

export async function updateSpace(spaceID: string, slug: string, title: string, summary: string, description: string, categories: string[], iconURL: string): Promise<void> {
  const userStore = useUserStore();
  if (!(await userStore.isAuthenticated)) {
    throw new Error("User is not logged in");
  }

  const response = await fetch(
    `/api/v1/spaces/${spaceID}`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      },
      body: JSON.stringify({
        slug: slug,
        title: title,
        summary: summary,
        description: description,
        categories: categories,
        icon_url: iconURL,
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to update space: " + await response.text());
  }
}

export async function changeSpaceStatus(spaceID: string, toStatus: string): Promise<void> {
  const userStore = useUserStore();
  if (!(await userStore.isAuthenticated)) {
    throw new Error("User is not logged in");
  }

  const response = await fetch(
    `/api/v1/spaces/${spaceID}/status`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      },
      body: JSON.stringify({
        to: toStatus
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to change space status: " + await response.text());
  }
}
