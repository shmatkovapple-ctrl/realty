import { api } from './client'
import { Notification } from '../types'

export const notificationsApi = {
  list: (params?: { unread_only?: boolean; page?: number; limit?: number }) => {
    const query = new URLSearchParams()
    if (params?.unread_only) query.set('unread_only', 'true')
    if (params?.page) query.set('page', String(params.page))
    if (params?.limit) query.set('limit', String(params.limit))
    return api.get<{ notifications: Notification[]; total: number; unread: number }>(
      `/notifications?${query}`
    )
  },

  markAsRead: (id: string) =>
    api.put(`/notifications/${id}/read`, {}),

  markAllAsRead: () =>
    api.put('/notifications/read-all', {}),
}