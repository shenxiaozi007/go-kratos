import { createContext, useContext, useState, useEffect } from 'react'
import { getStoredToken, setToken as saveToken, getProfile } from '../api/auth'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)

  const loadUser = async () => {
    if (!getStoredToken()) {
      setUser(null)
      setLoading(false)
      return
    }
    try {
      const data = await getProfile()
      setUser(data?.user ?? data ?? null)
    } catch (_) {
      setUser(null)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadUser()
  }, [])

  const loginSuccess = (token, userData) => {
    saveToken(token)
    setUser(userData ?? null)
    if (!userData) loadUser()
  }

  const logout = () => {
    saveToken(null)
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, loading, loadUser, loginSuccess, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
