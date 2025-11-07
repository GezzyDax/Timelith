export interface User {
  id: string
  username: string
  created_at: string
  updated_at: string
}

export interface Account {
  id: string
  phone: string
  status: 'active' | 'inactive' | 'error'
  last_login_at?: string
  error_message?: string
  created_at: string
  updated_at: string
}

export interface Template {
  id: string
  name: string
  content: string
  variables: string[]
  created_at: string
  updated_at: string
}

export interface Channel {
  id: string
  name: string
  chat_id: string
  type: 'channel' | 'group' | 'user'
  created_at: string
  updated_at: string
}

export interface Schedule {
  id: string
  name: string
  account_id: string
  template_id: string
  channel_id: string
  cron_expr: string
  timezone: string
  status: 'active' | 'paused' | 'completed'
  next_run_at?: string
  last_run_at?: string
  created_at: string
  updated_at: string
}

export interface JobLog {
  id: string
  schedule_id: string
  status: 'success' | 'failed' | 'retry'
  message?: string
  error?: string
  executed_at: string
  created_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface CreateAccountRequest {
  phone: string
}

export interface CreateTemplateRequest {
  name: string
  content: string
  variables: string[]
}

export interface CreateChannelRequest {
  name: string
  chat_id: string
  type: string
}

export interface CreateScheduleRequest {
  name: string
  account_id: string
  template_id: string
  channel_id: string
  cron_expr: string
  timezone: string
}
