import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { listShopItems, buy } from '../api/shop'
import Header from '../components/Header'
import { useAuth } from '../context/AuthContext'

export default function ShopPage() {
  const { user, loading: authLoading } = useAuth()
  const navigate = useNavigate()
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [buying, setBuying] = useState(null)

  useEffect(() => {
    if (!user) {
      setLoading(false)
      return
    }
    listShopItems()
      .then((res) => setItems(res?.items ?? []))
      .catch(() => setItems([]))
      .finally(() => setLoading(false))
  }, [user])

  const handleBuy = async (item) => {
    if (buying) return
    setBuying(item.id)
    try {
      const res = await buy(item.id, 1)
      if (res?.success) {
        navigate('/purchase-success', { state: { item, quantity: 1 } })
      } else {
        alert(res?.message || '购买失败')
      }
    } catch (e) {
      alert(e?.data?.message || e?.message || '购买失败')
    } finally {
      setBuying(null)
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
      <Header title="商店" showBack />
      <main className="p-4">
        {loading ? (
          <p className="text-center text-gray-500 py-8">加载中...</p>
        ) : items.length === 0 ? (
          <p className="text-center text-gray-500 py-8">暂无商品</p>
        ) : (
          <div className="grid grid-cols-2 gap-4">
            {items.map((item) => (
              <div key={item.id} className="glass rounded-2xl p-4 flex flex-col">
                <div className="aspect-square rounded-xl bg-gray-100 mb-3 flex items-center justify-center overflow-hidden">
                  {item.image_url ? (
                    <img
                      src={item.image_url}
                      alt=""
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <span className="material-symbols-outlined text-4xl text-gray-300">
                      inventory_2
                    </span>
                  )}
                </div>
                <h3 className="font-medium truncate">{item.name}</h3>
                <p className="text-xs text-gray-500 mt-0.5">{item.category}</p>
                <p className="text-primary font-medium mt-1">{item.price_coins ?? 0} 金币</p>
                <button
                  type="button"
                  onClick={() => handleBuy(item)}
                  disabled={!!buying}
                  className="mt-2 py-2 rounded-xl bg-primary text-white text-sm font-medium disabled:opacity-50"
                >
                  {buying === item.id ? '购买中...' : '购买'}
                </button>
              </div>
            ))}
          </div>
        )}
      </main>
    </div>
  )
}
