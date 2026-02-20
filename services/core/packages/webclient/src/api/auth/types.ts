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
