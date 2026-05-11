import { FC, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'

export const HomePage: FC = () => {
  const [query, setQuery] = useState('')
  const [city, setCity] = useState('')
  const navigate = useNavigate()

  const handleSearch = () => {
    const params = new URLSearchParams()
    if (query) params.set('q', query)
    if (city) params.set('city', city)
    navigate(`/listings?${params}`)
  }

  return (
    <div>
      <section className="bg-gradient-to-br from-blue-600 to-blue-800 rounded-2xl p-12 text-white mb-12">
        <h1 className="text-4xl font-bold mb-4">Найдите свою недвижимость</h1>
        <p className="text-blue-100 text-lg mb-8">Тысячи объявлений по всей России</p>
        <div className="bg-white rounded-xl p-4 flex flex-col md:flex-row gap-3 items-center">
  <input
    placeholder="Поиск по ключевым словам..."
    value={query}
    onChange={e => setQuery(e.target.value)}
    onKeyDown={e => e.key === 'Enter' && handleSearch()}
    className="flex-1 w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 text-gray-900"
  />
  <input
    placeholder="Город"
    value={city}
    onChange={e => setCity(e.target.value)}
    onKeyDown={e => e.key === 'Enter' && handleSearch()}
    className="w-full md:w-48 px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 text-gray-900"
  />
  <button
    onClick={handleSearch}
    className="w-full md:w-auto px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors text-sm"
  >
    Найти
  </button>
</div>
      </section>

      <section className="mb-12">
        <h2 className="text-2xl font-semibold text-gray-900 mb-6">Тип недвижимости</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {[
            { label: 'Квартиры',     icon: '🏢', type: 'apartment' },
            { label: 'Дома',         icon: '🏠', type: 'house' },
            { label: 'Коммерческая', icon: '🏪', type: 'commercial' },
            { label: 'Земля',        icon: '🌿', type: 'land' },
          ].map(item => (
            <button
              key={item.type}
              onClick={() => navigate(`/listings?type=${item.type}`)}
              className="bg-white rounded-xl border border-gray-200 p-6 text-center hover:border-blue-400 hover:shadow-md transition-all"
            >
              <div className="text-3xl mb-2">{item.icon}</div>
              <div className="font-medium text-gray-900">{item.label}</div>
            </button>
          ))}
        </div>
      </section>

      <section>
        <h2 className="text-2xl font-semibold text-gray-900 mb-6">Почему мы</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[
            { title: 'Проверенные объявления', desc: 'Каждое объявление проходит модерацию перед публикацией' },
            { title: 'Удобный поиск',          desc: 'Фильтры по цене, площади, району и другим параметрам' },
            { title: 'Быстрые сделки',         desc: 'Заявки на просмотр и оформление сделок онлайн' },
          ].map(item => (
            <div key={item.title} className="bg-white rounded-xl border border-gray-200 p-6">
              <h3 className="font-semibold text-gray-900 mb-2">{item.title}</h3>
              <p className="text-sm text-gray-500">{item.desc}</p>
            </div>
          ))}
        </div>
      </section>
    </div>
  )
}