import NextAuth, { NextAuthOptions, Account, DefaultUser } from 'next-auth'
import GoogleProvider from 'next-auth/providers/google'
import { auth } from '@/lib/auth/auth'
import { JWT } from 'next-auth/jwt'
import { Session } from 'next-auth'
import { User } from '@/lib/auth/types'

const authOptions: NextAuthOptions = {
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID!,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
    }),
  ],
  callbacks: {
    async signIn({ account }: { account: Account | null }) {
      if (account?.provider === 'google' && account.access_token) {
        try {
          await auth.handleCallback(account.access_token)
          return true
        } catch (error) {
          if (process.env.NODE_ENV === 'development') {
            // eslint-disable-next-line no-console
            console.error('Authentication error:', error)
          }
          return false
        }
      }
      return false
    },
    async jwt({ token, user, account }) {
      if (account?.provider === 'google' && account.access_token) {
        try {
          const response = await auth.handleCallback(account.access_token)
          return {
            ...token,
            accessToken: response.token,
            user: response.user,
          }
        } catch (error) {
          return token
        }
      }
      return token
    },
    async session({ session, token }) {
      return {
        ...session,
        accessToken: (token as JWT & { accessToken: string }).accessToken,
        user: (token as JWT & { user: User }).user,
      }
    },
  },
  pages: {
    signIn: '/auth/signin',
    error: '/auth/error',
  },
}

const handler = NextAuth(authOptions)

export { handler as GET, handler as POST } 