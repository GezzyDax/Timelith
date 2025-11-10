'use client'

import { useState } from 'react'
import { useTranslations } from 'next-intl'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Alert, AlertDescription } from '@/components/ui/alert'

interface DatabaseStepProps {
  onNext: (data: { useDockerDatabase: boolean; databaseUrl?: string }) => void
  onBack?: () => void
}

export function DatabaseStep({ onNext, onBack }: DatabaseStepProps) {
  const t = useTranslations('setup.database')
  const tCommon = useTranslations('common')
  const [useDocker, setUseDocker] = useState(true)
  const [databaseUrl, setDatabaseUrl] = useState('')
  const [error, setError] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!useDocker && !databaseUrl) {
      setError(t('errorNoUrl'))
      return
    }

    onNext({
      useDockerDatabase: useDocker,
      databaseUrl: useDocker ? undefined : databaseUrl,
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

            <div className="space-y-4">
              <div
                className="flex items-start space-x-3 p-4 border border-border rounded-lg cursor-pointer hover:bg-accent transition-colors"
                onClick={() => setUseDocker(true)}
              >
                <input
                  type="radio"
                  id="docker-db"
                  name="database-type"
                  checked={useDocker}
                  onChange={() => setUseDocker(true)}
                  className="mt-1 accent-primary"
                />
                <div className="flex-1">
                  <Label htmlFor="docker-db" className="cursor-pointer font-semibold">
                    {t('useDocker')}
                  </Label>
                  <p className="text-sm text-muted-foreground mt-1">
                    {t('useDockerDesc')}
                  </p>
                </div>
              </div>

              <div
                className="flex items-start space-x-3 p-4 border border-border rounded-lg cursor-pointer hover:bg-accent transition-colors"
                onClick={() => setUseDocker(false)}
              >
                <input
                  type="radio"
                  id="external-db"
                  name="database-type"
                  checked={!useDocker}
                  onChange={() => setUseDocker(false)}
                  className="mt-1 accent-primary"
                />
                <div className="flex-1">
                  <Label htmlFor="external-db" className="cursor-pointer font-semibold">
                    {t('useExternal')}
                  </Label>
                  <p className="text-sm text-muted-foreground mt-1">
                    {t('useExternalDesc')}
                  </p>
                </div>
              </div>
            </div>

            {!useDocker && (
              <div className="space-y-2">
                <Label htmlFor="database-url">{t('databaseUrl')}</Label>
                <Input
                  id="database-url"
                  type="text"
                  placeholder={t('databaseUrlPlaceholder')}
                  value={databaseUrl}
                  onChange={(e) => setDatabaseUrl(e.target.value)}
                  required={!useDocker}
                />
                <p className="text-sm text-muted-foreground">
                  {t('databaseUrlHelp')}
                </p>
              </div>
            )}

            <div className="flex justify-between pt-4">
              {onBack && (
                <Button type="button" variant="outline" onClick={onBack}>
                  {tCommon('back')}
                </Button>
              )}
              <Button type="submit" className="ml-auto">
                {t('nextButton')}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
