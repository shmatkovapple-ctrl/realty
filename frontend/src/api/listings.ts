import { api } from './client'
import { Listing, SearchResponse } from '../types'

export const listingsApi = {
  search: (params: Record<string, string | number>) => {
    const query = new URLSearchParams()
    Object.entries(params).forEach(([k, v]) => {
      if (v !== undefined && v !== '' && v !== 0) query.set(k, String(v))
    })
    return api.get<SearchResponse>(`/listings?${query}`)
  },

  autocomplete: (q: string, field = 'city') =>
    api.get<{ suggestions: string[] }>(
      `/listings/autocomplete?q=${encodeURIComponent(q)}&field=${field}`
    ),

  getById: (id: string) =>
    api.get<{ listing: Listing }>(`/listings/${id}`),

  create: (data: Partial<Listing>) =>
    api.post<{ listing: Listing }>('/listings', data),

  update: (id: string, data: Partial<Listing>) =>
    api.put<{ listing: Listing }>(`/listings/${id}`, data),

  delete: (id: string) =>
    api.delete(`/listings/${id}`),

  publish: (id: string) =>
    api.post(`/listings/${id}/publish`, {}),

  getMyListings: (page = 1, limit = 20) =>
    api.get<{ listings: Listing[]; total: number; page: number }>(
      `/listings/mine?page=${page}&limit=${limit}`
    ),

  getUploadUrl: (id: string, filename: string, content_type: string) =>
    api.post<{ upload_url: string; file_url: string }>(
      `/listings/${id}/upload-url`,
      { filename, content_type }
    ),
}