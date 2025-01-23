import { authOptions } from '@/lib/auth/auth'
import NextAuth from 'next-auth'

const handler = NextAuth(authOptions)

// Next.js Edge API Routes: https://nextjs.org/docs/app/building-your-application/routing/router-handlers#edge-and-nodejs-runtimes
export const dynamic = 'force-dynamic'
export { handler as GET, handler as POST } 