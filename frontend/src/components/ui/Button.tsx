import { ButtonHTMLAttributes, FC } from 'react'

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
}

const variants = {
  primary:   'bg-blue-600 hover:bg-blue-700 text-white',
  secondary: 'bg-white hover:bg-gray-50 text-gray-700 border border-gray-300',
  danger:    'bg-red-600 hover:bg-red-700 text-white',
  ghost:     'hover:bg-gray-100 text-gray-700',
}

const sizes = {
  sm: 'px-3 py-1.5 text-sm',
  md: 'px-4 py-2 text-sm',
  lg: 'px-6 py-3 text-base',
}

export const Button: FC<Props> = ({
  variant = 'primary',
  size = 'md',
  loading,
  children,
  className = '',
  disabled,
  ...props
}) => (
  <button
    className={`inline-flex items-center justify-center rounded-lg font-medium transition-colors
      disabled:opacity-50 disabled:cursor-not-allowed
      ${variants[variant]} ${sizes[size]} ${className}`}
    disabled={disabled || loading}
    {...props}
  >
    {loading ? (
      <span className="flex items-center gap-2">
        <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
        </svg>
        Загрузка...
      </span>
    ) : children}
  </button>
)