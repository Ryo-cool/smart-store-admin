'use client'

import { useSession, signOut } from 'next-auth/react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { cn } from '@/lib/utils'

interface HeaderProps extends React.HTMLAttributes<HTMLElement> {
  fixed?: boolean
  children?: React.ReactNode
}

export function Header({ fixed, children, ...props }: HeaderProps) {
  const { data: session, status } = useSession()

  return (
    <header className={cn(
      "sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60",
      fixed && "header-fixed fixed top-0 z-50"
    )} {...props}>
      <div className="container flex h-14 items-center">
        {children}
        <div className="mr-4 flex">
          <Link href="/" className="mr-6 flex items-center space-x-2">
            <span className="font-bold">NEXT MART 2030</span>
          </Link>
        </div>

        <div className="flex flex-1 items-center justify-between space-x-2 md:justify-end">
          <nav className="flex items-center space-x-2">
            {status === 'authenticated' && session?.user ? (
              <>
                <Button variant="ghost" asChild>
                  <Link href="/dashboard">ダッシュボード</Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/products">商品管理</Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/sales">売上管理</Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/deliveries">配送管理</Link>
                </Button>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                      <Avatar className="h-8 w-8">
                        <AvatarImage src={session.user.picture || ''} alt={session.user.name || ''} />
                        <AvatarFallback>{session.user.name?.[0] || 'U'}</AvatarFallback>
                      </Avatar>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem className="font-medium">
                      {session.user.name || 'ユーザー'}
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => signOut()}>
                      ログアウト
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </>
            ) : (
              <Button asChild>
                <Link href="/auth/signin">ログイン</Link>
              </Button>
            )}
          </nav>
        </div>
      </div>
    </header>
  )
}
