import { useNavigate } from 'react-router-dom'

export default function Header({ title, showBack, rightIcon, onRightClick }) {
  const navigate = useNavigate()

  return (
    <header className="sticky top-0 z-40 flex items-center justify-between h-14 px-4 bg-white/80 backdrop-blur-glass border-b border-gray-100">
      <div className="flex items-center gap-2 min-w-0">
        {showBack && (
          <button
            type="button"
            onClick={() => navigate(-1)}
            className="p-1 -ml-1 rounded-lg hover:bg-gray-100"
            aria-label="返回"
          >
            <span className="material-symbols-outlined text-2xl text-gray-700">arrow_back</span>
          </button>
        )}
        <h1 className="text-lg font-semibold truncate">{title || '萌宠之家'}</h1>
      </div>
      {rightIcon && (
        <button
          type="button"
          onClick={onRightClick}
          className="p-1 rounded-lg hover:bg-gray-100"
          aria-label="操作"
        >
          <span className="material-symbols-outlined text-2xl text-gray-700">{rightIcon}</span>
        </button>
      )}
    </header>
  )
}
