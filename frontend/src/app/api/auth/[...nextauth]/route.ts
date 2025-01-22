import NextAuth, { NextAuthOptions, Account } from 'next-auth'
import GoogleProvider from 'next-auth/providers/google'
import { auth } from '@/lib/auth/auth'
import { JWT } from 'next-auth/jwt'
import { Session } from 'next-auth'

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
            console.error('Authentication error:', error)
          }
          return false
        }
      }
      return false
    },
    async jwt({ token, user, account }: { token: JWT, user: any, account: Account | null }): Promise<JWT> {
      if (account && user && account.access_token) {
        const response = await auth.handleCallback(account.access_token)
        return {
          ...token,
          accessToken: response.token,
          user: response.user,
        }
      }
      return token
    },
    async session({ session, token }: { session: Session, token: JWT }): Promise<Session> {
      return {
        ...session,
        accessToken: token.accessToken,
        user: token.user,
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