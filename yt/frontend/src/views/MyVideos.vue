<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold">My Videos</h1>
      <router-link to="/upload" class="bg-red-600 hover:bg-red-700 text-white px-6 py-2 rounded-md font-medium">
        Upload Video
      </router-link>
    </div>
    
    <div v-if="loading" class="text-center">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-red-500"></div>
    </div>

    <div v-else-if="error" class="bg-red-900 border border-red-700 text-red-100 px-4 py-3 rounded mb-4">
      {{ error }}
    </div>

    <div v-else-if="videos.length === 0" class="text-center py-12">
      <div class="text-gray-500 mb-4">
        <svg class="w-16 h-16 mx-auto" fill="currentColor" viewBox="0 0 20 20">
          <path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z"/>
        </svg>
      </div>
      <h2 class="text-xl font-semibold mb-2">No videos uploaded yet</h2>
      <p class="text-gray-400 mb-4">Start by uploading your first video!</p>
      <router-link to="/upload" class="bg-red-600 hover:bg-red-700 text-white px-6 py-2 rounded-md font-medium">Upload Video</router-link>
    </div>

    <div v-else class="space-y-6">
      <div v-for="video in videos" :key="video.id" class="bg-gray-800 rounded-lg p-6 flex flex-col md:flex-row md:items-center space-y-4 md:space-y-0 md:space-x-6">
        <!-- Thumbnail -->
        <div class="flex-shrink-0">
          <router-link :to="`/video/${video.id}`" class="block">
            <div class="w-full md:w-48 aspect-video bg-gray-700 rounded-lg overflow-hidden flex items-center justify-center">
              <img v-if="video.thumbnail_url" :src="video.thumbnail_url" :alt="video.title" class="w-full h-full object-cover">
              <div v-else class="text-gray-500">
                <svg class="w-12 h-12" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M10 12a2 2 0 100-4 2 2 0 000 4z"/>
                  <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd"/>
                </svg>
              </div>
            </div>
          </router-link>
        </div>

        <!-- Video Info -->
        <div class="flex-grow">
          <router-link :to="`/video/${video.id}`" class="block hover:text-red-400 transition-colors">
            <h2 class="text-xl font-semibold mb-2">{{ video.title }}</h2>
          </router-link>
          <p v-if="video.description" class="text-gray-400 text-sm mb-2 line-clamp-2">{{ video.description }}</p>
          <div class="flex items-center space-x-4 text-xs text-gray-500">
            <span>{{ video.views }} views</span>
            <span>{{ formatDate(video.created_at) }}</span>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex-shrink-0 flex space-x-2">
          <button
            @click="editVideo(video)"
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors"
          >
            Edit
          </button>
          <button
            @click="deleteVideo(video)"
            class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { videoAPI } from '../services/api'

export default {
  name: 'MyVideos',
  data() {
    return {
      videos: [],
      loading: true,
      error: null
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
        const response = await videoAPI.getMyVideos()
        this.videos = response.data || []
      } catch (error) {
        console.error('Failed to fetch videos:', error)
        this.error = 'Failed to load videos. Please try again later.'
      } finally {
        this.loading = false
      }
    },
    editVideo(video) {
      // TODO: Implement edit functionality
      alert(`Edit functionality for "${video.title}" coming soon!`)
    },
    async deleteVideo(video) {
      if (!confirm(`Are you sure you want to delete "${video.title}"? This action cannot be undone.`)) {
        return
      }
      
      try {
        await videoAPI.delete(video.id)
        this.videos = this.videos.filter(v => v.id !== video.id)
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