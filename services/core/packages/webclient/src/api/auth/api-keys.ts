import {useUserStore} from "@/stores/user.ts";
import {type ApiKey, IDP_API_BASE_URL} from "@/api/auth/types.ts";

export async function getApiKeys(userid: string): Promise<ApiKey[]> {
  const userStore = useUserStore();
  if (!(await userStore.isAuthenticated)) {
    throw new Error("User is not logged in");
  }

  const resp = await fetch(`${IDP_API_BASE_URL}/api-keys`, {
    method: "GET",
    headers: {
      "Accept": "application/json",
      "Authorization": `Bearer ${userStore.token}`,
    },
  });

  if (resp.status !== 200) {
    throw new Error(`Failed to get api keys (code ${resp.status} "${resp.statusText}"): ${await resp.text()}`);
  }

  const apiKeys = await resp.json();

  // replace created_at and last_used_at with Date objects
  for (const apiKey of apiKeys) {
    apiKey.created_at = new Date(apiKey.created_at);
    if (apiKey.last_used_at) {
      apiKey.last_used_at = new Date(apiKey.last_used_at);
    }
  }

  return apiKeys as ApiKey[];
}

export async function createApiKey(description: string): Promise<string> {
  const userStore = useUserStore();
  if (!(await userStore.isAuthenticated)) {
    throw new Error("User is not logged in");
  }

  const resp = await fetch(`${IDP_API_BASE_URL}/api-keys`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Accept": "application/json",
      "Authorization": `Bearer ${userStore.token}`,
    },
    body: JSON.stringify({
      description: description,
    }),
  });

  if (resp.status !== 200) {
    throw new Error(`Failed to create api key (code ${resp.status} "${resp.statusText}"): ${await resp.text()}`);
  }

  return resp.text()
}

export async function deleteApiKey(keyId: string): Promise<void> {
  const userStore = useUserStore();
  if (!(await userStore.isAuthenticated)) {
    throw new Error("User is not logged in");
  }

  const resp = await fetch(`${IDP_API_BASE_URL}/api-keys/${keyId}`, {
    method: "DELETE",
    headers: {
      "Accept": "application/json",
      "Authorization": `Bearer ${userStore.token}`,
    },
  });

  if (resp.status !== 204) {
    throw new Error(`Failed to delete api key (code ${resp.status} "${resp.statusText}"): ${await resp.text()}`);
  }
}
