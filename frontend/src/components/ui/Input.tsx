import { InputHTMLAttributes, FC } from 'react'

interface Props extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
}

export const Input: FC<Props> = ({ label, error, className = '', ...props }) => (
  <div className="flex flex-col gap-1">
    {label && <label className="text-sm font-medium text-gray-700">{label}</label>}
    <input
      className={`w-full px-3 py-2 border rounded-lg text-sm outline-none transition-colors
        border-gray-300 focus:border-blue-500 focus:ring-1 focus:ring-blue-500
        disabled:bg-gray-50 disabled:text-gray-500
        ${error ? 'border-red-500' : ''} ${className}`}
      {...props}
    />
    {error && <span className="text-xs text-red-500">{error}</span>}
  </div>
)