import createMiddleware from 'next-intl/middleware'

export default createMiddleware({
  // A list of all locales that are supported
  locales: ['en', 'ru'],

  // Used when no locale matches
  defaultLocale: 'en',

  // Disable automatic locale detection based on browser settings
  localeDetection: true,
})

export const config = {
  // Match only internationalized pathnames
  matcher: ['/', '/(ru|en)/:path*'],
}
