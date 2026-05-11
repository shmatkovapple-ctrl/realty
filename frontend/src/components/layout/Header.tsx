import { FC } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../../store/auth'
import { Button } from '../ui/Button'

export const Header: FC = () => {
  const { isAuthenticated, user, logout } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/')
  }

  return (
    <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2">
          <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
            <span className="text-white font-bold text-sm">Н</span>
          </div>
          <span className="font-semibold text-gray-900 text-lg">НедвижимостьРФ</span>
        </Link>

        <nav className="hidden md:flex items-center gap-6">
          <Link to="/listings" className="text-sm text-gray-600 hover:text-gray-900 transition-colors">
            Каталог
          </Link>
          {isAuthenticated && (
            <Link to="/favorites" className="text-sm text-gray-600 hover:text-gray-900 transition-colors">
              Избранное
            </Link>
          )}
        </nav>

        <div className="flex items-center gap-3">
          {isAuthenticated ? (
  <>
    <Link to="/dashboard">
      <Button variant="ghost" size="sm">
        {user?.first_name || 'Кабинет'}
      </Button>
    </Link>
    <Button variant="secondary" size="sm" onClick={handleLogout}>
      Выйти
    </Button>
  </>
) : (
  <>
    <Link to="/login">
      <Button variant="ghost" size="sm">Войти</Button>
    </Link>
    <Link to="/register">
      <Button size="sm">Зарегистрироваться</Button>
    </Link>
  </>
)}
        </div>
      </div>
    </header>
  )
}