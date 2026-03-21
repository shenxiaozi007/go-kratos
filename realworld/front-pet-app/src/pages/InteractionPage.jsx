import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { getMyPet, interact } from '../api/pet'
import Header from '../components/Header'
import { useAuth } from '../context/AuthContext'

const ACTIONS = [
  { key: 'STROKE', label: '抚摸', icon: 'touch_app' },
  { key: 'TEASE', label: '逗猫棒', icon: 'sports_esports' },
  { key: 'TREAT', label: '零食', icon: 'cookie' },
  { key: 'BATH', label: '洗澡', icon: 'bathtub' },
]

export default function InteractionPage() {
  const { user, loading: authLoading } = useAuth()
  const navigate = useNavigate()
  const [pet, setPet] = useState(null)
  const [loading, setLoading] = useState(true)
  const [acting, setActing] = useState(null)

  const loadPet = () => {
    getMyPet()
      .then((res) => setPet(res?.pet ?? null))
      .catch(() => setPet(null))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    if (!user) {
      setPet(null)
      setLoading(false)
      return
    }
    loadPet()
  }, [user])

  const handleAction = async (action) => {
    if (!pet || acting) return
    setActing(action)
    try {
      const res = await interact(action)
      if (res?.pet) setPet(res.pet)
    } catch (_) {}
    setActing(null)
  }

  if (authLoading || loading) {
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

  if (!pet) {
    return (
      <div className="min-h-screen flex flex-col">
        <Header title="互动" showBack />
        <div className="flex-1 flex flex-col items-center justify-center p-6">
          <p className="text-gray-600 mb-4">还没有宠物，先去领养吧</p>
          <button
            type="button"
            onClick={() => navigate('/adoption')}
            className="px-6 py-3 rounded-xl bg-primary text-white font-medium"
          >
            去领养
          </button>
        </div>
      </div>
    )
  }

  const affection = pet.affection ?? 0
  const percent = Math.min(100, Math.round((affection / 1000) * 100))

  return (
    <div className="min-h-screen pb-20 flex flex-col">
      <Header title="互动" showBack />
      <main className="flex-1 flex flex-col p-4">
        <div className="glass rounded-2xl p-4 text-center mb-4">
          <h2 className="text-lg font-semibold">{pet.name || '未命名'}</h2>
          <p className="text-sm text-gray-500 mt-1">{pet.mood || '心情不错'}</p>
          <div className="mt-3 h-2 bg-gray-200 rounded-full overflow-hidden">
            <div
              className="h-full bg-primary rounded-full transition-all duration-300"
              style={{ width: `${percent}%` }}
            />
          </div>
          <p className="text-xs text-gray-500 mt-1">亲密度 {affection}/1000</p>
        </div>
        <div className="flex-1 flex items-center justify-center rounded-2xl bg-gray-100/80 min-h-[200px]">
          {pet.avatar_url ? (
            <img src={pet.avatar_url} alt="" className="max-h-64 object-contain" />
          ) : (
            <span className="material-symbols-outlined text-8xl text-gray-300">pets</span>
          )}
        </div>
        <div className="grid grid-cols-4 gap-2 mt-4">
          {ACTIONS.map(({ key, label, icon }) => (
            <button
              key={key}
              type="button"
              onClick={() => handleAction(key)}
              disabled={!!acting}
              className="flex flex-col items-center gap-1 py-4 rounded-xl glass hover:bg-primary/10 disabled:opacity-50"
            >
              <span className="material-symbols-outlined text-3xl text-primary">{icon}</span>
              <span className="text-sm">{label}</span>
            </button>
          ))}
        </div>
      </main>
    </div>
  )
}
