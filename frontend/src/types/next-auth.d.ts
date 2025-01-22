import { User } from '@/lib/auth/types'
import NextAuth, { DefaultSession, DefaultUser } from 'next-auth'
import { JWT as NextAuthJWT } from 'next-auth/jwt'

declare module 'next-auth' {
  interface Session extends DefaultSession {
    accessToken: string
    user: User & DefaultUser
  }

  interface User extends DefaultUser, Omit<User, keyof DefaultUser> {}
}

declare module 'next-auth/jwt' {
  interface JWT extends NextAuthJWT {
    accessToken: string
    user: User
  }
} 