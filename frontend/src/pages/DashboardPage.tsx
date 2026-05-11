import { FC } from 'react'
import { Link } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { authApi } from '../api/auth'
import { notificationsApi } from '../api/notifications'
import { listingsApi } from '../api/listings'
import { Card } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { Badge } from '../components/ui/Badge'
import { Listing, ListingTypeLabels, ListingStatusLabels } from '../types'

export const DashboardPage: FC = () => {
  const { data: profileData } = useQuery({
    queryKey: ['profile'],
    queryFn: authApi.getProfile,
  })

  const { data: notifData } = useQuery({
    queryKey: ['notifications'],
    queryFn: () => notificationsApi.list({ limit: 5 }),
  })

  const profile = profileData?.profile
  const isSellerOrAgent = profile?.role === 'seller' || profile?.role === 'agent'
  const queryClient = useQueryClient()

  const { data: myListingsData, isLoading: listingsLoading } = useQuery({
    queryKey: ['my-listings'],
    queryFn: () => listingsApi.getMyListings(),
    enabled: isSellerOrAgent,
  })

  const publishMutation = useMutation({
    mutationFn: (id: string) => listingsApi.publish(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['my-listings'] }),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => listingsApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['my-listings'] }),
  })

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Личный кабинет</h1>
        <Link to="/listings/new">
          <Button>Разместить объявление</Button>
        </Link>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <Card className="p-6">
          <div className="flex items-center gap-4">
            <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
              <span className="text-blue-600 font-bold text-lg">
                {profile?.first_name?.[0] || profile?.email?.[0]?.toUpperCase() || '?'}
              </span>
            </div>
            <div>
              <div className="font-semibold text-gray-900">
                {profile?.first_name
                  ? `${profile.first_name} ${profile.last_name}`
                  : 'Пользователь'}
              </div>
              <div className="text-sm text-gray-500">{profile?.email}</div>
              <Badge variant="blue" className="mt-1">
                {profile?.role === 'seller' ? 'Продавец'
                  : profile?.role === 'agent' ? 'Агент' : 'Покупатель'}
              </Badge>
            </div>
          </div>
        </Card>

        <Card className="p-6">
          <div className="text-sm text-gray-500 mb-1">Уведомления</div>
          <div className="text-3xl font-bold text-gray-900">{notifData?.unread || 0}</div>
          <div className="text-sm text-gray-500">непрочитанных</div>
        </Card>

        <Card className="p-6">
          <div className="text-sm text-gray-500 mb-1">Избранное</div>
          <div className="text-3xl font-bold text-gray-900">—</div>
          <Link to="/favorites" className="text-sm text-blue-600 hover:underline">Смотреть</Link>
        </Card>
      </div>

      {notifData?.notifications?.length ? (
        <Card className="p-6">
          <h2 className="font-semibold text-gray-900 mb-4">Последние уведомления</h2>
          <div className="flex flex-col divide-y divide-gray-100">
            {notifData.notifications.map(n => (
              <div key={n.id} className={`py-3 flex justify-between items-start ${!n.is_read ? 'font-medium' : ''}`}>
                <div>
                  <div className="text-sm text-gray-900">{n.title}</div>
                  <div className="text-xs text-gray-500 mt-0.5">{n.body}</div>
                </div>
                {!n.is_read && <Badge variant="blue">Новое</Badge>}
              </div>
            ))}
          </div>
        </Card>
      ) : null}

      {isSellerOrAgent && (
        <Card className="p-6 mt-6">
          <h2 className="font-semibold text-gray-900 mb-4">Мои объявления</h2>
          {listingsLoading ? (
            <div className="flex flex-col gap-3">
              {[...Array(3)].map((_, i) => (
                <div key={i} className="h-14 bg-gray-100 rounded-lg animate-pulse" />
              ))}
            </div>
          ) : myListingsData?.listings?.length ? (
            <div className="flex flex-col divide-y divide-gray-100">
              {myListingsData.listings.map((l: Listing) => (
                <MyListingRow
                  key={l.id}
                  listing={l}
                  onPublish={() => publishMutation.mutate(l.id)}
                  onDelete={() => deleteMutation.mutate(l.id)}
                />
              ))}
            </div>
          ) : (
            <div className="text-center py-6 text-sm text-gray-500">У вас пока нет объявлений</div>
          )}
        </Card>
      )}
    </div>
  )
}

const statusVariant: Record<number, 'gray' | 'green' | 'yellow' | 'red'> = {
  1: 'yellow', 2: 'green', 3: 'gray', 4: 'red',
}

const MyListingRow: FC<{ listing: Listing; onPublish: () => void; onDelete: () => void }> = ({
  listing, onPublish, onDelete,
}) => (
  <div className="py-3 flex items-center justify-between gap-4">
    <Link to={`/listings/${listing.id}`} className="flex-1 min-w-0">
      <div className="font-medium text-gray-900 truncate">{listing.title}</div>
      <div className="flex items-center gap-2 mt-0.5 flex-wrap">
        <Badge variant={statusVariant[listing.status] ?? 'gray'}>
          {ListingStatusLabels[listing.status as keyof typeof ListingStatusLabels] ?? 'Неизвестно'}
        </Badge>
        <span className="text-sm text-gray-500">
          {new Intl.NumberFormat('ru-RU').format(listing.price)} ₽
        </span>
        {listing.type > 0 && (
          <span className="text-xs text-gray-400">
            {ListingTypeLabels[listing.type as keyof typeof ListingTypeLabels]}
          </span>
        )}
      </div>
    </Link>
    <div className="flex gap-2 flex-shrink-0">
      {listing.status === 1 && (
        <Button size="sm" variant="secondary" onClick={onPublish}>Опубликовать</Button>
      )}
      <Button size="sm" variant="danger" onClick={onDelete}>Удалить</Button>
    </div>
  </div>
)