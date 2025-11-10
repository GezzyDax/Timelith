'use client'
export const dynamic = 'force-dynamic'


import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Navbar } from '@/components/layout/Navbar'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { api } from '@/lib/api'
import { Plus, Trash2 } from 'lucide-react'

export default function TemplatesPage() {
  const queryClient = useQueryClient()
  const [showCreate, setShowCreate] = useState(false)
  const [name, setName] = useState('')
  const [content, setContent] = useState('')

  const { data: templates } = useQuery({
    queryKey: ['templates'],
    queryFn: () => api.getTemplates(),
  })

  const createMutation = useMutation({
    mutationFn: () => api.createTemplate({ name, content, variables: [] }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['templates'] })
      setShowCreate(false)
      setName('')
      setContent('')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.deleteTemplate(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['templates'] })
    },
  })

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <main className="container mx-auto py-8">
        <div className="flex items-center justify-between mb-8">
          <h1 className="text-3xl font-bold">Message Templates</h1>
          <Button onClick={() => setShowCreate(!showCreate)}>
            <Plus className="mr-2 h-4 w-4" />
            Add Template
          </Button>
        </div>

        {showCreate && (
          <Card className="mb-6">
            <CardHeader>
              <CardTitle>Create Template</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={(e) => { e.preventDefault(); createMutation.mutate(); }} className="space-y-4">
                <div className="space-y-2">
                  <Label>Name</Label>
                  <Input value={name} onChange={(e) => setName(e.target.value)} required />
                </div>
                <div className="space-y-2">
                  <Label>Content</Label>
                  <textarea
                    className="w-full min-h-[100px] rounded-md border border-input bg-background px-3 py-2"
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    required
                  />
                </div>
                <Button type="submit">Create</Button>
              </form>
            </CardContent>
          </Card>
        )}

        <div className="grid gap-4">
          {templates?.map((template) => (
            <Card key={template.id}>
              <CardContent className="flex items-start justify-between p-6">
                <div className="flex-1">
                  <h3 className="font-semibold">{template.name}</h3>
                  <p className="text-sm text-muted-foreground mt-2">{template.content}</p>
                </div>
                <Button variant="destructive" size="icon" onClick={() => deleteMutation.mutate(template.id)}>
                  <Trash2 className="h-4 w-4" />
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      </main>
    </div>
  )
}
