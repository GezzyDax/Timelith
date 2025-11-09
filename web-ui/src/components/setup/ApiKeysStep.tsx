'use client'

import { useState } from 'react'
import { useTranslations } from 'next-intl'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Alert, AlertDescription } from '@/components/ui/alert'

interface ApiKeysStepProps {
  onNext: (data: { telegramAppId: string; telegramAppHash: string }) => void
  onBack: () => void
  isLoading?: boolean
}

export function ApiKeysStep({ onNext, onBack, isLoading }: ApiKeysStepProps) {
  const t = useTranslations('setup.apiKeys')
  const tCommon = useTranslations('common')
  const [telegramAppId, setTelegramAppId] = useState('')
  const [telegramAppHash, setTelegramAppHash] = useState('')
  const [showHash, setShowHash] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    // Validation
    if (!telegramAppId) {
      setError(t('errorAppIdRequired'))
      return
    }

    if (!/^\d+$/.test(telegramAppId)) {
      setError(t('errorAppIdInvalid'))
      return
    }

    if (!telegramAppHash) {
      setError(t('errorAppHashRequired'))
      return
    }

    if (telegramAppHash.length !== 32) {
      setError(t('errorAppHashLength'))
      return
    }

    onNext({
      telegramAppId,
      telegramAppHash,
    })
  }

  return (
    <div className="w-full max-w-2xl mx-auto">
      <Card>
        <CardHeader>
          <CardTitle>{t('title')}</CardTitle>
          <CardDescription>{t('description')}</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <Alert>
              <AlertDescription>
                <div className="space-y-2">
                  <p className="font-semibold">{t('howToTitle')}</p>
                  <ol className="list-decimal list-inside space-y-1 text-sm">
                    <li>{t('howToStep1')}</li>
                    <li>{t('howToStep2')}</li>
                    <li>{t('howToStep3')}</li>
                    <li>{t('howToStep4')}</li>
                    <li>{t('howToStep5')}</li>
                  </ol>
                </div>
              </AlertDescription>
            </Alert>

            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="telegram-app-id">{t('appId')}</Label>
                <Input
                  id="telegram-app-id"
                  type="text"
                  placeholder={t('appIdPlaceholder')}
                  value={telegramAppId}
                  onChange={(e) => setTelegramAppId(e.target.value.replace(/\D/g, ''))}
                  required
                  disabled={isLoading}
                />
                <p className="text-sm text-muted-foreground">
                  {t('appIdHelp')}
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="telegram-app-hash">{t('appHash')}</Label>
                <div className="relative">
                  <Input
                    id="telegram-app-hash"
                    type={showHash ? 'text' : 'password'}
                    placeholder={t('appHashPlaceholder')}
                    value={telegramAppHash}
                    onChange={(e) => setTelegramAppHash(e.target.value.toLowerCase().replace(/[^a-f0-9]/g, ''))}
                    maxLength={32}
                    required
                    disabled={isLoading}
                  />
                  <button
                    type="button"
                    onClick={() => setShowHash(!showHash)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                    disabled={isLoading}
                  >
                    {showHash ? (
                      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                      </svg>
                    ) : (
                      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      </svg>
                    )}
                  </button>
                </div>
                <p className="text-sm text-muted-foreground">
                  {t('appHashHelp')}
                </p>
                {telegramAppHash && telegramAppHash.length !== 32 && (
                  <p className="text-sm text-destructive">
                    {t('currentLength', { length: telegramAppHash.length })}
                  </p>
                )}
              </div>
            </div>

            <div className="bg-primary/10 border border-primary/20 rounded-lg p-4">
              <div className="flex items-start">
                <svg className="w-5 h-5 text-primary mt-0.5 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <div className="text-sm">
                  <p className="font-semibold mb-1">{t('securityNote')}</p>
                  <p className="text-muted-foreground">{t('securityDesc')}</p>
                </div>
              </div>
            </div>

            <div className="flex justify-between pt-4">
              <Button type="button" variant="outline" onClick={onBack} disabled={isLoading}>
                {tCommon('back')}
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? (
                  <>
                    <svg className="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    {t('completing')}
                  </>
                ) : (
                  t('completeButton')
                )}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
