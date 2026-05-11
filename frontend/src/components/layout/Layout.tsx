import { FC, ReactNode } from 'react'
import { Header } from './Header'

interface Props {
  children: ReactNode
}

export const Layout: FC<Props> = ({ children }) => (
  <div className="min-h-screen bg-gray-50">
    <Header />
    <main className="max-w-7xl mx-auto px-4 py-8">
      {children}
    </main>
    <footer className="bg-white border-t border-gray-200 mt-16">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="flex flex-col md:flex-row justify-between gap-4">
          <div>
            <div className="font-semibold text-gray-900">НедвижимостьРФ</div>
            <div className="text-sm text-gray-500 mt-1">Покупка и продажа недвижимости</div>
          </div>
          <div className="text-sm text-gray-400">© 2026 НедвижимостьРФ. Все права защищены.</div>
        </div>
      </div>
    </footer>
  </div>
)