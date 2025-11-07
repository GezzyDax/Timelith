'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { api } from '@/lib/api'

export default function SetupPage() {
  const router = useRouter()

  const [step, setStep] = useState(1)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

  // Form state
  const [telegramAppId, setTelegramAppId] = useState('')
  const [telegramAppHash, setTelegramAppHash] = useState('')
  const [serverPort, setServerPort] = useState('8080')
  const [postgresPassword, setPostgresPassword] = useState('timelith_password')
  const [adminUsername, setAdminUsername] = useState('admin')
  const [adminPassword, setAdminPassword] = useState('')
  const [adminPasswordConfirm, setAdminPasswordConfirm] = useState('')
  const [environment, setEnvironment] = useState('production')

  const validateStep1 = () => {
    if (!telegramAppId || !telegramAppHash) {
      setError('Please fill in all Telegram API credentials')
      return false
    }
    if (isNaN(parseInt(telegramAppId))) {
      setError('Telegram App ID must be a number')
      return false
    }
    if (telegramAppHash.length < 32) {
      setError('Telegram App Hash is too short')
      return false
    }
    return true
  }

  const validateStep2 = () => {
    const port = parseInt(serverPort)
    if (isNaN(port) || port < 1 || port > 65535) {
      setError('Invalid port number')
      return false
    }
    if (!postgresPassword) {
      setError('PostgreSQL password is required')
      return false
    }
    return true
  }

  const validateStep3 = () => {
    if (adminUsername.length < 3) {
      setError('Admin username must be at least 3 characters')
      return false
    }
    if (adminPassword.length < 6) {
      setError('Admin password must be at least 6 characters')
      return false
    }
    if (adminPassword !== adminPasswordConfirm) {
      setError('Passwords do not match')
      return false
    }
    return true
  }

  const handleNext = () => {
    setError('')

    if (step === 1 && !validateStep1()) return
    if (step === 2 && !validateStep2()) return

    setStep(step + 1)
  }

  const handleBack = () => {
    setError('')
    setStep(step - 1)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!validateStep3()) return

    setLoading(true)

    try {
      await api.performSetup({
        telegram_app_id: telegramAppId,
        telegram_app_hash: telegramAppHash,
        server_port: serverPort,
        postgres_password: postgresPassword,
        admin_username: adminUsername,
        admin_password: adminPassword,
        environment: environment,
      })

      setSuccess(true)

      // Show success message and redirect
      setTimeout(() => {
        router.push('/login')
      }, 3000)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Setup failed. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  if (success) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <Card className="w-full max-w-2xl">
          <CardHeader>
            <CardTitle className="text-2xl text-green-600">Setup Complete!</CardTitle>
            <CardDescription>
              Your Timelith installation has been configured successfully.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Alert>
              <AlertDescription>
                <p className="font-semibold mb-2">Next steps:</p>
                <ol className="list-decimal list-inside space-y-1">
                  <li>The server will need to be restarted for changes to take effect</li>
                  <li>Stop the current process (Ctrl+C in terminal)</li>
                  <li>Start the server again with: <code className="bg-gray-100 px-1 rounded">go run cmd/server/main.go</code></li>
                  <li>Login with your admin credentials</li>
                </ol>
              </AlertDescription>
            </Alert>
            <p className="text-sm text-muted-foreground">
              Redirecting to login page in 3 seconds...
            </p>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <Card className="w-full max-w-2xl">
        <CardHeader>
          <CardTitle className="text-2xl">Timelith Setup Wizard</CardTitle>
          <CardDescription>
            Step {step} of 3: {step === 1 ? 'Telegram API' : step === 2 ? 'Server & Database' : 'Admin User'}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={step === 3 ? handleSubmit : (e) => { e.preventDefault(); handleNext(); }}>
            {/* Step 1: Telegram API */}
            {step === 1 && (
              <div className="space-y-4">
                <div className="bg-blue-50 p-4 rounded-lg mb-4">
                  <p className="text-sm text-blue-800">
                    Get your Telegram API credentials from{' '}
                    <a
                      href="https://my.telegram.org"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="underline font-semibold"
                    >
                      my.telegram.org
                    </a>
                  </p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="telegram_app_id">Telegram App ID *</Label>
                  <Input
                    id="telegram_app_id"
                    type="text"
                    placeholder="1234567"
                    value={telegramAppId}
                    onChange={(e) => setTelegramAppId(e.target.value)}
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="telegram_app_hash">Telegram App Hash *</Label>
                  <Input
                    id="telegram_app_hash"
                    type="text"
                    placeholder="0123456789abcdef0123456789abcdef"
                    value={telegramAppHash}
                    onChange={(e) => setTelegramAppHash(e.target.value)}
                    required
                  />
                  <p className="text-xs text-muted-foreground">
                    32-character hexadecimal string
                  </p>
                </div>
              </div>
            )}

            {/* Step 2: Server & Database */}
            {step === 2 && (
              <div className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="server_port">Server Port</Label>
                  <Input
                    id="server_port"
                    type="text"
                    value={serverPort}
                    onChange={(e) => setServerPort(e.target.value)}
                    required
                  />
                  <p className="text-xs text-muted-foreground">
                    Default: 8080
                  </p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="environment">Environment</Label>
                  <select
                    id="environment"
                    className="w-full px-3 py-2 border rounded-md"
                    value={environment}
                    onChange={(e) => setEnvironment(e.target.value)}
                  >
                    <option value="production">Production</option>
                    <option value="development">Development</option>
                  </select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="postgres_password">PostgreSQL Password *</Label>
                  <Input
                    id="postgres_password"
                    type="password"
                    value={postgresPassword}
                    onChange={(e) => setPostgresPassword(e.target.value)}
                    required
                  />
                  <p className="text-xs text-muted-foreground">
                    Password for PostgreSQL database user &quot;timelith&quot;
                  </p>
                </div>

                <div className="bg-yellow-50 p-4 rounded-lg">
                  <p className="text-sm text-yellow-800">
                    <strong>Note:</strong> Make sure PostgreSQL is running and accessible.
                    The database &quot;timelith&quot; with user &quot;timelith&quot; should be created before setup.
                  </p>
                </div>
              </div>
            )}

            {/* Step 3: Admin User */}
            {step === 3 && (
              <div className="space-y-4">
                <div className="bg-green-50 p-4 rounded-lg mb-4">
                  <p className="text-sm text-green-800">
                    Create the first administrator account to access the dashboard
                  </p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="admin_username">Admin Username *</Label>
                  <Input
                    id="admin_username"
                    type="text"
                    value={adminUsername}
                    onChange={(e) => setAdminUsername(e.target.value)}
                    required
                  />
                  <p className="text-xs text-muted-foreground">
                    Minimum 3 characters
                  </p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="admin_password">Admin Password *</Label>
                  <Input
                    id="admin_password"
                    type="password"
                    value={adminPassword}
                    onChange={(e) => setAdminPassword(e.target.value)}
                    required
                  />
                  <p className="text-xs text-muted-foreground">
                    Minimum 6 characters
                  </p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="admin_password_confirm">Confirm Password *</Label>
                  <Input
                    id="admin_password_confirm"
                    type="password"
                    value={adminPasswordConfirm}
                    onChange={(e) => setAdminPasswordConfirm(e.target.value)}
                    required
                  />
                </div>

                <div className="bg-blue-50 p-4 rounded-lg">
                  <p className="text-sm text-blue-800">
                    <strong>Security Keys:</strong> JWT_SECRET and ENCRYPTION_KEY will be
                    generated automatically with cryptographically secure random values.
                  </p>
                </div>
              </div>
            )}

            {error && (
              <Alert variant="destructive" className="mt-4">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <div className="flex justify-between mt-6">
              {step > 1 && (
                <Button type="button" variant="outline" onClick={handleBack}>
                  Back
                </Button>
              )}
              {step < 3 ? (
                <Button type="submit" className="ml-auto">
                  Next
                </Button>
              ) : (
                <Button type="submit" className="ml-auto" disabled={loading}>
                  {loading ? 'Setting up...' : 'Complete Setup'}
                </Button>
              )}
            </div>

            <div className="flex justify-center mt-4 space-x-2">
              {[1, 2, 3].map((s) => (
                <div
                  key={s}
                  className={`h-2 w-12 rounded-full ${
                    s === step ? 'bg-primary' : s < step ? 'bg-primary/50' : 'bg-gray-200'
                  }`}
                />
              ))}
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
