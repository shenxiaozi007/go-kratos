import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { listAchievements } from '../api/achievement'
import Header from '../components/Header'
import { useAuth } from '../context/AuthContext'

export default function AchievementsPage() {
  const { user, loading: authLoading } = useAuth()
  const navigate = useNavigate()
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!user) {
      setLoading(false)
      return
    }
    listAchievements()
      .then((res) => setItems(res?.items ?? []))
      .catch(() => setItems([]))
      .finally(() => setLoading(false))
  }, [user])

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-500">加载中...</p>
      </div>
    )
  }

  if (!user) {
    navigate('/login')
    return null
  }

  return (
    <div className="min-h-screen pb-20">
      <Header title="成就" showBack />
      <main className="p-4">
        {loading ? (
          <p className="text-center text-gray-500 py-8">加载中...</p>
        ) : items.length === 0 ? (
          <p className="text-center text-gray-500 py-8">暂无成就</p>
        ) : (
          <ul className="space-y-2">
            {items.map(({ achievement, unlocked, unlocked_at }) => (
              <li
                key={achievement?.id}
                className={`flex items-center gap-4 py-3 px-4 rounded-xl ${
                  unlocked ? 'glass' : 'bg-gray-50 opacity-75'
                }`}
              >
                <div className="w-12 h-12 rounded-xl bg-gray-200 flex items-center justify-center shrink-0">
                  {achievement?.icon_url ? (
                    <img
                      src={achievement.icon_url}
                      alt=""
                      className="w-full h-full object-cover rounded-xl"
                    />
                  ) : (
                    <span className="material-symbols-outlined text-2xl text-gray-400">
                      {unlocked ? 'military_tech' : 'lock'}
                    </span>
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-medium">{achievement?.name}</p>
                  <p className="text-xs text-gray-500">{achievement?.description}</p>
                  {unlocked && unlocked_at > 0 && (
                    <p className="text-xs text-primary mt-0.5">
                      已解锁 · {new Date(unlocked_at * 1000).toLocaleDateString()}
                    </p>
                  )}
                </div>
              </li>
            ))}
          </ul>
        )}
      </main>
    </div>
  )
}
