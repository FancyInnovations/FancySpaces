<script lang="ts" setup>

  import { useHead } from '@vueuse/head'
  import { onMounted } from 'vue'
  import { createToken, validateToken } from '@/api/auth/tokens'
  import { validateUser } from '@/api/auth/users'
  import router from '@/router'
  import { useNotificationStore } from '@/stores/notifications'
  import { useUserStore } from '@/stores/user'

  const userStore = useUserStore()
  const notifications = useNotificationStore()

  const email = ref('')
  const password = ref('')
  const showPassword = ref(false)

  const showForgotPasswordDialog = ref(false)
  const resetPasswordEmail = ref('')

  function emailRule (value: string) {
    if (!value) return 'E-Mail is required'
    return true
  }

  function passwordRule (value: string) {
    if (!value) return 'Password is required'
    return true
  }

  const isEverythingValid = computed(() => {
    return emailRule(email.value) === true
      && passwordRule(password.value) === true
  })

  onMounted(async () => {
    useHead({
      title: `FancySpaces - Login`,
      meta: [
        {
          name: 'description',
          content: 'Login to your FancySpaces account to access your projects.',
        },
      ],
    })

    if (await userStore.isAuthenticated) {
      await router.push('/')
    }
  })

  async function login () {
    if (!isEverythingValid.value) {
      return
    }

    // Get user
    let user
    try {
      user = await validateUser(email.value, password.value)
    } catch (error: any) {
      console.error(`User validation failed: ${error.message}`)
      notifications.error(error.message)
      return
    }
    userStore.setUser(user)

    // create token
    let token: string
    try {
      token = await createToken(userStore.user!.id, password.value)
    } catch (error: any) {
      console.error(`Token creation failed: ${error.message}`)
      notifications.error(error.message)
      return
    }
    userStore.setToken(token)

    // validate token
    try {
      const valid = validateToken(userStore.token!)
      if (!valid) {
        console.error('Token is invalid')
        notifications.error('Token is invalid')
        return
      }
    } catch (error: any) {
      console.error(`Token validation failed: ${error.message}`)
      notifications.error(error.message)
      return
    }

    notifications.info('Login successful!')
    window.location.href = '/users/' + user.name
  }

</script>

<template>
  <v-container>
    <v-row justify="center">
      <v-col md="4">
        <h1 class="text-center">Sign in with</h1>
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
        <v-divider class="my-4" />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="4">
        <h1 class="text-center">Or login with your account</h1>

        <p class="text-center">
          Don't have an account? <a href="/register">Register here</a>.
        </p>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="4">
        <v-text-field
          v-model="email"
          autofocus
          color="primary"
          label="E-Mail"
          :rules="[emailRule]"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="4">
        <v-text-field
          v-model="password"
          :append-inner-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
          color="primary"
          label="Password"
          :rules="[passwordRule]"
          :type="showPassword ? 'text' : 'password'"
          @click:append-inner="showPassword = !showPassword"
        />
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="4">
        <v-btn
          color="primary"
          :disabled="!isEverythingValid"
          @click="login"
        >
          Login
        </v-btn>

        <v-btn
          class="ml-4"
          disabled
          variant="text"
          @click="showForgotPasswordDialog = true"
        >
          Reset Password
        </v-btn>
      </v-col>
    </v-row>
  </v-container>

  <!-- Reset password dialog  -->

  <Dialog
    :persistent="true"
    :shown="showForgotPasswordDialog"
  >
    <v-card
      elevation="8"
      rounded="xl"
    >
      <v-card-title class="mx-2 mt-2">
        Reset Password
      </v-card-title>

      <v-card-text>
        <v-container>
          <v-row>
            <v-col>
              <p>
                Please enter your email address to reset your password.
              </p>
            </v-col>
          </v-row>

          <v-row>
            <v-col>
              <v-text-field
                v-model="resetPasswordEmail"
                color="primary"
                label="E-Mail"
                :rules="[emailRule]"
              />
            </v-col>
          </v-row>

          <v-row justify="center">
            <v-col class="text-center">
              <v-btn
                color="primary"
                :disabled="!emailRule(email)"
              >
                Send Reset Link
              </v-btn>
            </v-col>
          </v-row>
        </v-container>
      </v-card-text>

      <v-card-actions>
        <v-spacer />

        <v-btn @click="showForgotPasswordDialog = false">Close</v-btn>
      </v-card-actions>
    </v-card>
  </Dialog>
</template>

<style scoped>
a {
  text-decoration: underline;
}
</style>
