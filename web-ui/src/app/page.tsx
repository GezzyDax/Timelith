'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import { Navbar } from '@/components/layout/Navbar'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { api } from '@/lib/api'

export default function DashboardPage() {
  const router = useRouter()
  const [checking, setChecking] = useState(true)

  useEffect(() => {
    const checkSetup = async () => {
      try {
        const status = await api.checkSetupStatus()
        if (status.setup_required) {
          router.push('/setup')
          return
        }
      } catch (err) {
        // If setup check fails, might be in setup mode or API is down
        console.error('Setup check failed:', err)
      }

      const token = localStorage.getItem('token')
      if (!token) {
        router.push('/login')
      }
      setChecking(false)
    }

    checkSetup()
  }, [router])

  if (checking) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p>Loading...</p>
      </div>
    )
  }

  const { data: accounts } = useQuery({
    queryKey: ['accounts'],
    queryFn: () => api.getAccounts(),
  })

  const { data: templates } = useQuery({
    queryKey: ['templates'],
    queryFn: () => api.getTemplates(),
  })

  const { data: schedules } = useQuery({
    queryKey: ['schedules'],
    queryFn: () => api.getSchedules(),
  })

  const { data: logs } = useQuery({
    queryKey: ['logs'],
    queryFn: () => api.getAllLogs(),
  })

  const stats = [
    { name: 'Total Accounts', value: accounts?.length || 0, href: '/accounts' },
    { name: 'Templates', value: templates?.length || 0, href: '/templates' },
    { name: 'Active Schedules', value: schedules?.filter(s => s.status === 'active').length || 0, href: '/schedules' },
    { name: 'Recent Jobs', value: logs?.length || 0, href: '/logs' },
  ]

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <main className="container mx-auto py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <p className="text-muted-foreground mt-2">
            Welcome to Timelith - Telegram Account Manager
          </p>
        </div>

        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {stats.map((stat) => (
            <Card key={stat.name}>
              <CardHeader className="pb-2">
                <CardDescription>{stat.name}</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-3xl font-bold">{stat.value}</div>
              </CardContent>
            </Card>
          ))}
        </div>

        <div className="mt-8 grid gap-4 md:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle>Recent Activity</CardTitle>
              <CardDescription>Latest job executions</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                {logs?.slice(0, 5).map((log) => (
                  <div key={log.id} className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">
                      {new Date(log.executed_at).toLocaleString()}
                    </span>
                    <span
                      className={
                        log.status === 'success'
                          ? 'text-green-600'
                          : log.status === 'failed'
                          ? 'text-red-600'
                          : 'text-yellow-600'
                      }
                    >
                      {log.status}
                    </span>
                  </div>
                ))}
                {(!logs || logs.length === 0) && (
                  <p className="text-sm text-muted-foreground">No recent activity</p>
                )}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>System Status</CardTitle>
              <CardDescription>Current system health</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <div className="flex items-center justify-between text-sm">
                  <span>Active Accounts</span>
                  <span className="text-green-600">
                    {accounts?.filter(a => a.status === 'active').length || 0}
                  </span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span>Running Schedules</span>
                  <span className="text-green-600">
                    {schedules?.filter(s => s.status === 'active').length || 0}
                  </span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span>Failed Jobs (24h)</span>
                  <span className="text-red-600">
                    {logs?.filter(l => l.status === 'failed').length || 0}
                  </span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  )
}
