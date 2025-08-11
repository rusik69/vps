<template>
  <div v-if="loading" class="text-center">
    <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-red-500"></div>
  </div>

  <div v-else-if="error" class="bg-red-900 border border-red-700 text-red-100 px-4 py-3 rounded mb-4">
    {{ error }}
  </div>

  <div v-else-if="video" class="max-w-6xl mx-auto">
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Video Player Section -->
      <div class="lg:col-span-2">
        <div class="bg-gray-800 rounded-lg overflow-hidden mb-6">
          <div class="aspect-video bg-gray-900 flex items-center justify-center">
            <video
              v-if="video.url"
              :src="video.url"
              controls
              class="w-full h-full"
              :poster="video.thumbnail_url"
            >
              Your browser does not support the video tag.
            </video>
            <div v-else class="text-gray-500">
              <svg class="w-16 h-16" fill="currentColor" viewBox="0 0 20 20">
                <path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z"/>
              </svg>
            </div>
          </div>
        </div>

        <!-- Video Info -->
        <div class="mb-6">
          <h1 class="text-2xl font-bold mb-4">{{ video.title }}</h1>
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center space-x-4">
              <div class="flex items-center space-x-2">
                <div class="w-10 h-10 bg-red-600 rounded-full flex items-center justify-center">
                  {{ video.username.charAt(0).toUpperCase() }}
                </div>
                <div>
                  <p class="font-medium">{{ video.username }}</p>
                </div>
              </div>
            </div>
            <div class="text-sm text-gray-400">
              {{ video.views }} views • {{ formatDate(video.created_at) }}
            </div>
          </div>
        </div>

        <!-- Description -->
        <div v-if="video.description" class="bg-gray-800 rounded-lg p-6">
          <h2 class="text-lg font-semibold mb-3">Description</h2>
          <p class="text-gray-300 whitespace-pre-wrap">{{ video.description }}</p>
        </div>
      </div>

      <!-- Sidebar -->
      <div class="lg:col-span-1">
        <div class="bg-gray-800 rounded-lg p-6">
          <h2 class="text-lg font-semibold mb-4">Video Details</h2>
          <div class="space-y-3">
            <div>
              <span class="text-gray-400 text-sm">Uploaded by:</span>
              <p class="font-medium">{{ video.username }}</p>
            </div>
            <div>
              <span class="text-gray-400 text-sm">Views:</span>
              <p class="font-medium">{{ video.views.toLocaleString() }}</p>
            </div>
            <div>
              <span class="text-gray-400 text-sm">Upload Date:</span>
              <p class="font-medium">{{ new Date(video.created_at).toLocaleDateString() }}</p>
            </div>
          </div>
        </div>

        <!-- Actions (if owner) -->
        <div v-if="isOwner" class="bg-gray-800 rounded-lg p-6 mt-6">
          <h2 class="text-lg font-semibold mb-4">Manage Video</h2>
          <div class="space-y-3">
            <button
              @click="editVideo"
              class="w-full bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md font-medium transition-colors"
            >
              Edit Video
            </button>
            <button
              @click="deleteVideo"
              class="w-full bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded-md font-medium transition-colors"
            >
              Delete Video
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { videoAPI } from '../services/api'
import { authService } from '../services/auth'

export default {
  name: 'VideoDetail',
  data() {
    return {
      video: null,
      loading: true,
      error: null
    }
  },
  computed: {
    isOwner() {
      const user = authService.getUser()
      return user && this.video && user.id === this.video.user_id
    }
  },
  async mounted() {
    await this.fetchVideo()
  },
  methods: {
    async fetchVideo() {
      try {
        this.loading = true
        this.error = null
        const response = await videoAPI.getById(this.$route.params.id)
        this.video = response.data
      } catch (error) {
        console.error('Failed to fetch video:', error)
        if (error.response?.status === 404) {
          this.error = 'Video not found.'
        } else {
          this.error = 'Failed to load video. Please try again later.'
        }
      } finally {
        this.loading = false
      }
    },
    editVideo() {
      // TODO: Implement edit functionality
      alert('Edit functionality coming soon!')
    },
    async deleteVideo() {
      if (!confirm('Are you sure you want to delete this video? This action cannot be undone.')) {
        return
      }
      
      try {
        await videoAPI.delete(this.video.id)
        this.$router.push('/my-videos')
      } catch (error) {
        console.error('Failed to delete video:', error)
        alert('Failed to delete video. Please try again.')
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