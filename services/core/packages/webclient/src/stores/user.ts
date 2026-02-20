import type {User} from "@/api/auth/types.ts";
import {validateToken} from "@/api/auth/tokens.ts";
import {jwtDecode} from "jwt-decode";

export const useUserStore = defineStore('user', {
    state: () => ({
        user: null as User | null,
        token: null as string | null,
    }),

    getters: {
        isAuthenticated: async (state) => {
            // Check if token exists
            if (!state.token) {
                return false;
            }

            // Check if the token is expired
            if (getTokenTTL(state.token) <= 0) {
                state.user = null;
                state.token = null;
                return false;
            }

            // Validate the token with the backend
            const valid = await validateToken(state.token);
            if (!valid) {
                state.user = null;
                state.token = null;
                return false;
            }

            return true;
        },
        tokenTTL: (state) => {
            // Returns the time-to-live of the token in milliseconds
            return getTokenTTL(state.token);
        }
    },

    actions: {
        setUser(user: User) {
            this.user = user;
            localStorage.setItem("current_user", JSON.stringify(user));
        },

        loadUserFromStorage() {
            const userData = localStorage.getItem("current_user");
            if (userData) {
                try {
                    this.user = JSON.parse(userData) as User;
                } catch (e) {
                    console.error("Failed to parse user data from storage:", e);
                    this.user = null;
                }
            } else {
                this.user = null;
            }
        },

        clearUser() {
            this.user = null;
            localStorage.removeItem("current_user");

            this.token = null;
            localStorage.removeItem("auth_token");
        },

        setToken(token: string) {
            this.token = token;
            localStorage.setItem("auth_token", token);
        },

        loadTokenFromStorage() {
            const token = localStorage.getItem("auth_token");
            if (token) {
                this.token = token;
            } else {
                this.token = null;
            }
        },

        clearToken() {
            this.token = null;
        }
    }
});

function getTokenTTL(token: string | null): number {
    if (!token) {
        return 0;
    }

    const payload = jwtDecode(token);
    if (typeof payload !== 'object' || !payload.exp) {
        return 0; // Invalid token or no expiration
    }

    const expMs = payload.exp * 1000; // Convert to milliseconds
    const nowMs = Date.now();

    return Math.max(0, expMs - nowMs); // Return remaining time in milliseconds
}
