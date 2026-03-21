import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { getMyPet } from '../api/pet'
import Header from '../components/Header'
import SignInModal from '../components/SignInModal'
import { useAuth } from '../context/AuthContext'

export default function HomePage() {
  const { user, loading: authLoading } = useAuth()
  const navigate = useNavigate()
  const [pet, setPet] = useState(null)
  const [petLoading, setPetLoading] = useState(true)
  const [signModalOpen, setSignModalOpen] = useState(false)

  useEffect(() => {
    if (!user) {
      setPet(null)
      setPetLoading(false)
      return
    }
    getMyPet()
      .then((res) => setPet(res?.pet ?? null))
      .catch(() => setPet(null))
      .finally(() => setPetLoading(false))
  }, [user])

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-500">加载中...</p>
      </div>
    )
  }

  if (!user) {
    return (
      <div className="min-h-screen flex flex-col">
        <Header title="萌宠之家" />
        <div className="flex-1 flex flex-col items-center justify-center p-6 text-center">
          <p className="text-gray-600 mb-4">登录后查看你的宠物</p>
          <button
            type="button"
            onClick={() => navigate('/login')}
            className="px-6 py-3 rounded-xl bg-primary text-white font-medium"
          >
            去登录
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen pb-20">
      <Header title="萌宠之家" />
      <main className="p-4">
        {petLoading ? (
          <div className="py-12 text-center text-gray-500">加载中...</div>
        ) : !pet ? (
          <div className="glass rounded-2xl p-6 text-center">
            <p className="text-gray-600 mb-4">还没有宠物，快去领养一只吧～</p>
            <button
              type="button"
              onClick={() => navigate('/adoption')}
              className="px-6 py-3 rounded-xl bg-primary text-white font-medium"
            >
              去领养
            </button>
          </div>
        ) : (
          <>
            <div className="glass rounded-2xl p-5 mb-4">
              <div className="flex items-center gap-4">
                <div className="w-20 h-20 rounded-2xl bg-gray-200 flex items-center justify-center overflow-hidden">
                  {pet.avatar_url ? (
                    <img src={pet.avatar_url} alt="" className="w-full h-full object-cover" />
                  ) : (
                    <span className="material-symbols-outlined text-4xl text-gray-400">pets</span>
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <h2 className="text-xl font-semibold truncate">{pet.name || '未命名'}</h2>
                  <p className="text-sm text-gray-500">{pet.mood || '心情不错'}</p>
                  <p className="text-xs text-primary mt-1">亲密度 {pet.affection ?? 0}/1000</p>
                </div>
              </div>
              <div className="grid grid-cols-3 gap-2 mt-4 text-center text-sm">
                <div className="rounded-xl bg-gray-50 py-2">
                  <p className="text-gray-500">饱食度</p>
                  <p className="font-medium">{pet.fullness ?? 0}%</p>
                </div>
                <div className="rounded-xl bg-gray-50 py-2">
                  <p className="text-gray-500">心情值</p>
                  <p className="font-medium">{pet.happiness ?? 0}%</p>
                </div>
                <div className="rounded-xl bg-gray-50 py-2">
                  <p className="text-gray-500">清洁度</p>
                  <p className="font-medium">{pet.cleanliness ?? 0}%</p>
                </div>
              </div>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <button
                type="button"
                onClick={() => navigate('/inventory')}
                className="flex items-center justify-center gap-2 py-4 rounded-xl glass"
              >
                <span className="material-symbols-outlined text-primary">restaurant</span>
                <span>喂食</span>
              </button>
              <button
                type="button"
                onClick={() => setSignModalOpen(true)}
                className="flex items-center justify-center gap-2 py-4 rounded-xl glass"
              >
                <span className="material-symbols-outlined text-primary">calendar_today</span>
                <span>签到</span>
              </button>
              <button
                type="button"
                onClick={() => navigate('/shop')}
                className="flex items-center justify-center gap-2 py-4 rounded-xl glass"
              >
                <span className="material-symbols-outlined text-primary">storefront</span>
                <span>商店</span>
              </button>
              <button
                type="button"
                onClick={() => navigate('/interaction')}
                className="flex items-center justify-center gap-2 py-4 rounded-xl glass"
              >
                <span className="material-symbols-outlined text-primary">favorite</span>
                <span>互动</span>
              </button>
            </div>
            <button
              type="button"
              onClick={() => navigate('/adoption')}
              className="w-full mt-4 py-3 rounded-xl border-2 border-dashed border-gray-300 text-gray-500"
            >
              + 领养新宠物
            </button>
          </>
        )}
      </main>
      <SignInModal open={signModalOpen} onClose={() => setSignModalOpen(false)} />
    </div>
  )
}
