import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { getProfile } from '../api/user'
import { getMyPet, listPets } from '../api/pet'
import Header from '../components/Header'
import SignInModal from '../components/SignInModal'
import { useAuth } from '../context/AuthContext'
import { logout } from '../api/auth'

export default function ProfilePage() {
  const { user, loading: authLoading, loadUser } = useAuth()
  const navigate = useNavigate()
  const [profile, setProfile] = useState(null)
  const [pets, setPets] = useState([])
  const [signModalOpen, setSignModalOpen] = useState(false)

  useEffect(() => {
    if (!user) {
      setProfile(null)
      setPets([])
      return
    }
    Promise.all([getProfile().catch(() => null), listPets().catch(() => ({ pets: [] }))]).then(
      ([p, petRes]) => {
        setProfile(p?.user ?? p ?? null)
        setPets(petRes?.pets ?? [])
      }
    )
  }, [user])

  const handleLogout = () => {
    logout()
    loadUser()
    navigate('/')
  }

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
        <Header title="我的" />
        <div className="flex-1 flex flex-col items-center justify-center p-6">
          <p className="text-gray-600 mb-4">登录后查看个人中心</p>
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

  const u = profile ?? user
  const nickname = u.nickname || u.username || '用户'
  const avatar = u.avatar_url || ''
  const level = u.level ?? 1
  const coins = u.coins ?? 0
  const heart = u.heart_points ?? 0

  return (
    <div className="min-h-screen pb-20">
      <Header title="我的" />
      <main className="p-4">
        <div className="glass rounded-2xl p-5 flex items-center gap-4 mb-4">
          <div className="w-16 h-16 rounded-2xl bg-gray-200 flex items-center justify-center overflow-hidden shrink-0">
            {avatar ? (
              <img src={avatar} alt="" className="w-full h-full object-cover" />
            ) : (
              <span className="material-symbols-outlined text-3xl text-gray-400">person</span>
            )}
          </div>
          <div className="flex-1 min-w-0">
            <h2 className="text-lg font-semibold truncate">{nickname}</h2>
            <p className="text-sm text-gray-500">Lv.{level}</p>
            <p className="text-xs text-primary mt-0.5">
              {coins} 金币 · {heart} 爱心
            </p>
          </div>
        </div>
        <div className="space-y-2">
          <button
            type="button"
            onClick={() => navigate('/edit-profile')}
            className="w-full flex items-center gap-3 py-3 px-4 rounded-xl glass"
          >
            <span className="material-symbols-outlined text-gray-600">edit</span>
            <span>编辑资料</span>
          </button>
          <button
            type="button"
            onClick={() => setSignModalOpen(true)}
            className="w-full flex items-center gap-3 py-3 px-4 rounded-xl glass"
          >
            <span className="material-symbols-outlined text-gray-600">calendar_today</span>
            <span>每日签到</span>
          </button>
          <button
            type="button"
            onClick={() => navigate('/achievements')}
            className="w-full flex items-center gap-3 py-3 px-4 rounded-xl glass"
          >
            <span className="material-symbols-outlined text-gray-600">military_tech</span>
            <span>成就</span>
          </button>
        </div>
        <div className="mt-6">
          <h3 className="text-sm font-medium text-gray-500 mb-2">我的宠物</h3>
          <div className="space-y-2">
            {pets.length === 0 ? (
              <button
                type="button"
                onClick={() => navigate('/adoption')}
                className="w-full py-4 rounded-xl border-2 border-dashed border-gray-300 text-gray-500 flex items-center justify-center gap-2"
              >
                <span className="material-symbols-outlined">add</span>
                添加宠物
              </button>
            ) : (
              pets.map((p) => (
                <button
                  key={p.id}
                  type="button"
                  onClick={() => navigate('/interaction')}
                  className="w-full flex items-center gap-3 py-3 px-4 rounded-xl glass text-left"
                >
                  <div className="w-12 h-12 rounded-xl bg-gray-200 overflow-hidden shrink-0">
                    {p.avatar_url ? (
                      <img src={p.avatar_url} alt="" className="w-full h-full object-cover" />
                    ) : (
                      <span className="material-symbols-outlined text-2xl text-gray-400 block text-center leading-12">
                        pets
                      </span>
                    )}
                  </div>
                  <span className="font-medium">{p.name || '未命名'}</span>
                </button>
              ))
            )}
          </div>
        </div>
        <button
          type="button"
          onClick={handleLogout}
          className="w-full mt-8 py-3 rounded-xl border border-gray-200 text-gray-600"
        >
          退出登录
        </button>
      </main>
      <SignInModal open={signModalOpen} onClose={() => setSignModalOpen(false)} />
    </div>
  )
}
