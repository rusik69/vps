const TOKEN_KEY = 'youtube_clone_token'
const USER_KEY = 'youtube_clone_user'

export const authService = {
  login(token, user) {
    localStorage.setItem(TOKEN_KEY, token)
    localStorage.setItem(USER_KEY, JSON.stringify(user))
  },

  logout() {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  },

  getToken() {
    return localStorage.getItem(TOKEN_KEY)
  },

  getUser() {
    const user = localStorage.getItem(USER_KEY)
    if (!user) return null
    
    try {
      return JSON.parse(user)
    } catch (error) {
      // If JSON parsing fails, return null and clean up invalid data
      this.logout()
      return null
    }
  },

  isAuthenticated() {
    return !!this.getToken()
  }
}