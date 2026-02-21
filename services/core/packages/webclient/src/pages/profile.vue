<script lang="ts" setup>

import {onMounted} from "vue";
import {useUserStore} from "@/stores/user.ts";
import {updateUser} from "@/api/auth/users.ts";
import {useNotificationStore} from "@/stores/notifications.ts";
import {useConfirmationStore} from "@/stores/confirmation.ts";
import type {User} from "@/api/auth/types.ts";
import {useHead} from "@vueuse/head";

const userStore = useUserStore();
const notifications = useNotificationStore();
const confirmationStore = useConfirmationStore();

const user = ref<User>(userStore.user!);
const activeSince = computed(() => {
    const date = new Date(userStore.user!.created_at);
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
});


const showEditUsernameDialog = ref(false);
const editedUsername = ref(userStore.user!.name);

const showEditEmailDialog = ref(false);
const editedEmail = ref(userStore.user!.email);

const showEditPasswordDialog = ref(false);
const editedPassword = ref('');
const editedPasswordConfirm = ref('');

onMounted(() => {
  useHead({
    title: `FancySpaces - Profile`,
    meta: [
      {
        name: 'description',
        content: 'View and manage your profile information..'
      }
    ]
  });
});

async function saveEditedUsername() {
    try {
        await updateUser(userStore.user!.id, editedUsername.value, "", "");
    } catch (error: any) {
        console.error(`Failed to update username: ${error.message}`);
        notifications.error(error.message);
    }

    userStore.user!.name = editedUsername.value;

    showEditUsernameDialog.value = false;
    notifications.info("Username updated successfully!");
}

function cancelEditUsername() {
    showEditUsernameDialog.value = false;
    editedUsername.value = userStore.user!.name;
}

async function saveEditedEmail() {
    try {
        await updateUser(userStore.user!.id, "", editedEmail.value, "");
    } catch (error: any) {
        console.error(`Failed to update email: ${error.message}`);
        notifications.error(error.message);
    }

    userStore.user!.email = editedEmail.value;

    showEditEmailDialog.value = false;
    notifications.info("E-Mail updated successfully!");
}

function cancelEditEmail() {
    showEditEmailDialog.value = false;
    editedEmail.value = userStore.user!.email;
}

async function saveEditedPassword() {
    if (editedPassword.value !== editedPasswordConfirm.value) {
        notifications.error("Passwords do not match!");
        return;
    }

    try {
        await updateUser(userStore.user!.id, "", "", editedPassword.value);
    } catch (error: any) {
        console.error(`Failed to update password: ${error.message}`);
        notifications.error(error.message);
    }

    showEditPasswordDialog.value = false;
    notifications.info("Password updated successfully!");
}

function cancelEditPassword() {
    showEditPasswordDialog.value = false;
    editedPassword.value = '';
    editedPasswordConfirm.value = '';
}

function deleteAccount() {
    confirmationStore.confirmation = {
        shown: true,
        persistent: true,
        title: "Delete Account",
        text: "Are you sure you want to delete your account? This action cannot be undone.",
        yesText: "Delete",
        onConfirm: () => {
            notifications.error("NOT IMPLEMENTED YET");
        }
    };
}

</script>

<template>
    <Dialog
        :shown="showEditUsernameDialog"
        persistent
    >
        <v-card>
            <v-card-title>Edit username</v-card-title>

            <v-card-text>
                <p class="text-body-1 mb-4">Please enter your new username:</p>

                <v-text-field
                    v-model="editedUsername"
                    autofocus
                    color="primary"
                    label="Username"
                    width="350px"
                />
            </v-card-text>

            <v-card-actions>
                <v-spacer/>

                <v-btn @click="saveEditedUsername">Save</v-btn>

                <v-btn @click="cancelEditUsername">Cancel</v-btn>
            </v-card-actions>
        </v-card>
    </Dialog>

    <Dialog
        :shown="showEditEmailDialog"
        persistent
    >
        <v-card>
            <v-card-title>Edit E-Mail</v-card-title>

            <v-card-text>
                <p class="text-body-1 mb-4">Please enter your new E-Mail:</p>

                <v-text-field
                    v-model="editedEmail"
                    autofocus
                    color="primary"
                    label="E-Mail"
                    width="350px"
                />
            </v-card-text>

            <v-card-actions>
                <v-spacer/>

                <v-btn @click="saveEditedEmail">Save</v-btn>

                <v-btn @click="cancelEditEmail">Cancel</v-btn>
            </v-card-actions>
        </v-card>
    </Dialog>

    <Dialog
        :shown="showEditPasswordDialog"
        persistent
    >
        <v-card>
            <v-card-title>Edit Password</v-card-title>

            <v-card-text>
                <p class="text-body-1 mb-4">Please enter your new password:</p>

                <v-text-field
                    v-model="editedPassword"
                    autofocus
                    color="primary"
                    label="New Password"
                    type="password"
                    width="350px"
                />

                <v-text-field
                    v-model="editedPasswordConfirm"
                    color="primary"
                    label="Confirm New Password"
                    type="password"
                    width="350px"
                />
            </v-card-text>

            <v-card-actions>
                <v-spacer/>

                <v-btn @click="saveEditedPassword">Save</v-btn>

                <v-btn @click="cancelEditPassword">Cancel</v-btn>
            </v-card-actions>
        </v-card>
    </Dialog>

    <v-container>
        <v-row justify="center">
            <v-col md="5">
                <h1 class="text-center">Your Profile</h1>
                <p class="text-body-1 text-center">Your account exists since {{ activeSince }}.</p>
                <p class="text-body-1 text-center">Roles: {{ user.roles.join(", ") }}</p>
            </v-col>
        </v-row>

        <v-row justify="center">
            <v-col md="5">
                <v-text-field
                    v-model="user.id"
                    color="primary"
                    label="User ID"
                    readonly
                />
            </v-col>
        </v-row>

        <v-row justify="center">
            <v-col md="5">
                <v-text-field
                    v-model="user.name"
                    append-inner-icon="mdi-pencil"
                    color="primary"
                    label="Username"
                    readonly
                    @click:append-inner="showEditUsernameDialog = true"
                />
            </v-col>
        </v-row>

        <v-row justify="center">
            <v-col md="5">
                <v-text-field
                    v-model="user.email"
                    append-inner-icon="mdi-pencil"
                    color="primary"
                    label="E-Mail"
                    readonly
                    @click:append-inner="showEditEmailDialog = true"
                />
            </v-col>
        </v-row>

        <v-row justify="center">
            <v-col md="5">
                <v-btn
                    class="mr-4"
                    color="primary"
                    @click="showEditPasswordDialog = true"
                >
                    Change Password
                </v-btn>

                <v-btn
                    class="ml-4"
                    color="error"
                    variant="outlined"
                    @click="deleteAccount"
                >
                    Delete Account
                </v-btn>
            </v-col>
        </v-row>
    </v-container>
</template>

<style scoped>

</style>
