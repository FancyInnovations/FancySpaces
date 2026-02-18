import type {SpaceSecret} from "@/api/secrets/types.ts";

export async function getSecret(spaceId: string, key: string): Promise<SpaceSecret> {
  const response = await fetch(
    `/api/v1/spaces/${spaceId}/secrets/${key}`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch secret: " + await response.text());
  }

  const secret = await response.json();
  secret.created_at = new Date(secret.created_at);
  secret.updated_at = new Date(secret.updated_at);

  return secret as SpaceSecret;
}

export async function getSecretDecrypted(spaceId: string, key: string): Promise<string> {
  const response = await fetch(
    `/api/v1/spaces/${spaceId}/secrets/${key}/decrypted`,
    {
      method: "GET",
      headers: {
        "Accept": "text/plain",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch decrypted secret: " + await response.text());
  }

  return await response.text();
}

export async function getAllSecrets(spaceId: string): Promise<SpaceSecret[]> {
  const response = await fetch(
    `/api/v1/spaces/${spaceId}/secrets`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch secrets: " + await response.text());
  }

  const secrets = await response.json();
  return secrets.map((secret: any) => ({
    ...secret,
    created_at: new Date(secret.created_at),
    updated_at: new Date(secret.updated_at),
  })) as SpaceSecret[];
}

export async function createSecret(spaceId: string, key: string, value: string, description: string): Promise<void> {
  const response = await fetch(
    `/api/v1/spaces/${spaceId}/secrets`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      },
      body: JSON.stringify({
        key: key,
        value: value,
        description: description,
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to create secret: " + await response.text());
  }
}

export async function updateSecret(spaceId: string, key: string, value: string, description: string): Promise<void> {
  const response = await fetch(
    `/api/v1/spaces/${spaceId}/secrets/${key}`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      },
      body: JSON.stringify({
        value: value,
        description: description,
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to update secret: " + await response.text());
  }
}

export async function deleteSecret(spaceId: string, key: string): Promise<void> {
  const response = await fetch(
    `/api/v1/spaces/${spaceId}/secrets/${key}`,
    {
      method: "DELETE",
      headers: {
        "Authorization": localStorage.getItem("fs_api_key") || "",
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to delete secret: " + await response.text());
  }
}
