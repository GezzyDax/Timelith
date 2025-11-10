'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { api } from '@/lib/api'
import { DatabaseStep } from '@/components/setup/DatabaseStep'
import { AdminStep } from '@/components/setup/AdminStep'
import { ApiKeysStep } from '@/components/setup/ApiKeysStep'

type SetupStep = 'database' | 'admin' | 'apikeys' | 'complete'

export default function SetupPage() {
  const router = useRouter()

  const [currentStep, setCurrentStep] = useState<SetupStep>('database')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleDatabaseNext = async (data: { useDockerDatabase: boolean; databaseUrl?: string }) => {
    setError('')
    setLoading(true)

    try {
      await api.setupDatabase({
        use_docker_database: data.useDockerDatabase,
        database_url: data.databaseUrl,
      })

      setCurrentStep('admin')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to configure database')
    } finally {
      setLoading(false)
    }
  }

  const handleAdminNext = async (data: { username: string; password: string }) => {
    setError('')
    setLoading(true)

    try {
      await api.setupAdmin(data)

      setCurrentStep('apikeys')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create admin user')
    } finally {
      setLoading(false)
    }
  }

  const handleApiKeysComplete = async (data: { telegramAppId: string; telegramAppHash: string }) => {
    setError('')
    setLoading(true)

    try {
      await api.setupComplete({
        telegram_app_id: data.telegramAppId,
        telegram_app_hash: data.telegramAppHash,
      })

      setCurrentStep('complete')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to complete setup')
    } finally {
      setLoading(false)
    }
  }

  const handleBack = () => {
    setError('')
    if (currentStep === 'admin') {
      setCurrentStep('database')
    } else if (currentStep === 'apikeys') {
      setCurrentStep('admin')
    }
  }

  if (currentStep === 'complete') {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background p-4">
        <Card className="w-full max-w-2xl">
          <CardHeader>
            <CardTitle className="text-2xl text-green-600">Setup Complete!</CardTitle>
            <CardDescription>
              Your Timelith installation has been configured successfully.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Alert className="bg-green-50 border-green-200">
              <AlertDescription>
                <div className="space-y-3">
                  <p className="font-semibold text-green-800">
                    ✓ All configuration has been saved to the database
                  </p>
                  <p className="text-green-700">
                    Your settings are now live - no restart required! The server automatically reloaded the configuration.
                  </p>
                  <div className="pt-2">
                    <p className="font-semibold mb-2 text-green-800">You can now:</p>
                    <ul className="list-disc list-inside space-y-1 text-sm text-green-700">
                      <li>Log in with your admin credentials</li>
                      <li>Start adding Telegram accounts</li>
                      <li>Create message templates</li>
                      <li>Configure channels and schedules</li>
                      <li>Manage all settings through the admin panel</li>
                    </ul>
                  </div>
                </div>
              </AlertDescription>
            </Alert>

            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <div className="flex items-start">
                <svg className="w-5 h-5 text-blue-600 mt-0.5 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <div className="text-sm text-blue-800">
                  <p className="font-semibold mb-1">Hot Reload Enabled</p>
                  <p>All changes you make through the admin panel will take effect immediately without requiring a server restart.</p>
                </div>
              </div>
            </div>

            <div className="flex justify-center pt-4">
              <button
                onClick={() => router.push('/login')}
                className="px-6 py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors"
              >
                Go to Login →
              </button>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <div className="w-full max-w-4xl">
        {/* Progress indicator */}
        <div className="mb-8">
          <div className="flex items-center justify-between max-w-2xl mx-auto">
            <div className="flex flex-col items-center flex-1">
              <div className={`w-10 h-10 rounded-full flex items-center justify-center font-semibold ${
                currentStep === 'database' ? 'bg-primary text-white' :
                ['admin', 'apikeys'].includes(currentStep) ? 'bg-green-500 text-white' : 'bg-gray-200 text-gray-600'
              }`}>
                {['admin', 'apikeys'].includes(currentStep) ? '✓' : '1'}
              </div>
              <span className="text-sm mt-2 font-medium">Database</span>
            </div>
            <div className={`h-1 flex-1 mx-4 ${
              ['admin', 'apikeys'].includes(currentStep) ? 'bg-green-500' : 'bg-gray-200'
            }`} />
            <div className="flex flex-col items-center flex-1">
              <div className={`w-10 h-10 rounded-full flex items-center justify-center font-semibold ${
                currentStep === 'admin' ? 'bg-primary text-white' :
                currentStep === 'apikeys' ? 'bg-green-500 text-white' : 'bg-gray-200 text-gray-600'
              }`}>
                {currentStep === 'apikeys' ? '✓' : '2'}
              </div>
              <span className="text-sm mt-2 font-medium">Admin User</span>
            </div>
            <div className={`h-1 flex-1 mx-4 ${
              currentStep === 'apikeys' ? 'bg-green-500' : 'bg-gray-200'
            }`} />
            <div className="flex flex-col items-center flex-1">
              <div className={`w-10 h-10 rounded-full flex items-center justify-center font-semibold ${
                currentStep === 'apikeys' ? 'bg-primary text-white' : 'bg-gray-200 text-gray-600'
              }`}>
                3
              </div>
              <span className="text-sm mt-2 font-medium">API Keys</span>
            </div>
          </div>
        </div>

        {/* Error display */}
        {error && (
          <div className="max-w-2xl mx-auto mb-4">
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          </div>
        )}

        {/* Step components */}
        {currentStep === 'database' && (
          <DatabaseStep onNext={handleDatabaseNext} />
        )}

        {currentStep === 'admin' && (
          <AdminStep onNext={handleAdminNext} onBack={handleBack} />
        )}

        {currentStep === 'apikeys' && (
          <ApiKeysStep
            onNext={handleApiKeysComplete}
            onBack={handleBack}
            isLoading={loading}
          />
        )}
      </div>
    </div>
  )
}
