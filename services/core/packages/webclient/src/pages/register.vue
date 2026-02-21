<script lang="ts" setup>

import {onMounted} from "vue";
import {useUserStore} from "@/stores/user.ts";
import {registerUser, validateUser} from "@/api/auth/users.ts";
import {createToken, validateToken} from "@/api/auth/tokens.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import router from "@/router";
import {useHead} from "@vueuse/head";

const userStore = useUserStore();
const notifications = useNotificationStore();

const username = ref('');
const email = ref('');
const password = ref('');
const showPassword = ref(false);
const repeatedPassword = ref('');
const showRepeatedPassword = ref(false);

const usernameRule = (value: string) => {
    if (!value) return 'Username is required';

    if (value.length < 5) return 'Username must be at least 5 characters long';

    if (value.length > 25) return 'Username must not exceed 25 characters';
    return true;
};

const emailRule = (value: string) => {
    if (!value) return 'E-Mail is required';

    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailPattern.test(value)) return 'Invalid E-Mail format';

    return true;
};

const passwordRule = (value: string) => {
    if (!value) return 'Password is required';

    if (value.length < 16) return 'Password must be at least 16 characters long';

    if (value.length > 64) return 'Password must not exceed 64 characters';

    if (!/[A-Z]/.test(value)) return 'Password must contain at least one uppercase letter';

    if (!/[a-z]/.test(value)) return 'Password must contain at least one lowercase letter';

    if (!/[0-9]/.test(value)) return 'Password must contain at least one number';

    if (!/[!@#$%^&*(),.?":{}|<>]/.test(value)) return 'Password must contain at least one special character';

    return true;
};

const repeatedPasswordRule = (value: string) => {
    if (!value) return 'Please repeat your password';

    if (value !== password.value) return 'Passwords do not match';

    return true;
};

const isEverythingValid = computed(() => {
    return usernameRule(username.value) === true &&
        emailRule(email.value) === true &&
        passwordRule(password.value) === true &&
        repeatedPasswordRule(repeatedPassword.value) === true;
});

onMounted(async () => {
  useHead({
    title: `FancySpaces - Register`,
    meta: [
      {
        name: 'description',
        content: 'Create a new account on FancySpaces to manage your projects and collaborate with your team.'
      }
    ]
  });

    if (await userStore.isAuthenticated) {
        await router.push("/");
    }
});

async function register() {
    if (!isEverythingValid.value) {
        return;
    }

    // create user
    try {
        await registerUser(username.value, email.value, password.value);
    } catch (error: any) {
        console.error(`Registration failed: ${error.message}`);
        notifications.error(error.message);
        return;
    }

    // get user info
    let user;
    try {
        user = await validateUser(email.value, password.value);
    } catch (error: any) {
        console.error(`User validation failed: ${error.message}`);
        notifications.error(error.message);
        return;
    }
    userStore.setUser(user);

    // create token
    let token: string;
    try {
        token = await createToken(userStore.user!.id, password.value);
    } catch (error: any) {
        console.error(`Token creation failed: ${error.message}`);
        notifications.error(error.message);
        return;
    }
    userStore.setToken(token);

    // validate token
    try {
        const valid = validateToken(userStore.token!);
        if (!valid) {
            console.error('Token is invalid');
            notifications.error('Token is invalid');
            return;
        }
    } catch (error: any) {
        console.error(`Token validation failed: ${error.message}`);
        notifications.error(error.message);
        return;
    }

    notifications.info("Registration successful!");
    await router.push("/");
}

</script>

<template>
    <v-container>
        <v-row justify="center">
            <v-col md="4">
                <h1 class="text-center">Sign up with</h1>
            </v-col>
        </v-row>

        <v-row justify="center">
            <v-col class="d-flex justify-space-evenly" md="4">
                <v-btn
                    color="primary"
                    variant="outlined"
                >
                    <v-icon left>mdi-google</v-icon>
                    Google
                </v-btn>

                <v-btn
                    class="mx-2"
                    color="primary"
                    variant="outlined"
                >
                    <v-icon left>mdi-github</v-icon>
                    GitHub
                </v-btn>

                <v-btn
                    color="primary"
                    variant="outlined"
                >
                    <v-icon left>mdi-chat</v-icon>
                    Discord
                </v-btn>
            </v-col>
        </v-row>

        <v-row justify="center">
            <v-col md="4">
                <v-divider class="my-4"/>
            </v-col>
        </v-row>

        <v-row justify="center">
            <v-col md="4">
                <h1 class="text-center">Or create an account yourself</h1>
                <p class="text-center">Already got an account? <a href="/login">Sign in here</a>.</p>
            </v-col>
        </v-row>

        <v-row justify="center">
            <v-col md="2">
                <v-text-field
                    v-model="username"
                    :rules="[usernameRule]"
                    autofocus
                    color="primary"
                    label="Username"
                />
            </v-col>

            <v-col md="2">
                <v-text-field
                    v-model="email"
                    :rules="[emailRule]"
                    color="primary"
                    label="E-Mail"
                />
            </v-col>
        </v-row>

        <v-row justify="center">
            <v-col md="4">
                <v-text-field
                    v-model="password"
                    :append-inner-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                    :rules="[passwordRule]"
                    :type="showPassword ? 'text' : 'password'"
                    color="primary"
                    label="Password"
                    @click:append-inner="showPassword = !showPassword"
                />
            </v-col>
        </v-row>

        <v-row justify="center">
            <v-col md="4">
                <v-text-field
                    v-model="repeatedPassword"
                    :append-inner-icon="showRepeatedPassword ? 'mdi-eye' : 'mdi-eye-off'"
                    :rules="[repeatedPasswordRule]"
                    :type="showRepeatedPassword ? 'text' : 'password'"
                    color="primary"
                    label="Repeat Password"
                    @click:append-inner="showRepeatedPassword = !showRepeatedPassword"
                />
            </v-col>
        </v-row>

        <v-row justify="center">
            <v-col md="4">
                <p>By creating an account, you agree to FancySpaces' <a href="">Terms</a> and <a href="">Privacy Policy</a>.</p>
            </v-col>
        </v-row>

        <v-row justify="center">
            <v-col md="4">
                <v-btn
                    :disabled="!isEverythingValid"
                    color="primary"
                    @click="register"
                >
                    Register
                </v-btn>
            </v-col>
        </v-row>
    </v-container>
</template>

<style scoped>
a {
  text-decoration: underline;
}
</style>
