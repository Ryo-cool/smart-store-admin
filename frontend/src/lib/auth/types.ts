import type { DefaultSession } from 'next-auth'

declare module 'next-auth' {
  interface Session extends DefaultSession {
    accessToken: string | undefined
    user: {
      id: string
      role: string
      picture?: string
    } & DefaultSession['user']
  }
}

declare module 'next-auth/jwt' {
  interface JWT extends DefaultSession {
    accessToken: string | undefined
    user: {
      id: string
      role: string
      picture?: string
    } & DefaultSession['user']
  }
}

export type Role = 'admin' | 'staff' | 'viewer'

export interface User {
  id: string
  email: string | null
  name: string | null
  picture: string | null
  role: Role
}

export interface Session {
  user: User
  expires: string
}

export interface AuthResponse {
  accessToken: string
  user: {
    id: string
    email: string
    name: string
    role: Role
    picture?: string
  }
} 