export interface User {
  user_id: string
  email: string
  first_name: string
  last_name: string
  avatar_url: string
  role: string
  status: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  user_id?: string
  profile?: User
}

export interface Address {
  country: string
  city: string
  district: string
  street: string
  building: string
  lat: number
  lng: number
}

export interface Listing {
  id: string
  seller_id: string
  agent_id?: string
  type: number
  status: number
  title: string
  description: string
  price: number
  currency: string
  area_sqm: number
  rooms: number
  floor: number
  floors_total: number
  address: Address
  media_urls: string[]
}

export interface SearchHit {
  listing_id: string
  score: number
  title: string
  price: number
  city: string
  district: string
  area_sqm: number
  rooms: number
  preview_url: string
  lat: number
  lng: number
}

export interface SearchResponse {
  hits: SearchHit[]
  total: number
  page: number
}

export interface Notification {
  id: string
  user_id: string
  title: string
  body: string
  link: string
  is_read: boolean
  created_at: string
}

export const ListingTypeLabels: Record<number, string> = {
  1: 'Квартира',
  2: 'Дом',
  3: 'Коммерческая',
  4: 'Земля',
}

export const ListingStatusLabels: Record<number, string> = {
  1: 'Черновик',
  2: 'Опубликовано',
  3: 'В архиве',
  4: 'Продано',
}