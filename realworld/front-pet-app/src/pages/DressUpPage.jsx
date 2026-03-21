import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { getMyPet } from '../api/pet'
import { getPetAppearance, updatePetAppearance } from '../api/pet'
import { getInventory } from '../api/shop'
import Header from '../components/Header'
import { useAuth } from '../context/AuthContext'

export default function DressUpPage() {
  const { user, loading: authLoading } = useAuth()
  const navigate = useNavigate()
  const [pet, setPet] = useState(null)
  const [appearance, setAppearance] = useState({ head: 0, body: 0, neck: 0 })
  const [inventory, setInventory] = useState([])
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (!user) {
      setLoading(false)
      return
    }
    getMyPet()
      .then((r) => r?.pet ?? null)
      .then((p) => {
        setPet(p)
        if (p) return getPetAppearance(p.id)
        return null
      })
      .then((res) => {
        if (res) {
          setAppearance({
            head: res.head_item_id ?? 0,
            body: res.body_item_id ?? 0,
            neck: res.neck_item_id ?? 0,
          })
        }
      })
      .catch(() => {})
    getInventory()
      .then((r) => (r?.items ?? []).filter((i) => i.category === '装扮'))
      .then(setInventory)
      .catch(() => setInventory([]))
      .finally(() => setLoading(false))
  }, [user])

  const handleSave = async () => {
    if (!pet || saving) return
    setSaving(true)
    try {
      await updatePetAppearance(pet.id, {
        head_item_id: appearance.head,
        body_item_id: appearance.body,
        neck_item_id: appearance.neck,
      })
      navigate(-1)
    } catch (e) {
      alert(e?.data?.message || e?.message || '保存失败')
    } finally {
      setSaving(false)
    }
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
        <Header title="换装" showBack />
        <p className="text-center text-gray-500 py-12">还没有宠物</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen pb-20">
      <Header title="换装" showBack />
      <main className="p-4">
        <div className="rounded-2xl bg-gray-100 aspect-square max-h-64 flex items-center justify-center mb-6">
          {pet.avatar_url ? (
            <img src={pet.avatar_url} alt="" className="max-h-full object-contain" />
          ) : (
            <span className="material-symbols-outlined text-8xl text-gray-300">pets</span>
          )}
        </div>
        <p className="text-sm text-gray-500 mb-4">从背包选择装扮（需先购买装扮类商品）</p>
        {inventory.length === 0 && (
          <p className="text-gray-500 text-sm">背包中暂无装扮，去商店购买吧</p>
        )}
        <button
          type="button"
          onClick={handleSave}
          disabled={saving}
          className="w-full mt-6 py-3 rounded-xl bg-primary text-white font-medium disabled:opacity-50"
        >
          {saving ? '保存中...' : '保存搭配'}
        </button>
      </main>
    </div>
  )
}
