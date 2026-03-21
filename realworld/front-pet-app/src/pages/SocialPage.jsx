import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { listFriends, getRanking, listFriendRequests } from '../api/friend'
import Header from '../components/Header'
import { useAuth } from '../context/AuthContext'

const TABS = ['好友', '排行榜', '申请']

export default function SocialPage() {
  const { user, loading: authLoading } = useAuth()
  const navigate = useNavigate()
  const [tab, setTab] = useState(0)
  const [friends, setFriends] = useState([])
  const [ranking, setRanking] = useState([])
  const [requests, setRequests] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!user) {
      setLoading(false)
      return
    }
    const load = () => {
      setLoading(true)
      Promise.all([
        listFriends().catch(() => ({ friends: [] })),
        getRanking('heart_points', 20).catch(() => ({ ranking: [] })),
        listFriendRequests().catch(() => ({ requests: [] })),
      ]).then(([f, r, req]) => {
        setFriends(f?.friends ?? [])
        setRanking(r?.ranking ?? [])
        setRequests(req?.requests ?? [])
        setLoading(false)
      })
    }
    load()
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

  const list =
    tab === 0 ? friends : tab === 1 ? ranking : requests
  const isRequests = tab === 2

  return (
    <div className="min-h-screen pb-20">
      <Header title="社交" />
      <div className="flex border-b border-gray-100">
        {TABS.map((label, i) => (
          <button
            key={label}
            type="button"
            onClick={() => setTab(i)}
            className={`flex-1 py-3 text-sm font-medium ${
              tab === i ? 'text-primary border-b-2 border-primary' : 'text-gray-500'
            }`}
          >
            {label}
          </button>
        ))}
      </div>
      <main className="p-4">
        {loading ? (
          <p className="text-center text-gray-500 py-8">加载中...</p>
        ) : list.length === 0 ? (
          <p className="text-center text-gray-500 py-8">暂无数据</p>
        ) : (
          <ul className="space-y-2">
            {tab === 1 &&
              ranking.map((u, i) => (
                <li
                  key={u.id}
                  className="flex items-center gap-3 py-3 px-4 rounded-xl glass"
                >
                  <span className="text-lg font-semibold text-gray-400 w-6">#{i + 1}</span>
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
                    <p className="text-xs text-gray-500">爱心 {u.heart_points ?? 0}</p>
                  </div>
                </li>
              ))}
            {tab === 0 &&
              friends.map((u) => (
                <li key={u.id} className="flex items-center gap-3 py-3 px-4 rounded-xl glass">
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
                    <p className="text-xs text-gray-500">Lv.{u.level ?? 1}</p>
                  </div>
                </li>
              ))}
            {isRequests &&
              requests.map((r) => (
                <li key={r.id} className="flex items-center gap-3 py-3 px-4 rounded-xl glass">
                  <div className="w-10 h-10 rounded-full bg-gray-200 overflow-hidden shrink-0">
                    {r.from_user?.avatar_url ? (
                      <img
                        src={r.from_user.avatar_url}
                        alt=""
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <span className="material-symbols-outlined text-2xl text-gray-400 block text-center leading-10">
                        person
                      </span>
                    )}
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-medium truncate">
                      {r.from_user?.nickname || r.from_user?.username || '用户'}
                    </p>
                    <p className="text-xs text-gray-500">{r.message || '请求加你为好友'}</p>
                  </div>
                  <button
                    type="button"
                    onClick={() => navigate(`/requests`)}
                    className="text-sm text-primary"
                  >
                    处理
                  </button>
                </li>
              ))}
          </ul>
        )}
      </main>
      <button
        type="button"
        onClick={() => navigate('/ranking')}
        className="fixed bottom-20 right-4 w-12 h-12 rounded-full bg-primary text-white shadow-lg flex items-center justify-center"
      >
        <span className="material-symbols-outlined">leaderboard</span>
      </button>
    </div>
  )
}
