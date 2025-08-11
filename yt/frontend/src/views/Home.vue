<template>
  <div>
    <h1 class="text-3xl font-bold mb-8">Latest Videos</h1>
    
    <div v-if="loading" class="text-center">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-red-500"></div>
    </div>

    <div v-else-if="error" class="bg-red-900 border border-red-700 text-red-100 px-4 py-3 rounded mb-4">
      {{ error }}
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <div v-for="video in videos" :key="video.id" class="bg-gray-800 rounded-lg overflow-hidden shadow-lg hover:shadow-xl transition-shadow">
        <router-link :to="`/video/${video.id}`" class="block">
          <div class="aspect-video bg-gray-700 flex items-center justify-center">
            <img v-if="video.thumbnail_url" :src="video.thumbnail_url" :alt="video.title" class="w-full h-full object-cover">
            <div v-else class="text-gray-500">
              <svg class="w-12 h-12" fill="currentColor" viewBox="0 0 20 20">
                <path d="M10 12a2 2 0 100-4 2 2 0 000 4z"/>
                <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd"/>
              </svg>
            </div>
          </div>
          <div class="p-4">
            <h2 class="text-lg font-semibold mb-2 line-clamp-2">{{ video.title }}</h2>
            <p class="text-gray-400 text-sm mb-2">by {{ video.username }}</p>
            <div class="flex justify-between text-xs text-gray-500">
              <span>{{ video.views }} views</span>
              <span>{{ formatDate(video.created_at) }}</span>
            </div>
          </div>
        </router-link>
      </div>
    </div>

    <div v-if="!loading && videos.length === 0" class="text-center py-12">
      <div class="text-gray-500 mb-4">
        <svg class="w-16 h-16 mx-auto" fill="currentColor" viewBox="0 0 20 20">
          <path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z"/>
        </svg>
      </div>
      <h2 class="text-xl font-semibold mb-2">No videos available</h2>
      <p class="text-gray-400 mb-4">Be the first to upload a video!</p>
      <router-link v-if="isAuthenticated" to="/upload" class="bg-red-600 hover:bg-red-700 text-white px-6 py-2 rounded-md font-medium">Upload Video</router-link>
    </div>
  </div>
</template>

<script>
import { videoAPI } from '../services/api'
import { authService } from '../services/auth'

export default {
  name: 'Home',
  data() {
    return {
      videos: [],
      loading: true,
      error: null
    }
  },
  computed: {
    isAuthenticated() {
      return authService.isAuthenticated()
    }
  },
  async mounted() {
    await this.fetchVideos()
  },
  methods: {
    async fetchVideos() {
      try {
        this.loading = true
        this.error = null
        const response = await videoAPI.getAll()
        this.videos = response.data || []
      } catch (error) {
        console.error('Failed to fetch videos:', error)
        this.error = 'Failed to load videos. Please try again later.'
      } finally {
        this.loading = false
      }
    },
    formatDate(dateString) {
      const date = new Date(dateString)
      const now = new Date()
      const diff = Math.floor((now - date) / 1000)
      
      if (diff < 60) return 'Just now'
      if (diff < 3600) return `${Math.floor(diff / 60)} minutes ago`
      if (diff < 86400) return `${Math.floor(diff / 3600)} hours ago`
      if (diff < 2592000) return `${Math.floor(diff / 86400)} days ago`
      
      return date.toLocaleDateString()
    }
  }
}
</script>