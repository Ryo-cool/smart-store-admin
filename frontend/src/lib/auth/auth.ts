import { AuthResponse, User } from './types'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export const auth = {
  login: async (): Promise<string> => {
    const response = await fetch(`${API_URL}/api/auth/google`)
    const data = await response.json()
    return data.url
  },

  getUser: async (token: string): Promise<User> => {
    const response = await fetch(`${API_URL}/api/auth/me`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    if (!response.ok) {
      throw new Error('Failed to get user')
    }
    return response.json()
  },

  handleCallback: async (code: string): Promise<AuthResponse> => {
    const response = await fetch(`${API_URL}/api/auth/google/callback?code=${code}`)
    if (!response.ok) {
      throw new Error('Failed to authenticate')
    }
    return response.json()
  },

  logout: async (): Promise<void> => {
    await fetch(`${API_URL}/api/auth/logout`, {
      method: 'POST',
    })
  },
} 