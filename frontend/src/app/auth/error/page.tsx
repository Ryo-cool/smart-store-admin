import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import Link from 'next/link'

export default function ErrorPage() {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <Card className="w-[400px]">
        <CardHeader>
          <CardTitle>認証エラー</CardTitle>
          <CardDescription>
            ログイン処理中にエラーが発生しました。
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Button asChild className="w-full">
            <Link href="/auth/signin">
              ログイン画面に戻る
            </Link>
          </Button>
        </CardContent>
      </Card>
    </div>
  )
} 