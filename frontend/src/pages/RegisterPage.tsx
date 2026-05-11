import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { authApi } from '../api/auth'
import { useAuthStore } from '../store/auth'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Card } from '../components/ui/Card'

const validateEmail = (email: string) =>
  /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)

const validatePhone = (phone: string) =>
  !phone || /^\+7\d{10}$/.test(phone)

const validatePassword = (password: string) =>
  password.length >= 8

export const RegisterPage = () => {
  const [form, setForm] = useState({
    email: '',
    password: '',
    confirmPassword: '',
    phone: '',
    role: 'buyer',
  })
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [serverError, setServerError] = useState('')
  const [loading, setLoading] = useState(false)
  const { setTokens } = useAuthStore()
  const navigate = useNavigate()

  const validate = () => {
    const e: Record<string, string> = {}

    if (!form.email) {
      e.email = 'Email обязателен'
    } else if (!validateEmail(form.email)) {
      e.email = 'Неверный формат email'
    }

    if (!form.password) {
      e.password = 'Пароль обязателен'
    } else if (!validatePassword(form.password)) {
      e.password = 'Пароль должен быть не менее 8 символов'
    }

    if (!form.confirmPassword) {
      e.confirmPassword = 'Подтвердите пароль'
    } else if (form.password !== form.confirmPassword) {
      e.confirmPassword = 'Пароли не совпадают'
    }

    if (!validatePhone(form.phone)) {
      e.phone = 'Формат: +79001234567'
    }

    setErrors(e)
    return Object.keys(e).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setServerError('')
    if (!validate()) return

    setLoading(true)
    try {
      const res = await authApi.register({
        email: form.email,
        password: form.password,
        phone: form.phone,
        role: form.role,
      })
      setTokens(res.access_token, res.refresh_token)
      navigate('/dashboard')
    } catch (err: any) {
      setServerError(err.message || 'Ошибка регистрации')
    } finally {
      setLoading(false)
    }
  }

  const set = (field: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setForm({ ...form, [field]: e.target.value })
    if (errors[field]) setErrors({ ...errors, [field]: '' })
  }

  return (
    <div className="min-h-[calc(100vh-200px)] flex items-center justify-center">
      <Card className="w-full max-w-md p-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">Регистрация</h1>
        <p className="text-gray-500 text-sm mb-6">Создайте аккаунт для размещения объявлений</p>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <Input
            label="Email"
            type="email"
            placeholder="ivan@example.ru"
            value={form.email}
            onChange={set('email')}
            error={errors.email}
            required
          />
          <Input
            label="Телефон"
            type="tel"
            placeholder="+79001234567"
            value={form.phone}
            onChange={set('phone')}
            error={errors.phone}
          />
          <Input
            label="Пароль"
            type="password"
            placeholder="Минимум 8 символов"
            value={form.password}
            onChange={set('password')}
            error={errors.password}
            required
          />
          <Input
            label="Подтвердите пароль"
            type="password"
            placeholder="Повторите пароль"
            value={form.confirmPassword}
            onChange={set('confirmPassword')}
            error={errors.confirmPassword}
            required
          />

          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium text-gray-700">Я являюсь</label>
            <select
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none focus:border-blue-500"
              value={form.role}
              onChange={set('role')}
            >
              <option value="buyer">Покупателем</option>
              <option value="seller">Продавцом</option>
              <option value="agent">Агентом</option>
            </select>
          </div>

          {serverError && (
            <div className="bg-red-50 text-red-600 text-sm px-4 py-3 rounded-lg">
              {serverError}
            </div>
          )}

          <Button type="submit" loading={loading} className="w-full mt-2">
            Зарегистрироваться
          </Button>
        </form>

        <p className="text-center text-sm text-gray-500 mt-6">
          Уже есть аккаунт?{' '}
          <Link to="/login" className="text-blue-600 hover:underline font-medium">Войти</Link>
        </p>
      </Card>
    </div>
  )
}