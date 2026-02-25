export async function createToken(email: string, password: string): Promise<string> {
    const resp = await fetch("https://idp.fancyspaces.net/idp/api/v1/tokens/create", {
        method: "POST",
        headers: {
            "Accept": "text/plain",
            "Authorization": `Basic ${btoa(email + ":" + password)}`,
        },
    });

    if (resp.status !== 201) {
        throw new Error(`Failed to create token (code ${resp.status} "${resp.statusText}"): ${await resp.text()}`);
    }

    const token = resp.text();
    if (!token) {
        throw new Error("Failed to create token: No token returned");
    }

    return token;
}

export async function validateToken(token: string): Promise<boolean> {
    const resp = await fetch("https://idp.fancyspaces.net/idp/api/v1/tokens/validate", {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${token}`,
        },
    });

    if (!resp.ok) {
        // TODO handle specific error cases
        return false;
    }

    return true;
}

export async function refreshToken(token: string): Promise<string> {
    const resp = await fetch("https://idp.fancyspaces.net/idp/api/v1/tokens/refresh", {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${token}`,
        },
    });

    if (resp.status !== 200) {
        throw new Error(`Failed to refresh token (code ${resp.status} "${resp.statusText}"): ${await resp.text()}`);
    }

    const newToken = resp.text();
    if (!newToken) {
        throw new Error("Failed to refresh token: No token returned");
    }

    return newToken;
}
