import { Header } from './header'

interface RootLayoutProps {
  children: React.ReactNode
}

export function RootLayout({ children }: RootLayoutProps) {
  return (
    <div className="relative min-h-screen">
      <Header />
      <main className="container py-6">{children}</main>
    </div>
  )
} 