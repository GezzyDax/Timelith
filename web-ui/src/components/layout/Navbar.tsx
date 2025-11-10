'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useTranslations } from 'next-intl'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { ThemeToggle } from '@/components/theme/theme-toggle'
import { LanguageToggle } from '@/components/theme/language-toggle'

export function Navbar() {
  const pathname = usePathname()
  const router = useRouter()
  const t = useTranslations()

  const navigation = [
    { name: t('dashboard.title'), href: '/' },
    { name: t('accounts.title'), href: '/accounts' },
    { name: t('templates.title'), href: '/templates' },
    { name: t('channels.title'), href: '/channels' },
    { name: t('schedules.title'), href: '/schedules' },
  ]

  const handleLogout = () => {
    localStorage.removeItem('token')
    router.push('/login')
  }

  return (
    <nav className="border-b bg-background">
      <div className="flex h-16 items-center px-4 container mx-auto">
        <div className="flex items-center space-x-8">
          <Link href="/" className="font-bold text-xl text-primary">
            Timelith
          </Link>
          <div className="hidden md:flex space-x-4">
            {navigation.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  'text-sm font-medium transition-colors hover:text-primary',
                  pathname === item.href
                    ? 'text-foreground'
                    : 'text-muted-foreground'
                )}
              >
                {item.name}
              </Link>
            ))}
          </div>
        </div>
        <div className="ml-auto flex items-center gap-2">
          <ThemeToggle />
          <LanguageToggle />
          <Button variant="ghost" onClick={handleLogout}>
            {t('common.logout') || 'Logout'}
          </Button>
        </div>
      </div>
    </nav>
  )
}
