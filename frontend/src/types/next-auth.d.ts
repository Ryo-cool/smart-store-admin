import { User as CustomUser } from '@/lib/auth/types'
import NextAuth, { DefaultSession, DefaultUser } from 'next-auth'
import { JWT as NextAuthJWT } from 'next-auth/jwt'

declare module 'next-auth' {
  interface Session extends DefaultSession {
    accessToken: string
    user: {
      id: string
      email: string | null
      name: string | null
      picture: string | null
      role: CustomUser['role']
    } & DefaultSession['user']
  }

  interface User extends DefaultUser {
    id: string
    email: string | null
    name: string | null
    picture: string | null
    role: CustomUser['role']
  }
}

declare module 'next-auth/jwt' {
  interface JWT extends NextAuthJWT {
    accessToken: string
    user: {
      id: string
      email: string | null
      name: string | null
      picture: string | null
      role: CustomUser['role']
    }
  }
} 