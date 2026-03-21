import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { getInventory, useItem } from '../api/shop'
import { getMyPet } from '../api/pet'
import Header from '../components/Header'
import { useAuth } from '../context/AuthContext'

export default function InventoryPage() {
  const { user, loading: authLoading } = useAuth()
  const navigate = useNavigate()
  const [items, setItems] = useState([])
  const [pet, setPet] = useState(null)
  const [loading, setLoading] = useState(true)
  const [using, setUsing] = useState(null)

  useEffect(() => {
    if (!user) {
      setLoading(false)
      return
    }
    Promise.all([
      getInventory().then((r) => r?.items ?? []).catch(() => []),
      getMyPet().then((r) => r?.pet ?? null).catch(() => null),
    ]).then(([inv, p]) => {
      setItems(inv)
      setPet(p)
      setLoading(false)
    })
  }, [user])

  const handleUse = async (item) => {
    if (using) return
    setUsing(item.item_id)
    try {
      const res = await useItem(item.item_id, pet?.id ?? 0)
      if (res?.success) {
        setItems((prev) =>
          prev.map((i) =>
            i.item_id === item.item_id ? { ...i, quantity: (i.quantity ?? 1) - 1 } : i
          ).filter((i) => (i.quantity ?? 0) > 0)
        )
      } else {
        alert(res?.message || '使用失败')
      }
    } catch (e) {
      alert(e?.data?.message || e?.message || '使用失败')
    } finally {
      setUsing(null)
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
      <Header title="背包" showBack />
      <main className="p-4">
        {loading ? (
          <p className="text-center text-gray-500 py-8">加载中...</p>
        ) : items.length === 0 ? (
          <p className="text-center text-gray-500 py-8">背包空空如也，去商店逛逛吧</p>
        ) : (
          <ul className="space-y-2">
            {items.map((item) => (
              <li
                key={item.item_id}
                className="flex items-center gap-4 py-3 px-4 rounded-xl glass"
              >
                <div className="w-14 h-14 rounded-xl bg-gray-100 flex items-center justify-center overflow-hidden shrink-0">
                  {item.image_url ? (
                    <img src={item.image_url} alt="" className="w-full h-full object-cover" />
                  ) : (
                    <span className="material-symbols-outlined text-3xl text-gray-400">
                      inventory_2
                    </span>
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-medium">{item.name}</p>
                  <p className="text-xs text-gray-500">
                    {item.category} · x{item.quantity ?? 0}
                  </p>
                </div>
                <button
                  type="button"
                  onClick={() => handleUse(item)}
                  disabled={!!using}
                  className="py-2 px-4 rounded-xl bg-primary text-white text-sm disabled:opacity-50"
                >
                  {using === item.item_id ? '使用中...' : '使用'}
                </button>
              </li>
            ))}
          </ul>
        )}
      </main>
    </div>
  )
}
