import { FC, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { authApi } from '../api/auth'
import { useAuthStore } from '../store/auth'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Card } from '../components/ui/Card'

export const LoginPage: FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { setTokens, setUser } = useAuthStore()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const res = await authApi.login({ email, password })
      setTokens(res.access_token, res.refresh_token)
      if (res.profile) setUser(res.profile)
      navigate('/dashboard')
    } catch (err: any) {
      setError(err.message || 'Ошибка входа')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-[calc(100vh-200px)] flex items-center justify-center">
      <Card className="w-full max-w-md p-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">Вход в аккаунт</h1>
        <p className="text-gray-500 text-sm mb-6">Введите ваши данные для входа</p>
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <Input label="Email" type="email" placeholder="ivan@example.ru"
            value={email} onChange={e => setEmail(e.target.value)} required />
          <Input label="Пароль" type="password" placeholder="••••••••"
            value={password} onChange={e => setPassword(e.target.value)} required />
          {error && (
            <div className="bg-red-50 text-red-600 text-sm px-4 py-3 rounded-lg">{error}</div>
          )}
          <Button type="submit" loading={loading} className="w-full mt-2">
            Войти
          </Button>
        </form>
        <p className="text-center text-sm text-gray-500 mt-6">
          Нет аккаунта?{' '}
          <Link to="/register" className="text-blue-600 hover:underline font-medium">
            Зарегистрироваться
          </Link>
        </p>
      </Card>
    </div>
  )
}