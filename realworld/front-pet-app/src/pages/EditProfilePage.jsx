import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { getProfile, updateProfile } from '../api/user'
import Header from '../components/Header'
import { useAuth } from '../context/AuthContext'

export default function EditProfilePage() {
  const { user, loading: authLoading, loadUser } = useAuth()
  const navigate = useNavigate()
  const [nickname, setNickname] = useState('')
  const [signature, setSignature] = useState('')
  const [avatarUrl, setAvatarUrl] = useState('')
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (!user) return
    getProfile()
      .then((res) => {
        const u = res?.user ?? res
        if (u) {
          setNickname(u.nickname || u.username || '')
          setSignature(u.signature || '')
          setAvatarUrl(u.avatar_url || '')
        }
      })
      .catch(() => {})
  }, [user])

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (saving) return
    setSaving(true)
    try {
      await updateProfile({ nickname: nickname.trim(), signature: signature.trim(), avatar_url: avatarUrl.trim() })
      loadUser()
      navigate('/profile')
    } catch (e) {
      alert(e?.data?.message || e?.message || '保存失败')
    } finally {
      setSaving(false)
    }
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

  return (
    <div className="min-h-screen pb-20">
      <Header title="编辑资料" showBack />
      <main className="p-4">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">昵称</label>
            <input
              type="text"
              value={nickname}
              onChange={(e) => setNickname(e.target.value)}
              className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary"
              placeholder="昵称"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">个性签名</label>
            <textarea
              value={signature}
              onChange={(e) => setSignature(e.target.value)}
              className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary resize-none"
              rows={3}
              placeholder="写点什么吧"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">头像 URL</label>
            <input
              type="url"
              value={avatarUrl}
              onChange={(e) => setAvatarUrl(e.target.value)}
              className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary"
              placeholder="https://..."
            />
          </div>
          <button
            type="submit"
            disabled={saving}
            className="w-full py-3 rounded-xl bg-primary text-white font-medium disabled:opacity-50"
          >
            {saving ? '保存中...' : '保存'}
          </button>
        </form>
      </main>
    </div>
  )
}
