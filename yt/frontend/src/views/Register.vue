<template>
  <div class="max-w-md mx-auto">
    <div class="bg-gray-800 rounded-lg shadow-lg p-8">
      <h1 class="text-2xl font-bold text-center mb-8">Create Account</h1>
      
      <div v-if="error" class="bg-red-900 border border-red-700 text-red-100 px-4 py-3 rounded mb-4">
        {{ error }}
      </div>

      <form @submit.prevent="register" class="space-y-6">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-300 mb-2">Username</label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            required
            minlength="3"
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent"
            placeholder="Choose a username"
          >
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-gray-300 mb-2">Email</label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            required
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent"
            placeholder="Enter your email"
          >
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-300 mb-2">Password</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            required
            minlength="6"
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent"
            placeholder="Choose a password"
          >
        </div>

        <div>
          <label for="confirmPassword" class="block text-sm font-medium text-gray-300 mb-2">Confirm Password</label>
          <input
            id="confirmPassword"
            v-model="form.confirmPassword"
            type="password"
            required
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent"
            placeholder="Confirm your password"
          >
        </div>

        <button
          type="submit"
          :disabled="loading || !isFormValid"
          class="w-full bg-red-600 hover:bg-red-700 disabled:bg-red-800 text-white py-2 px-4 rounded-md font-medium transition-colors flex items-center justify-center"
        >
          <div v-if="loading" class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
          {{ loading ? 'Creating Account...' : 'Create Account' }}
        </button>
      </form>

      <div class="text-center mt-6">
        <p class="text-gray-400">
          Already have an account?
          <router-link to="/login" class="text-red-500 hover:text-red-400 font-medium">Sign in</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { authAPI } from '../services/api'
import { authService } from '../services/auth'

export default {
  name: 'Register',
  data() {
    return {
      form: {
        username: '',
        email: '',
        password: '',
        confirmPassword: ''
      },
      loading: false,
      error: null
    }
  },
  computed: {
    isFormValid() {
      return (
        this.form.username.length >= 3 &&
        this.form.email.includes('@') &&
        this.form.password.length >= 6 &&
        this.form.password === this.form.confirmPassword
      )
    }
  },
  methods: {
    async register() {
      if (!this.isFormValid) {
        this.error = 'Please fill all fields correctly and ensure passwords match.'
        return
      }

      try {
        this.loading = true
        this.error = null
        
        const { confirmPassword, ...userData } = this.form
        const response = await authAPI.register(userData)
        const { token, user } = response.data
        
        authService.login(token, user)
        
        this.$router.push('/')
      } catch (error) {
        console.error('Registration failed:', error)
        this.error = error.response?.data?.message || error.message || 'Registration failed. Please try again.'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>