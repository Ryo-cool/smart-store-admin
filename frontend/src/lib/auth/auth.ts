import { AuthResponse, User } from './types'
import { NextAuthOptions } from 'next-auth'
import GoogleProvider from 'next-auth/providers/google'
import { Role } from './types'
import type { JWT } from 'next-auth/jwt'
import type { Session } from 'next-auth'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
const allowedDomains = process.env.ALLOWED_EMAIL_DOMAINS?.split(',') || []

export const authOptions: NextAuthOptions = {
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID || '',
      clientSecret: process.env.GOOGLE_CLIENT_SECRET || '',
    }),
  ],
  callbacks: {
    async signIn({ user }) {
      if (!user.email) return false
      if (allowedDomains.length === 0) return true
      return allowedDomains.some((domain) => user.email?.endsWith(domain))
    },
    async jwt({ token, user }) {
      if (user) {
        token.user = {
          ...user,
          role: 'admin' as Role,
        }
      }
      return token
    },
    async session({ session, token }: { session: Session; token: JWT }) {
      session.accessToken = token.accessToken as string
      session.user = token.user as typeof session.user
      return session
    },
  },
  pages: {
    signIn: '/auth/signin',
    error: '/auth/error',
  },
}

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