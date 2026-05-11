import { FC, ReactNode } from 'react'

interface Props {
  children: ReactNode
  className?: string
  onClick?: () => void
}

export const Card: FC<Props> = ({ children, className = '', onClick }) => (
  <div
    className={`bg-white rounded-xl border border-gray-200 shadow-sm
      ${onClick ? 'cursor-pointer hover:shadow-md transition-shadow' : ''}
      ${className}`}
    onClick={onClick}
  >
    {children}
  </div>
)