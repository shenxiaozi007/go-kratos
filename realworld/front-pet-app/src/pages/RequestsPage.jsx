import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { listFriendRequests, acceptFriendRequest, rejectFriendRequest } from '../api/friend'
import Header from '../components/Header'
import { useAuth } from '../context/AuthContext'

export default function RequestsPage() {
  const { user, loading: authLoading } = useAuth()
  const navigate = useNavigate()
  const [requests, setRequests] = useState([])
  const [loading, setLoading] = useState(true)
  const [acting, setActing] = useState(null)

  const load = () => {
    listFriendRequests()
      .then((res) => setRequests(res?.requests ?? []))
      .catch(() => setRequests([]))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    if (!user) {
      setLoading(false)
      return
    }
    setLoading(true)
    load()
  }, [user])

  const handleAccept = async (id) => {
    if (acting) return
    setActing(id)
    try {
      await acceptFriendRequest(id)
      setRequests((prev) => prev.filter((r) => r.id !== id))
    } catch (_) {}
    setActing(null)
  }

  const handleReject = async (id) => {
    if (acting) return
    setActing(id)
    try {
      await rejectFriendRequest(id)
      setRequests((prev) => prev.filter((r) => r.id !== id))
    } catch (_) {}
    setActing(null)
  }

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

  const pending = requests.filter((r) => r.status === 'PENDING')

  return (
    <div className="min-h-screen pb-20">
      <Header title="好友申请" showBack />
      <main className="p-4">
        {loading ? (
          <p className="text-center text-gray-500 py-8">加载中...</p>
        ) : pending.length === 0 ? (
          <p className="text-center text-gray-500 py-8">暂无待处理申请</p>
        ) : (
          <ul className="space-y-2">
            {pending.map((r) => (
              <li key={r.id} className="flex items-center gap-4 py-3 px-4 rounded-xl glass">
                <div className="w-12 h-12 rounded-full bg-gray-200 overflow-hidden shrink-0">
                  {r.from_user?.avatar_url ? (
                    <img
                      src={r.from_user.avatar_url}
                      alt=""
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <span className="material-symbols-outlined text-2xl text-gray-400 block text-center leading-12">
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
                <div className="flex gap-2">
                  <button
                    type="button"
                    onClick={() => handleAccept(r.id)}
                    disabled={!!acting}
                    className="py-2 px-3 rounded-lg bg-primary text-white text-sm disabled:opacity-50"
                  >
                    同意
                  </button>
                  <button
                    type="button"
                    onClick={() => handleReject(r.id)}
                    disabled={!!acting}
                    className="py-2 px-3 rounded-lg border border-gray-200 text-sm disabled:opacity-50"
                  >
                    拒绝
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </main>
    </div>
  )
}
