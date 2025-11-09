import axios, { AxiosInstance } from 'axios'
import type {
  Account,
  Channel,
  CreateAccountRequest,
  CreateChannelRequest,
  CreateScheduleRequest,
  CreateTemplateRequest,
  JobLog,
  LoginRequest,
  LoginResponse,
  Schedule,
  Template,
  User,
} from '@/types'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

class ApiClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: `${API_URL}/api`,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    // Add token to requests
    this.client.interceptors.request.use((config) => {
      const token = localStorage.getItem('token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    })

    // Handle 401 errors
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('token')
          window.location.href = '/login'
        }
        return Promise.reject(error)
      }
    )
  }

  // Auth
  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await this.client.post<LoginResponse>('/auth/login', data)
    return response.data
  }

  async register(data: LoginRequest): Promise<User> {
    const response = await this.client.post<User>('/auth/register', data)
    return response.data
  }

  // Accounts
  async getAccounts(): Promise<Account[]> {
    const response = await this.client.get<Account[]>('/accounts')
    return response.data
  }

  async getAccount(id: string): Promise<Account> {
    const response = await this.client.get<Account>(`/accounts/${id}`)
    return response.data
  }

  async createAccount(data: CreateAccountRequest): Promise<Account> {
    const response = await this.client.post<Account>('/accounts', data)
    return response.data
  }

  async deleteAccount(id: string): Promise<void> {
    await this.client.delete(`/accounts/${id}`)
  }

  // Templates
  async getTemplates(): Promise<Template[]> {
    const response = await this.client.get<Template[]>('/templates')
    return response.data
  }

  async getTemplate(id: string): Promise<Template> {
    const response = await this.client.get<Template>(`/templates/${id}`)
    return response.data
  }

  async createTemplate(data: CreateTemplateRequest): Promise<Template> {
    const response = await this.client.post<Template>('/templates', data)
    return response.data
  }

  async updateTemplate(id: string, data: CreateTemplateRequest): Promise<Template> {
    const response = await this.client.put<Template>(`/templates/${id}`, data)
    return response.data
  }

  async deleteTemplate(id: string): Promise<void> {
    await this.client.delete(`/templates/${id}`)
  }

  // Channels
  async getChannels(): Promise<Channel[]> {
    const response = await this.client.get<Channel[]>('/channels')
    return response.data
  }

  async getChannel(id: string): Promise<Channel> {
    const response = await this.client.get<Channel>(`/channels/${id}`)
    return response.data
  }

  async createChannel(data: CreateChannelRequest): Promise<Channel> {
    const response = await this.client.post<Channel>('/channels', data)
    return response.data
  }

  async updateChannel(id: string, data: CreateChannelRequest): Promise<Channel> {
    const response = await this.client.put<Channel>(`/channels/${id}`, data)
    return response.data
  }

  async deleteChannel(id: string): Promise<void> {
    await this.client.delete(`/channels/${id}`)
  }

  // Schedules
  async getSchedules(): Promise<Schedule[]> {
    const response = await this.client.get<Schedule[]>('/schedules')
    return response.data
  }

  async getSchedule(id: string): Promise<Schedule> {
    const response = await this.client.get<Schedule>(`/schedules/${id}`)
    return response.data
  }

  async createSchedule(data: CreateScheduleRequest): Promise<Schedule> {
    const response = await this.client.post<Schedule>('/schedules', data)
    return response.data
  }

  async updateSchedule(id: string, data: Partial<CreateScheduleRequest>): Promise<Schedule> {
    const response = await this.client.put<Schedule>(`/schedules/${id}`, data)
    return response.data
  }

  async updateScheduleStatus(id: string, status: string): Promise<Schedule> {
    const response = await this.client.patch<Schedule>(`/schedules/${id}/status`, { status })
    return response.data
  }

  async deleteSchedule(id: string): Promise<void> {
    await this.client.delete(`/schedules/${id}`)
  }

  async getScheduleLogs(id: string): Promise<JobLog[]> {
    const response = await this.client.get<JobLog[]>(`/schedules/${id}/logs`)
    return response.data
  }

  // Logs
  async getAllLogs(): Promise<JobLog[]> {
    const response = await this.client.get<JobLog[]>('/logs')
    return response.data
  }

  // Setup
  async checkSetupStatus(): Promise<{ setup_required: boolean }> {
    const response = await this.client.get<{ setup_required: boolean }>('/setup/status')
    return response.data
  }

  async performSetup(data: {
    telegram_app_id: string
    telegram_app_hash: string
    server_port: string
    postgres_password: string
    admin_username: string
    admin_password: string
    environment: string
  }): Promise<{ success: boolean; message: string }> {
    const response = await this.client.post<{ success: boolean; message: string }>('/setup', data)
    return response.data
  }

  // New 3-stage setup
  async setupDatabase(data: {
    use_docker_database: boolean
    database_url?: string
  }): Promise<{ success: boolean; message: string }> {
    const response = await this.client.post<{ success: boolean; message: string }>('/setup/database', data)
    return response.data
  }

  async setupAdmin(data: {
    username: string
    password: string
  }): Promise<{ success: boolean; message: string }> {
    const response = await this.client.post<{ success: boolean; message: string }>('/setup/admin', data)
    return response.data
  }

  async setupComplete(data: {
    telegram_app_id: string
    telegram_app_hash: string
  }): Promise<{ success: boolean; message: string }> {
    const response = await this.client.post<{ success: boolean; message: string }>('/setup/complete', data)
    return response.data
  }

  // Settings Management
  async getAllSettings(): Promise<any[]> {
    const response = await this.client.get<any[]>('/settings')
    return response.data
  }

  async getSettingsByCategory(category: string): Promise<any[]> {
    const response = await this.client.get<any[]>(`/settings/category/${category}`)
    return response.data
  }

  async createSetting(data: {
    key: string
    value: string
    encrypted: boolean
    category: string
    description?: string
  }): Promise<{ success: boolean; message: string }> {
    const response = await this.client.post<{ success: boolean; message: string }>('/settings', data)
    return response.data
  }

  async updateSetting(key: string, value: string): Promise<{ success: boolean; message: string }> {
    const response = await this.client.put<{ success: boolean; message: string }>(`/settings/${key}`, { value })
    return response.data
  }

  async deleteSetting(key: string): Promise<{ success: boolean; message: string }> {
    const response = await this.client.delete<{ success: boolean; message: string }>(`/settings/${key}`)
    return response.data
  }

  // User Management
  async getUsers(): Promise<User[]> {
    const response = await this.client.get<User[]>('/users')
    return response.data
  }

  async getUser(id: string): Promise<User> {
    const response = await this.client.get<User>(`/users/${id}`)
    return response.data
  }

  async createUser(data: { username: string; password: string }): Promise<User> {
    const response = await this.client.post<User>('/users', data)
    return response.data
  }

  async updateUser(id: string, data: { username?: string; password?: string }): Promise<User> {
    const response = await this.client.put<User>(`/users/${id}`, data)
    return response.data
  }

  async deleteUser(id: string): Promise<{ success: boolean; message: string }> {
    const response = await this.client.delete<{ success: boolean; message: string }>(`/users/${id}`)
    return response.data
  }
}

export const api = new ApiClient()
