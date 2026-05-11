import { FC, useState } from 'react'
import { useNavigate, useSearchParams, Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { listingsApi } from '../api/listings'
import { SearchHit } from '../types'
import { Card } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Badge } from '../components/ui/Badge'

export const ListingsPage: FC = () => {
  const [searchParams, setSearchParams] = useSearchParams()
  const navigate = useNavigate()

  const [filters, setFilters] = useState({
    q:         searchParams.get('q') || '',
    city:      searchParams.get('city') || '',
    type:      searchParams.get('type') || '',
    price_min: searchParams.get('price_min') || '',
    price_max: searchParams.get('price_max') || '',
    rooms:     searchParams.get('rooms') || '',
  })

  const params: Record<string, string | number> = { page: 1, limit: 20 }
  Object.entries(filters).forEach(([k, v]) => { if (v) params[k] = v })

  const { data, isLoading } = useQuery({
    queryKey: ['listings', filters],
    queryFn: () => listingsApi.search(params),
  })

  const applyFilters = () => {
    const p = new URLSearchParams()
    Object.entries(filters).forEach(([k, v]) => { if (v) p.set(k, v) })
    setSearchParams(p)
  }

  const resetFilters = () => {
    setFilters({ q: '', city: '', type: '', price_min: '', price_max: '', rooms: '' })
    setSearchParams({})
  }

  return (
    <div className="flex gap-8">
      <aside className="w-72 flex-shrink-0">
        <Card className="p-5 sticky top-24">
          <h2 className="font-semibold text-gray-900 mb-4">Фильтры</h2>
          <div className="flex flex-col gap-4">
            <Input label="Поиск" placeholder="Квартира, дом..."
              value={filters.q} onChange={e => setFilters({...filters, q: e.target.value})} />
            <Input label="Город" placeholder="Москва"
              value={filters.city} onChange={e => setFilters({...filters, city: e.target.value})} />
            <div className="flex flex-col gap-1">
              <label className="text-sm font-medium text-gray-700">Тип</label>
              <select
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none focus:border-blue-500"
                value={filters.type} onChange={e => setFilters({...filters, type: e.target.value})}
              >
                <option value="">Все типы</option>
                <option value="apartment">Квартира</option>
                <option value="house">Дом</option>
                <option value="commercial">Коммерческая</option>
                <option value="land">Земля</option>
              </select>
            </div>
            <div className="grid grid-cols-2 gap-2">
              <Input label="Цена от" placeholder="0"
                value={filters.price_min} onChange={e => setFilters({...filters, price_min: e.target.value})} />
              <Input label="Цена до" placeholder="∞"
                value={filters.price_max} onChange={e => setFilters({...filters, price_max: e.target.value})} />
            </div>
            <div className="flex flex-col gap-1">
              <label className="text-sm font-medium text-gray-700">Комнат</label>
              <select
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none focus:border-blue-500"
                value={filters.rooms} onChange={e => setFilters({...filters, rooms: e.target.value})}
              >
                <option value="">Любое</option>
                <option value="1">1</option>
                <option value="2">2</option>
                <option value="3">3</option>
                <option value="4">4+</option>
              </select>
            </div>
            <Button onClick={applyFilters} className="w-full">Применить</Button>
            <Button variant="ghost" onClick={resetFilters} className="w-full">Сбросить</Button>
          </div>
        </Card>
      </aside>

      <div className="flex-1">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-2xl font-bold text-gray-900">
            {isLoading ? 'Загрузка...' : `Найдено ${data?.total || 0} объявлений`}
          </h1>
          <Button onClick={() => navigate('/map')} variant="secondary" size="sm">На карте</Button>
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="bg-white rounded-xl border border-gray-200 h-64 animate-pulse" />
            ))}
          </div>
        ) : data?.hits?.length ? (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {data.hits.map(hit => <ListingCard key={hit.listing_id} hit={hit} />)}
          </div>
        ) : (
          <div className="text-center py-16 text-gray-500">
            <div className="text-4xl mb-4">🔍</div>
            <div className="text-lg font-medium">Объявлений не найдено</div>
            <div className="text-sm mt-2">Попробуйте изменить параметры поиска</div>
          </div>
        )}
      </div>
    </div>
  )
}

const ListingCard: FC<{ hit: SearchHit }> = ({ hit }) => (
  <Link to={`/listings/${hit.listing_id}`}>
    <Card className="overflow-hidden hover:shadow-md transition-shadow">
      <div className="bg-gray-100 h-48 flex items-center justify-center">
        {hit.preview_url ? (
          <img src={hit.preview_url} alt={hit.title} className="w-full h-full object-cover" />
        ) : (
          <span className="text-gray-400 text-4xl">🏠</span>
        )}
      </div>
      <div className="p-4">
        <div className="text-xl font-bold text-blue-600 mb-1">
          {new Intl.NumberFormat('ru-RU').format(hit.price)} ₽
        </div>
        <div className="font-medium text-gray-900 mb-1 line-clamp-2">{hit.title}</div>
        <div className="text-sm text-gray-500 mb-3">
          {hit.city}{hit.district ? `, ${hit.district}` : ''}
        </div>
        <div className="flex gap-2 flex-wrap">
          {hit.rooms > 0 && <Badge variant="blue">{hit.rooms} комн.</Badge>}
          {hit.area_sqm > 0 && <Badge variant="gray">{hit.area_sqm} м²</Badge>}
        </div>
      </div>
    </Card>
  </Link>
)