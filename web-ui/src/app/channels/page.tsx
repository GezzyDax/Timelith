'use client'

import { Navbar } from '@/components/layout/Navbar'

export default function ChannelsPage() {
  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <main className="container mx-auto py-8">
        <h1 className="text-3xl font-bold">Channels</h1>
        <p className="text-muted-foreground mt-2">Manage Telegram channels and groups</p>
      </main>
    </div>
  )
}
