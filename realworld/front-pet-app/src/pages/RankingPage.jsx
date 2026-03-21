import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { getRanking } from '../api/friend'
import Header from '../components/Header'
import { useAuth } from '../context/AuthContext'

export default function RankingPage() {
  const { user, loading: authLoading } = useAuth()
  const navigate = useNavigate()
  const [type, setType] = useState('heart_points')
  const [list, setList] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!user) {
      setLoading(false)
      return
    }
    setLoading(true)
    getRanking(type, 50)
      .then((res) => setList(res?.ranking ?? []))
      .catch(() => setList([]))
      .finally(() => setLoading(false))
  }, [user, type])

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
      <Header title="排行榜" showBack />
      <div className="flex gap-2 p-4 border-b border-gray-100">
        <button
          type="button"
          onClick={() => setType('heart_points')}
          className={`flex-1 py-2 rounded-xl text-sm font-medium ${
            type === 'heart_points' ? 'bg-primary text-white' : 'glass'
          }`}
        >
          爱心积分
        </button>
        <button
          type="button"
          onClick={() => setType('level')}
          className={`flex-1 py-2 rounded-xl text-sm font-medium ${
            type === 'level' ? 'bg-primary text-white' : 'glass'
          }`}
        >
          等级
        </button>
      </div>
      <main className="p-4">
        {loading ? (
          <p className="text-center text-gray-500 py-8">加载中...</p>
        ) : (
          <ul className="space-y-2">
            {list.map((u, i) => (
              <li key={u.id} className="flex items-center gap-4 py-3 px-4 rounded-xl glass">
                <span className="text-lg font-bold text-gray-400 w-8">#{i + 1}</span>
                <div className="w-10 h-10 rounded-full bg-gray-200 overflow-hidden shrink-0">
                  {u.avatar_url ? (
                    <img src={u.avatar_url} alt="" className="w-full h-full object-cover" />
                  ) : (
                    <span className="material-symbols-outlined text-2xl text-gray-400 block text-center leading-10">
                      person
                    </span>
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate">{u.nickname || u.username}</p>
                  <p className="text-xs text-gray-500">
                    Lv.{u.level ?? 1} · 爱心 {u.heart_points ?? 0}
                  </p>
                </div>
              </li>
            ))}
          </ul>
        )}
      </main>
    </div>
  )
}
