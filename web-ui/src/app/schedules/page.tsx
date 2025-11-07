'use client'

import { useQuery } from '@tanstack/react-query'
import { Navbar } from '@/components/layout/Navbar'
import { Card, CardContent } from '@/components/ui/card'
import { api } from '@/lib/api'
import { formatDate } from '@/lib/utils'

export default function SchedulesPage() {
  const { data: schedules } = useQuery({
    queryKey: ['schedules'],
    queryFn: () => api.getSchedules(),
  })

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <main className="container mx-auto py-8">
        <h1 className="text-3xl font-bold mb-8">Schedules</h1>
        <div className="grid gap-4">
          {schedules?.map((schedule) => (
            <Card key={schedule.id}>
              <CardContent className="p-6">
                <h3 className="font-semibold">{schedule.name}</h3>
                <p className="text-sm text-muted-foreground">Cron: {schedule.cron_expr}</p>
                <p className="text-sm text-muted-foreground">Next run: {formatDate(schedule.next_run_at)}</p>
                <span className={`inline-flex mt-2 px-2.5 py-0.5 rounded-full text-xs font-medium
                  ${schedule.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}`}>
                  {schedule.status}
                </span>
              </CardContent>
            </Card>
          ))}
        </div>
      </main>
    </div>
  )
}
