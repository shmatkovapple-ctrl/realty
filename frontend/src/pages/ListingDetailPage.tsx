import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { listingsApi } from '../api/listings'
import { api } from '../api/client'
import { Button } from '../components/ui/Button'
import { Badge } from '../components/ui/Badge'
import { Card } from '../components/ui/Card'
import { useAuthStore } from '../store/auth'
import { ListingTypeLabels, ListingStatusLabels } from '../types'

export const ListingDetailPage = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { isAuthenticated } = useAuthStore()

  const [showViewingForm, setShowViewingForm] = useState(false)
  const [comment, setComment] = useState('')
  const [viewingLoading, setViewingLoading] = useState(false)
  const [viewingSuccess, setViewingSuccess] = useState(false)
  const [viewingError, setViewingError] = useState('')
  const [inFavorites, setInFavorites] = useState(false)

  const { data, isLoading } = useQuery({
    queryKey: ['listing', id],
    queryFn: () => listingsApi.getById(id!),
    enabled: !!id,
  })

  const handleViewingRequest = async () => {
    setViewingLoading(true)
    setViewingError('')
    try {
      await api.post('/viewings', { listing_id: id, comment })
      setViewingSuccess(true)
      setShowViewingForm(false)
    } catch (err: any) {
      setViewingError(err.message || 'Ошибка отправки заявки')
    } finally {
      setViewingLoading(false)
    }
  }

  const handleFavorite = async () => {
    try {
      if (inFavorites) {
        await api.delete(`/favorites/${id}`)
        setInFavorites(false)
      } else {
        await api.post('/favorites', { listing_id: id })
        setInFavorites(true)
      }
    } catch (err: any) {
      console.error(err)
    }
  }

  if (isLoading) return (
    <div className="animate-pulse space-y-4">
      <div className="h-96 bg-gray-200 rounded-xl" />
      <div className="h-8 bg-gray-200 rounded w-1/2" />
    </div>
  )

  const listing = data?.listing
  if (!listing) return (
    <div className="text-center py-16 text-gray-500">Объявление не найдено</div>
  )

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <div className="lg:col-span-2">
        <div className="bg-gray-100 rounded-xl h-96 flex items-center justify-center mb-6 overflow-hidden">
          {listing.media_urls?.length ? (
            <img src={listing.media_urls[0]} alt={listing.title} className="w-full h-full object-cover" />
          ) : (
            <span className="text-gray-400 text-6xl">🏠</span>
          )}
        </div>

        <div className="flex gap-2 mb-4">
          <Badge variant="blue">{ListingTypeLabels[listing.type] || 'Объект'}</Badge>
          <Badge variant={listing.status === 2 ? 'green' : 'gray'}>
            {ListingStatusLabels[listing.status] || 'Статус'}
          </Badge>
        </div>

        <h1 className="text-2xl font-bold text-gray-900 mb-2">{listing.title}</h1>
        <div className="text-3xl font-bold text-blue-600 mb-6">
          {new Intl.NumberFormat('ru-RU').format(listing.price)} ₽
        </div>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
          {listing.area_sqm > 0 && (
            <div className="bg-gray-50 rounded-lg p-3 text-center">
              <div className="text-lg font-semibold text-gray-900">{listing.area_sqm} м²</div>
              <div className="text-xs text-gray-500">Площадь</div>
            </div>
          )}
          {listing.rooms > 0 && (
            <div className="bg-gray-50 rounded-lg p-3 text-center">
              <div className="text-lg font-semibold text-gray-900">{listing.rooms}</div>
              <div className="text-xs text-gray-500">Комнат</div>
            </div>
          )}
          {listing.floor > 0 && (
            <div className="bg-gray-50 rounded-lg p-3 text-center">
              <div className="text-lg font-semibold text-gray-900">{listing.floor}/{listing.floors_total}</div>
              <div className="text-xs text-gray-500">Этаж</div>
            </div>
          )}
          {listing.address?.city && (
            <div className="bg-gray-50 rounded-lg p-3 text-center">
              <div className="text-lg font-semibold text-gray-900">{listing.address.city}</div>
              <div className="text-xs text-gray-500">Город</div>
            </div>
          )}
        </div>

        {listing.description && (
          <div className="mb-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-3">Описание</h2>
            <p className="text-gray-600 leading-relaxed whitespace-pre-wrap">{listing.description}</p>
          </div>
        )}

        {listing.address && (
          <div>
            <h2 className="text-lg font-semibold text-gray-900 mb-3">Адрес</h2>
            <p className="text-gray-600">
              {[listing.address.country, listing.address.city, listing.address.district,
                listing.address.street, listing.address.building].filter(Boolean).join(', ')}
            </p>
          </div>
        )}
      </div>

      <div>
        <Card className="p-6 sticky top-24">
          <div className="text-2xl font-bold text-blue-600 mb-4">
            {new Intl.NumberFormat('ru-RU').format(listing.price)} ₽
          </div>

          {isAuthenticated ? (
            <div className="flex flex-col gap-3">
              {viewingSuccess ? (
                <div className="bg-green-50 text-green-700 text-sm px-4 py-3 rounded-lg text-center">
                  Заявка отправлена! Продавец свяжется с вами.
                </div>
              ) : showViewingForm ? (
                <div className="flex flex-col gap-3">
                  <div className="flex flex-col gap-1">
                    <label className="text-sm font-medium text-gray-700">Комментарий</label>
                    <textarea
                      rows={3}
                      placeholder="Удобное время для просмотра, вопросы..."
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none focus:border-blue-500 resize-none"
                      value={comment}
                      onChange={e => setComment(e.target.value)}
                    />
                  </div>
                  {viewingError && (
                    <div className="bg-red-50 text-red-600 text-sm px-3 py-2 rounded-lg">
                      {viewingError}
                    </div>
                  )}
                  <Button className="w-full" loading={viewingLoading} onClick={handleViewingRequest}>
                    Отправить заявку
                  </Button>
                  <Button variant="ghost" className="w-full" onClick={() => setShowViewingForm(false)}>
                    Отмена
                  </Button>
                </div>
              ) : (
                <>
                  <Button className="w-full" onClick={() => setShowViewingForm(true)}>
                    Записаться на просмотр
                  </Button>
                  <Button variant="secondary" className="w-full" onClick={handleFavorite}>
                    {inFavorites ? '❤️ В избранном' : '🤍 В избранное'}
                  </Button>
                </>
              )}
            </div>
          ) : (
            <div className="flex flex-col gap-3">
              <p className="text-sm text-gray-500 text-center mb-2">
                Войдите чтобы связаться с продавцом
              </p>
              <Button className="w-full" onClick={() => navigate('/login')}>Войти</Button>
              <Button variant="secondary" className="w-full" onClick={() => navigate('/register')}>
                Зарегистрироваться
              </Button>
            </div>
          )}
        </Card>
      </div>
    </div>
  )
}