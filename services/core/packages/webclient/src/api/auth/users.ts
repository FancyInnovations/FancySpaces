import type {User} from "@/api/auth/types.ts";
import {useUserStore} from "@/stores/user.ts";

export async function registerUser(username: string, email: string, password: string): Promise<void> {
    const resp = await fetch("/idp/api/v1/users/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Accept": "application/json",
        },
        body: JSON.stringify({
            provider: "basic",
            name: username,
            email: email,
            password: password,
        })
    });

    if (resp.status !== 201) {
        throw new Error(`Failed to register user (code ${resp.status} "${resp.statusText}"): ${await resp.text()}`);
    }
}

export async function validateUser(email: string, password: string): Promise<User> {
    const resp = await fetch("/idp/api/v1/users/validate", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Accept": "application/json",
        },
        body: JSON.stringify({
            user: email,
            password: password,
        })
    });

    if (resp.status !== 200) {
        throw new Error(`Failed to validate user (code ${resp.status} "${resp.statusText}"): ${await resp.text()}`);
    }

    let user: User;
    try {
        user = await resp.json();
    } catch (e) {
        throw new Error(`Failed to parse user data: ${e}`);
    }

    return user;
}

export async function updateUser(userid: string, name: string, email: string, password: string): Promise<void> {
    const userStore = useUserStore();
    if (!userStore.isAuthenticated) {
        throw new Error("User is not logged in");
    }

    const resp = await fetch(`/idp/api/v1/users/${userid}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Accept": "application/json",
            "Authorization": `Bearer ${userStore.token}`,
        },
        body: JSON.stringify({
            name: name,
            email: email,
            password: password,
        })
    });

    if (resp.status !== 200) {
        throw new Error(`Failed to update user (code ${resp.status} "${resp.statusText}"): ${await resp.text()}`);
    }
}

export async function verifyUser(code: string): Promise<void> {
    const userStore = useUserStore();
    if (!userStore.isAuthenticated) {
        throw new Error("User is not logged in");
    }

    const resp = await fetch(`/idp/api/v1/users/verify/check`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Accept": "application/json",
            "Authorization": `Bearer ${userStore.token}`,
        },
        body: code
    });

    if (resp.status !== 200) {
        throw new Error(`Failed to verify user (code ${resp.status} "${resp.statusText}"): ${await resp.text()}`);
    }
}

export async function resendVerificationCode(): Promise<void> {
    const userStore = useUserStore();
    if (!userStore.isAuthenticated) {
        throw new Error("User is not logged in");
    }

    const resp = await fetch(`/idp/api/v1/users/verify/resend`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Accept": "application/json",
            "Authorization": `Bearer ${userStore.token}`,
        },
        body: ""
    });

    if (resp.status !== 200) {
        throw new Error(`Failed to resend verification code (code ${resp.status} "${resp.statusText}"): ${await resp.text()}`);
    }
}
