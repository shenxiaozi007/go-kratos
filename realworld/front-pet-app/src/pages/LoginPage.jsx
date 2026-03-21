import { useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { login, register } from '../api/auth'
import Header from '../components/Header'
import { useAuth } from '../context/AuthContext'

export default function LoginPage() {
  const [search] = useSearchParams()
  const isRegister = search.get('register') === '1'
  const navigate = useNavigate()
  const { loginSuccess } = useAuth()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    if (!username.trim() || !password) {
      setError('请输入用户名和密码')
      return
    }
    setLoading(true)
    try {
      if (isRegister) {
        await register(username.trim(), password)
      }
      const res = await login(username.trim(), password)
      const token = res?.token ?? res?.access_token
      if (token) {
        loginSuccess(token, res?.user ?? null)
        navigate(search.get('from') || '/')
      } else {
        setError('登录失败，请重试')
      }
    } catch (err) {
      setError(err?.data?.message || err?.message || '登录失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex flex-col">
      <Header title={isRegister ? '注册' : '登录'} showBack />
      <main className="flex-1 p-6 flex flex-col justify-center">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">用户名</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary focus:border-transparent"
              placeholder="请输入用户名"
              autoComplete="username"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">密码</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary focus:border-transparent"
              placeholder="请输入密码"
              autoComplete={isRegister ? 'new-password' : 'current-password'}
            />
          </div>
          {error && <p className="text-sm text-red-500">{error}</p>}
          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 rounded-xl bg-primary text-white font-medium disabled:opacity-50"
          >
            {loading ? '提交中...' : isRegister ? '注册并登录' : '登录'}
          </button>
        </form>
        <p className="mt-4 text-center text-sm text-gray-500">
          {isRegister ? (
            <>
              已有账号？{' '}
              <button type="button" onClick={() => navigate('/login')} className="text-primary">
                去登录
              </button>
            </>
          ) : (
            <>
              没有账号？{' '}
              <button
                type="button"
                onClick={() => navigate('/login?register=1')}
                className="text-primary"
              >
                去注册
              </button>
            </>
          )}
        </p>
      </main>
    </div>
  )
}
