import { api } from './client'
import { AuthResponse, User } from '../types'

export const authApi = {
  register: (data: { email: string; password: string; phone: string; role: string }) =>
    api.post<AuthResponse>('/auth/register', data),

  login: (data: { email: string; password: string }) =>
    api.post<AuthResponse>('/auth/login', data),

  logout: (refresh_token: string) =>
    api.post('/auth/logout', { refresh_token }),

  refresh: (refresh_token: string) =>
    api.post<AuthResponse>('/auth/refresh', { refresh_token }),

  getProfile: () =>
    api.get<{ profile: User }>('/profile'),

  updateProfile: (data: { first_name: string; last_name: string; bio: string }) =>
    api.put<{ profile: User }>('/profile', data),
}