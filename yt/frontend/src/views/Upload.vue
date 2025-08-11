<template>
  <div class="max-w-2xl mx-auto">
    <h1 class="text-3xl font-bold mb-8">Upload Video</h1>

    <div class="bg-gray-800 rounded-lg shadow-lg p-8">
      <div v-if="error" class="bg-red-900 border border-red-700 text-red-100 px-4 py-3 rounded mb-6">
        {{ error }}
      </div>

      <div v-if="success" class="bg-green-900 border border-green-700 text-green-100 px-4 py-3 rounded mb-6">
        Video uploaded successfully! <router-link to="/my-videos" class="underline">View your videos</router-link>
      </div>

      <form @submit.prevent="uploadVideo" class="space-y-6">
        <div>
          <label for="title" class="block text-sm font-medium text-gray-300 mb-2">Title *</label>
          <input
            id="title"
            v-model="form.title"
            type="text"
            required
            maxlength="255"
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent"
            placeholder="Enter video title"
          >
        </div>

        <div>
          <label for="description" class="block text-sm font-medium text-gray-300 mb-2">Description</label>
          <textarea
            id="description"
            v-model="form.description"
            rows="4"
            maxlength="2000"
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent"
            placeholder="Describe your video..."
          ></textarea>
          <p class="text-xs text-gray-500 mt-1">{{ form.description.length }}/2000 characters</p>
        </div>

        <!-- Upload Mode Toggle -->
        <div class="mb-6">
          <div class="flex items-center space-x-4">
            <label class="flex items-center">
              <input
                v-model="uploadMode"
                type="radio"
                value="url"
                class="sr-only"
              >
              <span class="flex items-center cursor-pointer">
                <span class="w-4 h-4 border-2 border-gray-400 rounded-full mr-2 flex items-center justify-center">
                  <span v-if="uploadMode === 'url'" class="w-2 h-2 bg-red-500 rounded-full"></span>
                </span>
                Video URL
              </span>
            </label>
            <label class="flex items-center">
              <input
                v-model="uploadMode"
                type="radio"
                value="file"
                class="sr-only"
              >
              <span class="flex items-center cursor-pointer">
                <span class="w-4 h-4 border-2 border-gray-400 rounded-full mr-2 flex items-center justify-center">
                  <span v-if="uploadMode === 'file'" class="w-2 h-2 bg-red-500 rounded-full"></span>
                </span>
                Upload File
              </span>
            </label>
          </div>
        </div>

        <!-- URL Input (when URL mode is selected) -->
        <div v-if="uploadMode === 'url'">
          <label for="url" class="block text-sm font-medium text-gray-300 mb-2">Video URL *</label>
          <input
            id="url"
            v-model="form.url"
            type="url"
            :required="uploadMode === 'url'"
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent"
            placeholder="https://example.com/video.mp4"
          >
          <p class="text-xs text-gray-500 mt-1">Direct link to your video file (MP4, WebM, etc.)</p>
        </div>

        <!-- File Upload (when file mode is selected) -->
        <div v-if="uploadMode === 'file'">
          <label for="videoFile" class="block text-sm font-medium text-gray-300 mb-2">Video File *</label>
          <div class="relative">
            <input
              id="videoFile"
              ref="fileInput"
              type="file"
              accept="video/*"
              :required="uploadMode === 'file'"
              @change="handleFileSelect"
              class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-red-600 file:text-white hover:file:bg-red-700"
            >
          </div>
          <p class="text-xs text-gray-500 mt-1">Max file size: 500MB. Supported formats: MP4, WebM, AVI, MOV, etc.</p>
          
          <!-- Upload Progress -->
          <div v-if="uploadProgress > 0 && uploadProgress < 100" class="mt-2">
            <div class="bg-gray-600 rounded-full h-2">
              <div 
                class="bg-red-500 h-2 rounded-full transition-all duration-300"
                :style="{ width: uploadProgress + '%' }"
              ></div>
            </div>
            <p class="text-xs text-gray-400 mt-1">Uploading: {{ uploadProgress }}%</p>
          </div>
        </div>

        <div>
          <label for="thumbnailUrl" class="block text-sm font-medium text-gray-300 mb-2">Thumbnail URL</label>
          <input
            id="thumbnailUrl"
            v-model="form.thumbnail_url"
            type="url"
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent"
            placeholder="https://example.com/thumbnail.jpg"
          >
          <p class="text-xs text-gray-500 mt-1">Optional: Direct link to thumbnail image</p>
        </div>

        <!-- Preview -->
        <div v-if="form.url || form.thumbnail_url" class="border border-gray-600 rounded-lg p-4">
          <h3 class="text-lg font-semibold mb-4">Preview</h3>
          
          <div class="bg-gray-700 rounded-lg overflow-hidden">
            <div class="aspect-video bg-gray-900 flex items-center justify-center">
              <img v-if="form.thumbnail_url" :src="form.thumbnail_url" :alt="form.title" class="w-full h-full object-cover" @error="thumbnailError = true">
              <div v-else-if="form.url" class="text-gray-500">
                <svg class="w-12 h-12" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z"/>
                </svg>
              </div>
            </div>
            <div v-if="form.title || form.description" class="p-4">
              <h4 v-if="form.title" class="font-semibold mb-2">{{ form.title }}</h4>
              <p v-if="form.description" class="text-sm text-gray-400 line-clamp-2">{{ form.description }}</p>
            </div>
          </div>
        </div>

        <div class="flex space-x-4">
          <button
            type="submit"
            :disabled="loading || !isFormValid"
            class="flex-1 bg-red-600 hover:bg-red-700 disabled:bg-red-800 text-white py-3 px-6 rounded-md font-medium transition-colors flex items-center justify-center"
          >
            <div v-if="loading" class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
            {{ loading ? 'Uploading...' : 'Upload Video' }}
          </button>
          
          <router-link
            to="/my-videos"
            class="bg-gray-600 hover:bg-gray-700 text-white py-3 px-6 rounded-md font-medium transition-colors text-center"
          >
            Cancel
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { videoAPI } from '../services/api'
import { authService } from '../services/auth'

export default {
  name: 'Upload',
  data() {
    return {
      uploadMode: 'url', // 'url' or 'file'
      form: {
        title: '',
        description: '',
        url: '',
        thumbnail_url: ''
      },
      selectedFile: null,
      uploadProgress: 0,
      loading: false,
      error: null,
      success: false,
      thumbnailError: false
    }
  },
  computed: {
    isFormValid() {
      const hasTitle = this.form.title.trim()
      if (this.uploadMode === 'url') {
        return hasTitle && this.form.url.trim()
      } else {
        return hasTitle && this.selectedFile
      }
    }
  },
  methods: {
    handleFileSelect(event) {
      const file = event.target.files[0]
      if (file) {
        // Check file size (500MB limit)
        if (file.size > 500 * 1024 * 1024) {
          this.error = 'File size must be less than 500MB'
          this.$refs.fileInput.value = ''
          return
        }
        
        // Check file type
        if (!file.type.startsWith('video/')) {
          this.error = 'Please select a valid video file'
          this.$refs.fileInput.value = ''
          return
        }
        
        this.selectedFile = file
        this.error = null
      }
    },

    async uploadFile() {
      const formData = new FormData()
      formData.append('video', this.selectedFile)
      
      try {
        const response = await fetch('/api/upload', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${authService.getToken()}`
          },
          body: formData
        })
        
        if (!response.ok) {
          throw new Error('Upload failed')
        }
        
        const result = await response.json()
        return result.url
      } catch (error) {
        throw new Error('Failed to upload file: ' + error.message)
      }
    },

    async uploadVideo() {
      if (!this.isFormValid) {
        this.error = 'Please fill in required fields.'
        return
      }

      try {
        this.loading = true
        this.error = null
        this.success = false
        this.uploadProgress = 0
        
        let videoUrl = this.form.url
        
        // If file mode, upload the file first
        if (this.uploadMode === 'file' && this.selectedFile) {
          videoUrl = await this.uploadFile()
          this.uploadProgress = 100
        }
        
        // Create video record
        const videoData = {
          title: this.form.title,
          description: this.form.description,
          url: videoUrl,
          thumbnail_url: this.form.thumbnail_url
        }
        
        await videoAPI.create(videoData)
        
        this.success = true
        this.resetForm()
        
        // Redirect after success
        setTimeout(() => {
          this.$router.push('/my-videos')
        }, 2000)
        
      } catch (error) {
        console.error('Failed to upload video:', error)
        this.error = error.response?.data?.message || error.message || 'Failed to upload video. Please try again.'
        this.uploadProgress = 0
      } finally {
        this.loading = false
      }
    },

    resetForm() {
      this.form = {
        title: '',
        description: '',
        url: '',
        thumbnail_url: ''
      }
      this.selectedFile = null
      this.uploadProgress = 0
      if (this.$refs.fileInput) {
        this.$refs.fileInput.value = ''
      }
    }
  }
}
</script>