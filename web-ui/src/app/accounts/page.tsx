'use client'

import { useEffect, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Navbar } from '@/components/layout/Navbar'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { api } from '@/lib/api'
import { formatDate } from '@/lib/utils'
import type { Account } from '@/types'
import { Plus, Trash2 } from 'lucide-react'

export const dynamic = 'force-dynamic'

const getErrorMessage = (error: unknown) => {
  if (error && typeof error === 'object' && 'response' in error) {
    const response = (error as any).response
    return response?.data?.error || response?.data?.message || 'Server error'
  }
  if (error instanceof Error) {
    return error.message
  }
  return 'Unexpected error'
}

export default function AccountsPage() {
  const queryClient = useQueryClient()
  const [showCreate, setShowCreate] = useState(false)
  const [phone, setPhone] = useState('')
  const [code, setCode] = useState('')
  const [password, setPassword] = useState('')
  const [wizardStep, setWizardStep] = useState<'phone' | 'code' | 'password'>('phone')
  const [currentAccount, setCurrentAccount] = useState<Account | null>(null)
  const [passwordHint, setPasswordHint] = useState<string | undefined>()
  const [formError, setFormError] = useState('')
  const [formMessage, setFormMessage] = useState('')
  const [alertMessage, setAlertMessage] = useState('')

  const { data: accounts, isLoading } = useQuery({
    queryKey: ['accounts'],
    queryFn: () => api.getAccounts(),
  })

  const resetWizard = () => {
    setWizardStep('phone')
    setPhone('')
    setCode('')
    setPassword('')
    setPasswordHint(undefined)
    setCurrentAccount(null)
    setFormError('')
    setFormMessage('')
  }

  const finalizeAccountLink = (message: string) => {
    queryClient.invalidateQueries({ queryKey: ['accounts'] })
    setAlertMessage(message)
    setShowCreate(false)
    resetWizard()
    setTimeout(() => setAlertMessage(''), 4000)
  }

  useEffect(() => {
    if (!showCreate) {
      resetWizard()
    }
  }, [showCreate])

  const sendCodeMutation = useMutation({
    mutationFn: (phone: string) => api.createAccount({ phone }),
    onSuccess: (account) => {
      setCurrentAccount(account)
      setWizardStep('code')
      setFormMessage(`Verification code sent to ${account.phone}`)
      setFormError('')
      setPasswordHint(undefined)
    },
    onError: (error) => {
      setFormError(getErrorMessage(error))
    },
  })

  const verifyCodeMutation = useMutation({
    mutationFn: (codeValue: string) => {
      if (!currentAccount) {
        throw new Error('No account to verify')
      }
      return api.verifyAccountCode(currentAccount.id, codeValue)
    },
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
      setCurrentAccount(result.account)
      setFormError('')

      if (result.requires_password) {
        setWizardStep('password')
        setPassword('')
        setPasswordHint(result.password_hint)
        setFormMessage('Enter the two-factor password to finish linking this account.')
      } else {
        finalizeAccountLink('Account linked successfully.')
      }
    },
    onError: (error) => {
      setFormError(getErrorMessage(error))
    },
  })

  const verifyPasswordMutation = useMutation({
    mutationFn: (pwd: string) => {
      if (!currentAccount) {
        throw new Error('No account to verify')
      }
      return api.verifyAccountPassword(currentAccount.id, pwd)
    },
    onSuccess: (account) => {
      setCurrentAccount(account)
      setFormError('')
      finalizeAccountLink('Two-factor password accepted. Account is active.')
    },
    onError: (error) => {
      setFormError(getErrorMessage(error))
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.deleteAccount(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
    },
  })

  const handleFormSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setFormError('')
    setFormMessage('')

    if (wizardStep === 'phone') {
      sendCodeMutation.mutate(phone)
      return
    }

    if (!currentAccount) {
      setFormError('Start with entering a phone number.')
      setWizardStep('phone')
      return
    }

    if (wizardStep === 'code') {
      verifyCodeMutation.mutate(code)
      return
    }

    verifyPasswordMutation.mutate(password)
  }

  const handleCancel = () => {
    setShowCreate(false)
    resetWizard()
  }

  const isSubmitting =
    sendCodeMutation.isPending || verifyCodeMutation.isPending || verifyPasswordMutation.isPending
  const submitLabel =
    wizardStep === 'phone'
      ? 'Send Code'
      : wizardStep === 'code'
        ? 'Verify Code'
        : 'Verify Password'
  const pendingLabel =
    wizardStep === 'phone'
      ? 'Sending...'
      : wizardStep === 'code'
        ? 'Verifying...'
        : 'Verifying...'

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <main className="container mx-auto py-8">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold">Telegram Accounts</h1>
            <p className="text-muted-foreground mt-2">
              Manage your Telegram accounts for sending messages
            </p>
          </div>
          <Button onClick={() => setShowCreate(!showCreate)}>
            <Plus className="mr-2 h-4 w-4" />
            Add Account
          </Button>
        </div>

        {alertMessage && (
          <div className="mb-4 rounded-md bg-emerald-50 px-4 py-2 text-sm text-emerald-700">
            {alertMessage}
          </div>
        )}

        {showCreate && (
          <Card className="mb-6">
            <CardHeader>
              <CardTitle>Link Telegram Account</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleFormSubmit} className="space-y-4">
                {wizardStep === 'phone' && (
                  <div className="space-y-2">
                    <Label htmlFor="phone">Phone Number</Label>
                    <Input
                      id="phone"
                      type="text"
                      placeholder="+1234567890"
                      value={phone}
                      onChange={(e) => setPhone(e.target.value)}
                      required
                    />
                  </div>
                )}

                {wizardStep === 'code' && (
                  <div className="space-y-2">
                    <Label htmlFor="code">
                      Verification Code{' '}
                      {currentAccount?.phone && (
                        <span className="text-xs text-muted-foreground">({currentAccount.phone})</span>
                      )}
                    </Label>
                    <Input
                      id="code"
                      type="text"
                      placeholder="12345"
                      value={code}
                      onChange={(e) => setCode(e.target.value)}
                      required
                    />
                    <p className="text-xs text-muted-foreground">
                      Enter the code you received via Telegram or SMS.
                    </p>
                  </div>
                )}

                {wizardStep === 'password' && (
                  <div className="space-y-2">
                    <Label htmlFor="password">
                      Two-Factor Password{' '}
                      {passwordHint && (
                        <span className="text-xs text-muted-foreground">(hint: {passwordHint})</span>
                      )}
                    </Label>
                    <Input
                      id="password"
                      type="password"
                      placeholder="Your Telegram cloud password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      required
                    />
                  </div>
                )}

                {formError && <p className="text-sm text-destructive">{formError}</p>}
                {formMessage && <p className="text-sm text-muted-foreground">{formMessage}</p>}

                <div className="flex space-x-2">
                  <Button type="submit" disabled={isSubmitting}>
                    {isSubmitting ? pendingLabel : submitLabel}
                  </Button>
                  <Button type="button" variant="outline" onClick={handleCancel}>
                    Cancel
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        )}

        {isLoading ? (
          <p>Loading...</p>
        ) : (
          <div className="grid gap-4">
            {accounts?.map((account) => {
              const statusClass =
                account.status === 'active'
                  ? 'bg-green-100 text-green-800'
                  : account.status === 'error'
                    ? 'bg-red-100 text-red-800'
                    : account.status === 'password_required'
                      ? 'bg-amber-100 text-amber-800'
                      : 'bg-gray-100 text-gray-800'

              return (
                <Card key={account.id}>
                  <CardContent className="flex items-center justify-between p-6">
                    <div>
                      <h3 className="font-semibold">{account.phone}</h3>
                      <div className="flex flex-wrap items-center gap-3 mt-2 text-sm text-muted-foreground">
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${statusClass}`}>
                          {account.status}
                        </span>
                        <span>Last login: {formatDate(account.last_login_at)}</span>
                        {account.login_code_sent_at && account.status !== 'active' && (
                          <span>Code sent: {formatDate(account.login_code_sent_at)}</span>
                        )}
                        <span>Created: {formatDate(account.created_at)}</span>
                        {account.two_factor_required && (
                          <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                            2FA pending
                          </span>
                        )}
                      </div>
                      {account.error_message && (
                        <p className="text-sm text-red-600 mt-1">{account.error_message}</p>
                      )}
                    </div>
                    <Button
                      variant="destructive"
                      size="icon"
                      onClick={() => deleteMutation.mutate(account.id)}
                      disabled={deleteMutation.isPending}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </CardContent>
                </Card>
              )
            })}
            {(!accounts || accounts.length === 0) && (
              <p className="text-center text-muted-foreground py-8">
                No accounts found. Add your first Telegram account to get started.
              </p>
            )}
          </div>
        )}
      </main>
    </div>
  )
}
