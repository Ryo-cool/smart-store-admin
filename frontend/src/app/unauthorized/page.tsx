'use client'

import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'

export default function UnauthorizedPage() {
  const router = useRouter()

  return (
    <div className="container flex h-screen w-screen flex-col items-center justify-center">
      <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
        <div className="flex flex-col space-y-2 text-center">
          <h1 className="text-2xl font-semibold tracking-tight">
            アクセス権限がありません
          </h1>
          <p className="text-sm text-muted-foreground">
            このページにアクセスするための権限が不足しています。
          </p>
        </div>
        <Button onClick={() => router.back()}>
          前のページに戻る
        </Button>
        <Button variant="outline" asChild>
          <a href="/dashboard">ダッシュボードへ</a>
        </Button>
      </div>
    </div>
  )
} 