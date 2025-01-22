import { withAuth } from 'next-auth/middleware'
import { NextResponse } from 'next/server'
import { ROLE_BASED_ROUTES, hasRequiredRole } from '@/lib/auth/utils'
import { Role } from '@/lib/auth/types'

export default withAuth(
  function middleware(req) {
    const token = req.nextauth.token
    const isAuth = !!token
    const isAuthPage = req.nextUrl.pathname.startsWith('/auth')

    if (isAuthPage) {
      if (isAuth) {
        return NextResponse.redirect(new URL('/dashboard', req.url))
      }
      return null
    }

    if (!isAuth) {
      let from = req.nextUrl.pathname
      if (req.nextUrl.search) {
        from += req.nextUrl.search
      }

      return NextResponse.redirect(
        new URL(`/auth/signin?from=${encodeURIComponent(from)}`, req.url)
      )
    }

    // 権限チェック
    const path = req.nextUrl.pathname
    const userRole = token.user?.role as Role

    // パスに対応する必要な権限を取得
    const requiredRole = Object.entries(ROLE_BASED_ROUTES).find(([route]) =>
      path.startsWith(route)
    )?.[1] as Role | undefined

    // 必要な権限が設定されていて、ユーザーの権限が不足している場合
    if (requiredRole && !hasRequiredRole(userRole, requiredRole)) {
      return NextResponse.redirect(new URL('/unauthorized', req.url))
    }
  },
  {
    callbacks: {
      authorized: ({ token }) => !!token,
    },
  }
)

export const config = {
  matcher: [
    // 認証が必要なパスを指定
    '/dashboard/:path*',
    '/products/:path*',
    '/sales/:path*',
    '/deliveries/:path*',
    '/settings/:path*',
    '/auth/:path*',
  ],
} 