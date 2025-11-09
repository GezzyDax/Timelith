'use client'

import { useRouter, usePathname } from 'next/navigation'
import { useLocale } from 'next-intl'

export function LanguageToggle() {
  const router = useRouter()
  const pathname = usePathname()
  const currentLocale = useLocale()

  const toggleLanguage = () => {
    const newLocale = currentLocale === 'en' ? 'ru' : 'en'

    // Remove current locale from pathname if present
    const pathnameWithoutLocale = pathname.replace(/^\/(en|ru)/, '')

    // Navigate to the same path with new locale
    router.push(`/${newLocale}${pathnameWithoutLocale || '/'}`)
  }

  return (
    <button
      onClick={toggleLanguage}
      className="px-3 py-2 rounded-md border border-border hover:bg-accent transition-colors font-medium text-sm"
      aria-label="Toggle language"
      title={currentLocale === 'en' ? 'Switch to Russian' : 'Переключить на английский'}
    >
      {currentLocale === 'en' ? (
        <span className="flex items-center gap-1">
          <span>🇷🇺</span>
          <span>RU</span>
        </span>
      ) : (
        <span className="flex items-center gap-1">
          <span>🇬🇧</span>
          <span>EN</span>
        </span>
      )}
    </button>
  )
}
