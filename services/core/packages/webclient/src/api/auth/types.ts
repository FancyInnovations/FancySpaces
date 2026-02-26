// export const IDP_API_BASE_URL = "https://idp.fancyspaces.net/idp/api/v1";
export const IDP_API_BASE_URL = "http://localhost:8083/idp/api/v1";

export interface User {
    id: string;
    provider: string;
    name: string;
    email: string;
    verified: boolean;
    password: boolean;
    roles: string[];
    created_at: Date;
    is_active: boolean;
    metadata: Record<string, string>;
}

export interface ApiKey {
  key_id: string;
  user_id: string;
  description: string;
  created_at: Date;
  last_used_at: Date | null;
}
