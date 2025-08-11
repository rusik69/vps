<template>
  <div class="max-w-md mx-auto">
    <div class="bg-gray-800 rounded-lg shadow-lg p-8">
      <h1 class="text-2xl font-bold text-center mb-8">Login to YouTube Clone</h1>
      
      <div v-if="error" class="bg-red-900 border border-red-700 text-red-100 px-4 py-3 rounded mb-4">
        {{ error }}
      </div>

      <form @submit.prevent="login" class="space-y-6">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-300 mb-2">Username</label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            required
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent"
            placeholder="Enter your username"
          >
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-300 mb-2">Password</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            required
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent"
            placeholder="Enter your password"
          >
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-red-600 hover:bg-red-700 disabled:bg-red-800 text-white py-2 px-4 rounded-md font-medium transition-colors flex items-center justify-center"
        >
          <div v-if="loading" class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
          {{ loading ? 'Signing in...' : 'Sign In' }}
        </button>
      </form>

      <div class="text-center mt-6">
        <p class="text-gray-400">
          Don't have an account?
          <router-link to="/register" class="text-red-500 hover:text-red-400 font-medium">Sign up</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { authAPI } from '../services/api'
import { authService } from '../services/auth'

export default {
  name: 'Login',
  data() {
    return {
      form: {
        username: '',
        password: ''
      },
      loading: false,
      error: null
    }
  },
  methods: {
    async login() {
      try {
        this.loading = true
        this.error = null
        
        const response = await authAPI.login(this.form)
        const { token, user } = response.data
        
        authService.login(token, user)
        
        this.$router.push('/')
      } catch (error) {
        console.error('Login failed:', error)
        this.error = error.response?.data?.message || error.message || 'Login failed. Please try again.'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>