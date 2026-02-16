export async function kvDelete(db: string, coll: string, key: string): Promise<void> {
  const response = await fetch(
    `/storage/api/v1/databases/${db}/collections/${coll}/kv/2020`,
    {
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      },
      body: JSON.stringify({
        key: key
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to delete key: " + await response.text());
  }
}

export async function kvDeleteMultiple(db: string, coll: string, keys: string[]): Promise<void> {
  const response = await fetch(
    `/storage/api/v1/databases/${db}/collections/${coll}/kv/2021`,
    {
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      },
      body: JSON.stringify({
        keys: keys
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to delete multiple keys: " + await response.text());
  }
}

export async function kvDeleteAll(db: string, coll: string): Promise<void> {
  const response = await fetch(
    `/storage/api/v1/databases/${db}/collections/${coll}/kv/2022`,
    {
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      },
    },
  );

  if (!response.ok) {
    throw new Error("Failed to delete all keys: " + await response.text());
  }
}

export async function kvExists(db: string, coll: string, key: string): Promise<boolean> {
  const response = await fetch(
    `/storage/api/v1/databases/${db}/collections/${coll}/kv/2030`,
    {
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      },
      body: JSON.stringify({
        key: key
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to check if key exists: " + await response.text());
  }

  const data = await response.json();

  return data.exists;
}

export async function kvGet(db: string, coll: string, key: string): Promise<any> {
  const response = await fetch(
    `/storage/api/v1/databases/${db}/collections/${coll}/kv/2031`,
    {
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      },
      body: JSON.stringify({
        key: key
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to get key: " + await response.text());
  }

  const data = await response.json();

  return data.value;
}

export async function kvGetMultiple(db: string, coll: string, keys: string[]): Promise<Record<string, any>> {
  const response = await fetch(
    `/storage/api/v1/databases/${db}/collections/${coll}/kv/2031`,
    {
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      },
      body: JSON.stringify({
        keys: keys
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to get multiple key: " + await response.text());
  }

  const data = await response.json();

  return data.values;
}

export async function kvGetAll(db: string, coll: string): Promise<Record<string, any>> {
  const response = await fetch(
    `/storage/api/v1/databases/${db}/collections/${coll}/kv/2032`,
    {
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      },
    },
  );

  if (!response.ok) {
    throw new Error("Failed to get all keys: " + await response.text());
  }

  const data = await response.json();

  return data.values;
}

export async function kvKeys(db: string, coll: string): Promise<string[]> {
  const response = await fetch(
    `/storage/api/v1/databases/${db}/collections/${coll}/kv/2037`,
    {
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch kv keys: " + await response.text());
  }

  const data = await response.json();

  return data.keys;
}

export async function kvCount(db: string, coll: string): Promise<number> {
  const response = await fetch(
    `/storage/api/v1/databases/${db}/collections/${coll}/kv/2038`,
    {
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch kv count: " + await response.text());
  }

  const data = await response.json();

  return data.count;
}

export async function kvSize(db: string, coll: string): Promise<number> {
  const response = await fetch(
    `/storage/api/v1/databases/${db}/collections/${coll}/kv/2039`,
    {
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Authorization": localStorage.getItem("fs_api_key") || "",
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch kv size: " + await response.text());
  }

  const data = await response.json();

  return data.size;
}
