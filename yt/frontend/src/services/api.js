import axios from 'axios'
import { authService } from './auth'

const api = axios.create({
  baseURL: process.env.NODE_ENV === 'production' ? '/api' : 'http://localhost:8080/api',
  timeout: 10000
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = authService.getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      authService.logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const videoAPI = {
  getAll: () => api.get('/videos'),
  getById: (id) => api.get(`/videos/${id}`),
  create: (video) => api.post('/videos', video),
  update: (id, video) => api.put(`/videos/${id}`, video),
  delete: (id) => api.delete(`/videos/${id}`),
  getMyVideos: () => api.get('/my-videos')
}

export const authAPI = {
  login: (credentials) => api.post('/auth/login', credentials),
  register: (userData) => api.post('/auth/register', userData)
}

export default api