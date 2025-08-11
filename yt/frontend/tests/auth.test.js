import { authService } from '../src/services/auth'

// Mock localStorage
const localStorageMock = (() => {
  let store = {}
  return {
    getItem: jest.fn((key) => store[key] || null),
    setItem: jest.fn((key, value) => {
      store[key] = value.toString()
    }),
    removeItem: jest.fn((key) => {
      delete store[key]
    }),
    clear: jest.fn(() => {
      store = {}
    })
  }
})()

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
})

describe('authService', () => {
  beforeEach(() => {
    // Clear all mocks before each test
    jest.clearAllMocks()
  })

  describe('login', () => {
    it('should store token and user in localStorage', () => {
      const token = 'test-token'
      const user = { id: 1, username: 'testuser' }

      authService.login(token, user)

      expect(localStorage.setItem).toHaveBeenCalledWith('youtube_clone_token', token)
      expect(localStorage.setItem).toHaveBeenCalledWith('youtube_clone_user', JSON.stringify(user))
    })
  })

  describe('logout', () => {
    it('should remove token and user from localStorage', () => {
      authService.logout()

      expect(localStorage.removeItem).toHaveBeenCalledWith('youtube_clone_token')
      expect(localStorage.removeItem).toHaveBeenCalledWith('youtube_clone_user')
    })
  })

  describe('getToken', () => {
    it('should return token from localStorage', () => {
      const token = 'test-token'
      localStorage.setItem('youtube_clone_token', token)

      const result = authService.getToken()

      expect(localStorage.getItem).toHaveBeenCalledWith('youtube_clone_token')
      expect(result).toBe(token)
    })

    it('should return null if no token exists', () => {
      localStorage.clear()

      const result = authService.getToken()

      expect(result).toBeNull()
    })
  })

  describe('getUser', () => {
    it('should return parsed user from localStorage', () => {
      const user = { id: 1, username: 'testuser' }
      localStorage.setItem('youtube_clone_user', JSON.stringify(user))

      const result = authService.getUser()

      expect(localStorage.getItem).toHaveBeenCalledWith('youtube_clone_user')
      expect(result).toEqual(user)
    })

    it('should return null if no user exists', () => {
      localStorage.clear()

      const result = authService.getUser()

      expect(result).toBeNull()
    })

    it('should return null if user data is invalid JSON', () => {
      // Manually set invalid JSON to test error handling
      const originalGetItem = localStorage.getItem
      localStorage.getItem = jest.fn().mockReturnValue('invalid-json')

      const result = authService.getUser()

      expect(result).toBeNull()
      
      // Restore original function
      localStorage.getItem = originalGetItem
    })
  })

  describe('isAuthenticated', () => {
    it('should return true if token exists', () => {
      localStorage.setItem('youtube_clone_token', 'test-token')

      const result = authService.isAuthenticated()

      expect(result).toBe(true)
    })

    it('should return false if no token exists', () => {
      localStorage.clear()

      const result = authService.isAuthenticated()

      expect(result).toBe(false)
    })

    it('should return false if token is empty string', () => {
      localStorage.setItem('youtube_clone_token', '')

      const result = authService.isAuthenticated()

      expect(result).toBe(false)
    })
  })
})