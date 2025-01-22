import { NextAuthOptions } from 'next-auth'
import NextAuth from 'next-auth/next'
import GoogleProvider from 'next-auth/providers/google'
import { User, Role } from '@/lib/auth/types'

const allowedDomains = process.env.ALLOWED_EMAIL_DOMAINS?.split(',') || []

export const authOptions: NextAuthOptions = {
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID!,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
    }),
  ],
  callbacks: {
    async signIn({ user, account }) {
      if (account?.provider !== 'google') return false

      const email = user.email
      if (!email) return false

      const domain = email.split('@')[1]
      if (!domain) return false

      return allowedDomains.includes(domain)
    },
    async jwt({ token, user }) {
      if (user) {
        token.user = {
          id: user.id,
          name: user.name,
          email: user.email,
          picture: user.image ?? null,
          role: 'admin' as Role,
        }
      }
      return token
    },
    async session({ session, token }) {
      session.accessToken = token.accessToken as string
      session.user = token.user
      return session
    },
  },
  pages: {
    signIn: '/auth/signin',
  },
}

const handler = NextAuth(authOptions)

export { handler as GET, handler as POST } 