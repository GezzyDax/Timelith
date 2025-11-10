'use client'

import { useState } from 'react'
import { useTranslations } from 'next-intl'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Alert, AlertDescription } from '@/components/ui/alert'

interface AdminStepProps {
  onNext: (data: { username: string; password: string }) => void
  onBack: () => void
}

export function AdminStep({ onNext, onBack }: AdminStepProps) {
  const t = useTranslations('setup.admin')
  const tCommon = useTranslations('common')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [showPassword, setShowPassword] = useState(false)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    // Validation
    if (!username) {
      setError(t('errorUsernameRequired'))
      return
    }

    if (username.length < 3) {
      setError(t('errorUsernameShort'))
      return
    }

    if (!password) {
      setError(t('errorPasswordRequired'))
      return
    }

    if (password.length < 6) {
      setError(t('errorPasswordShort'))
      return
    }

    if (password !== confirmPassword) {
      setError(t('errorPasswordMismatch'))
      return
    }

    onNext({ username, password })
  }

  const getPasswordStrength = (pwd: string) => {
    if (pwd.length === 0) return { strength: 0, label: '', color: '' }
    if (pwd.length < 6) return { strength: 25, label: t('strengthTooShort'), color: 'bg-destructive' }

    let strength = 0
    if (pwd.length >= 8) strength += 25
    if (pwd.length >= 12) strength += 25
    if (/[a-z]/.test(pwd) && /[A-Z]/.test(pwd)) strength += 25
    if (/[0-9]/.test(pwd)) strength += 12.5
    if (/[^a-zA-Z0-9]/.test(pwd)) strength += 12.5

    if (strength < 40) return { strength, label: t('strengthWeak'), color: 'bg-destructive' }
    if (strength < 70) return { strength, label: t('strengthMedium'), color: 'bg-yellow-500' }
    return { strength, label: t('strengthStrong'), color: 'bg-success' }
  }

  const passwordStrength = getPasswordStrength(password)

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
              <div className="space-y-2">
                <Label htmlFor="username">{t('username')}</Label>
                <Input
                  id="username"
                  type="text"
                  placeholder={t('usernamePlaceholder')}
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  autoComplete="username"
                  required
                />
                <p className="text-sm text-muted-foreground">
                  {t('usernameHelp')}
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="password">{t('password')}</Label>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? 'text' : 'password'}
                    placeholder={t('passwordPlaceholder')}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    autoComplete="new-password"
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                  >
                    {showPassword ? (
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
                {password && (
                  <div className="space-y-1">
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">{t('passwordStrength')}</span>
                      <span className="font-medium">{passwordStrength.label}</span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-2">
                      <div
                        className={`h-2 rounded-full transition-all ${passwordStrength.color}`}
                        style={{ width: `${passwordStrength.strength}%` }}
                      />
                    </div>
                  </div>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="confirm-password">{t('confirmPassword')}</Label>
                <Input
                  id="confirm-password"
                  type={showPassword ? 'text' : 'password'}
                  placeholder={t('passwordPlaceholder')}
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  autoComplete="new-password"
                  required
                />
              </div>
            </div>

            <div className="flex justify-between pt-4">
              <Button type="button" variant="outline" onClick={onBack}>
                {tCommon('back')}
              </Button>
              <Button type="submit">
                {t('nextButton')}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
