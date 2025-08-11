<template>
  <div id="app" class="min-h-screen bg-gray-900 text-white">
    <nav class="bg-gray-800 shadow-lg">
      <div class="max-w-7xl mx-auto px-4">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link to="/" class="flex items-center">
              <h1 class="text-xl font-bold text-red-500">YouTube Clone</h1>
            </router-link>
          </div>
          <div class="flex items-center space-x-4">
            <template v-if="isAuthenticated">
              <router-link to="/my-videos" class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium">My Videos</router-link>
              <router-link to="/upload" class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium">Upload</router-link>
              <button @click="logout" class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium">Logout</button>
            </template>
            <template v-else>
              <router-link to="/login" class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium">Login</router-link>
              <router-link to="/register" class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium">Register</router-link>
            </template>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto py-6 px-4">
      <router-view />
    </main>
  </div>
</template>

<script>
import { authService } from './services/auth'

export default {
  name: 'App',
  computed: {
    isAuthenticated() {
      return authService.isAuthenticated()
    }
  },
  methods: {
    logout() {
      authService.logout()
      this.$router.push('/')
    }
  }
}
</script>