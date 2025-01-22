import { withAuth } from 'next-auth/middleware'

export default withAuth({
  pages: {
    signIn: '/auth/signin',
  },
})

export const config = {
  matcher: [
    // 認証が必要なパスを指定
    '/dashboard/:path*',
    '/products/:path*',
    '/sales/:path*',
    '/deliveries/:path*',
    '/settings/:path*',
  ],
} 