import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { listingsApi } from '../api/listings'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Card } from '../components/ui/Card'

const compressImage = (file: File, maxWidth = 1280, maxHeight = 960, quality = 0.85): Promise<Blob> => {
  return new Promise((resolve) => {
    const img = new Image()
    const url = URL.createObjectURL(file)
    img.onload = () => {
      URL.revokeObjectURL(url)
      let { width, height } = img

      if (width > maxWidth || height > maxHeight) {
        const ratio = Math.min(maxWidth / width, maxHeight / height)
        width = Math.round(width * ratio)
        height = Math.round(height * ratio)
      }

      const canvas = document.createElement('canvas')
      canvas.width = width
      canvas.height = height

      const ctx = canvas.getContext('2d')!
      ctx.drawImage(img, 0, 0, width, height)

      canvas.toBlob(
        (blob) => resolve(blob!),
        'image/jpeg',
        quality
      )
    }
    img.src = url
  })
}

export const CreateListingPage = () => {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [createdId, setCreatedId] = useState<string | null>(null)
  const [photoLoading, setPhotoLoading] = useState(false)
  const [form, setForm] = useState({
    title: '',
    description: '',
    price: '',
    area_sqm: '',
    rooms: '',
    floor: '',
    floors_total: '',
    type: '1',
    country: 'Россия',
    city: '',
    district: '',
    street: '',
    building: '',
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const res = await listingsApi.create({
        type: parseInt(form.type),
        title: form.title,
        description: form.description,
        price: parseFloat(form.price),
        area_sqm: parseFloat(form.area_sqm),
        rooms: parseInt(form.rooms),
        floor: parseInt(form.floor),
        floors_total: parseInt(form.floors_total),
        address: {
          country: form.country,
          city: form.city,
          district: form.district,
          street: form.street,
          building: form.building,
          lat: 0,
          lng: 0,
        },
      })
      setCreatedId(res.listing.id)
    } catch (err: any) {
      setError(err.message || 'Ошибка создания объявления')
    } finally {
      setLoading(false)
    }
  }

  const handlePhotoUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files
    if (!files?.length || !createdId) return
    setPhotoLoading(true)
    try {
      const mediaUrls: string[] = []
      const token = localStorage.getItem('access_token')

      for (const file of Array.from(files)) {
        const compressed = await compressImage(file)
        const filename = file.name.replace(/\.[^.]+$/, '') + '.jpg'

        const formData = new FormData()
        formData.append('photo', compressed, filename)

        const res = await fetch(
          `http://localhost:8080/api/v1/listings/${createdId}/upload`,
          {
            method: 'POST',
            headers: { Authorization: `Bearer ${token}` },
            body: formData,
          }
        )
        const data = await res.json()
        if (data.file_url) mediaUrls.push(data.file_url)
      }

      if (mediaUrls.length > 0) {
        await listingsApi.update(createdId, { media_urls: mediaUrls })
        await listingsApi.publish(createdId)
      }

      navigate(`/listings/${createdId}`)
    } catch {
      navigate(`/listings/${createdId}`)
    } finally {
      setPhotoLoading(false)
    }
  }

  const f = (field: string) => ({
    value: form[field as keyof typeof form],
    onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) =>
      setForm({ ...form, [field]: e.target.value }),
  })

  if (createdId) {
    return (
      <div className="max-w-2xl mx-auto">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">Объявление создано</h1>
        <p className="text-gray-500 mb-6">
          Добавьте фотографии чтобы объявление выглядело привлекательнее
        </p>

        <Card className="p-6">
          <h2 className="font-semibold text-gray-900 mb-4">Загрузить фотографии</h2>
          <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center">
            <div className="text-4xl mb-2">📷</div>
            <div className="text-sm text-gray-500 mb-4">
              Выберите одно или несколько фото объекта.
              Фото автоматически сожмётся до 1280×960px
            </div>
            <input
              type="file"
              accept="image/*"
              multiple
              className="hidden"
              id="photo-upload"
              onChange={handlePhotoUpload}
            />
            <Button
              type="button"
              variant="secondary"
              loading={photoLoading}
              onClick={() => document.getElementById('photo-upload')?.click()}
            >
              {photoLoading ? 'Загрузка...' : 'Выбрать фото'}
            </Button>
          </div>
          <Button
            type="button"
            variant="ghost"
            className="w-full mt-3"
            onClick={async () => {
              try { await listingsApi.publish(createdId!) } catch {}
              navigate(`/listings/${createdId}`)
            }}
          >
            Пропустить (без фото)
          </Button>
        </Card>
      </div>
    )
  }

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Разместить объявление</h1>

      <form onSubmit={handleSubmit} className="flex flex-col gap-6">
        <Card className="p-6">
          <h2 className="font-semibold text-gray-900 mb-4">Основная информация</h2>
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-1">
              <label className="text-sm font-medium text-gray-700">Тип недвижимости</label>
              <select
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none focus:border-blue-500"
                {...f('type')}
              >
                <option value="1">Квартира</option>
                <option value="2">Дом</option>
                <option value="3">Коммерческая</option>
                <option value="4">Земля</option>
              </select>
            </div>
            <Input
              label="Заголовок объявления"
              placeholder="Уютная квартира в центре города"
              required
              {...f('title')}
            />
            <div className="flex flex-col gap-1">
              <label className="text-sm font-medium text-gray-700">Описание</label>
              <textarea
                rows={4}
                placeholder="Подробное описание объекта..."
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none focus:border-blue-500 resize-none"
                value={form.description}
                onChange={e => setForm({ ...form, description: e.target.value })}
              />
            </div>
          </div>
        </Card>

        <Card className="p-6">
          <h2 className="font-semibold text-gray-900 mb-4">Цена и характеристики</h2>
          <div className="grid grid-cols-2 gap-4">
            <Input label="Цена (₽)" placeholder="5000000" type="number" required {...f('price')} />
            <Input label="Площадь (м²)" placeholder="50" type="number" required {...f('area_sqm')} />
            <Input label="Комнат" placeholder="2" type="number" {...f('rooms')} />
            <Input label="Этаж" placeholder="5" type="number" {...f('floor')} />
            <Input label="Этажей в доме" placeholder="10" type="number" {...f('floors_total')} />
          </div>
        </Card>

        <Card className="p-6">
          <h2 className="font-semibold text-gray-900 mb-4">Адрес</h2>
          <div className="grid grid-cols-2 gap-4">
            <Input label="Страна" {...f('country')} />
            <Input label="Город" placeholder="Москва" required {...f('city')} />
            <Input label="Район" placeholder="Арбат" {...f('district')} />
            <Input label="Улица" placeholder="ул. Арбат" {...f('street')} />
            <Input label="Дом" placeholder="10" {...f('building')} />
          </div>
        </Card>

        {error && (
          <div className="bg-red-50 text-red-600 text-sm px-4 py-3 rounded-lg">{error}</div>
        )}

        <div className="flex gap-3">
          <Button type="submit" loading={loading} className="flex-1">
            Создать объявление
          </Button>
          <Button type="button" variant="secondary" onClick={() => navigate('/dashboard')}>
            Отмена
          </Button>
        </div>
      </form>
    </div>
  )
}