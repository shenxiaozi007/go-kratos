import { NavLink } from 'react-router-dom'

const tabs = [
  { to: '/', label: '首页', icon: 'home' },
  { to: '/social', label: '社交', icon: 'group' },
  { to: '/interaction', label: '互动', icon: 'pets' },
  { to: '/profile', label: '我的', icon: 'person' },
]

export default function BottomNav() {
  return (
    <nav className="fixed bottom-0 left-0 right-0 max-w-app mx-auto flex justify-around items-center h-16 bg-white/90 backdrop-blur-glass border-t border-gray-200 z-50">
      {tabs.map(({ to, label, icon }) => (
        <NavLink
          key={to}
          to={to}
          className={({ isActive }) =>
            `flex flex-col items-center justify-center gap-0.5 py-2 px-4 rounded-xl transition-colors ${
              isActive ? 'text-primary font-medium' : 'text-gray-500'
            }`
          }
        >
          <span className="material-symbols-outlined text-2xl">{icon}</span>
          <span className="text-xs">{label}</span>
        </NavLink>
      ))}
    </nav>
  )
}
