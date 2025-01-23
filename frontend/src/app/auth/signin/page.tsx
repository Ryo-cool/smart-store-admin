'use client'

import { signIn } from 'next-auth/react'
import { useSearchParams } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Suspense } from 'react'

export default function SignIn() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <SignInContent />
    </Suspense>
  )
}

function SignInContent() {
  const searchParams = useSearchParams()
  const from = searchParams.get('from') || '/dashboard'

  return (
    <div className="container flex h-screen w-screen flex-col items-center justify-center">
      <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
        <div className="flex flex-col space-y-2 text-center">
          <h1 className="text-2xl font-semibold tracking-tight">
            ログイン
          </h1>
          <p className="text-sm text-muted-foreground">
            Googleアカウントでログインしてください
          </p>
        </div>
        <Button
          onClick={() => signIn('google', { callbackUrl: from })}
          className="w-full"
        >
          Googleでログイン
        </Button>
      </div>
    </div>
  )
} 