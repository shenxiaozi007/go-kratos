import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { getBreed } from '../api/breed'
import { createPet } from '../api/pet'
import Header from '../components/Header'
import { useAuth } from '../context/AuthContext'

const SPECIES_MAP = { CAT: '猫', DOG: '狗' }

export default function BreedDetailPage() {
  const { id } = useParams()
  const { user, loading: authLoading } = useAuth()
  const navigate = useNavigate()
  const [breed, setBreed] = useState(null)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [name, setName] = useState('')

  useEffect(() => {
    if (!user || !id) {
      setLoading(false)
      return
    }
    getBreed(id)
      .then((res) => setBreed(res?.breed ?? null))
      .catch(() => setBreed(null))
      .finally(() => setLoading(false))
  }, [user, id])

  const handleAdopt = async (e) => {
    e.preventDefault()
    if (!breed || !name.trim() || submitting) return
    setSubmitting(true)
    try {
      await createPet({
        name: name.trim(),
        species: breed.species === 'DOG' ? 2 : 1,
        breed_id: breed.id,
      })
      navigate('/')
    } catch (e) {
      alert(e?.data?.message || e?.message || '领养失败')
    } finally {
      setSubmitting(false)
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

  if (!breed) {
    return (
      <div className="min-h-screen flex flex-col">
        <Header title="品种详情" showBack />
        <p className="text-center text-gray-500 py-12">品种不存在</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen pb-20">
      <Header title="选择此品种领养" showBack />
      <main className="p-4">
        <div className="glass rounded-2xl p-5 mb-6">
          <div className="aspect-video rounded-xl bg-gray-100 mb-4 flex items-center justify-center overflow-hidden">
            {breed.image_url ? (
              <img src={breed.image_url} alt="" className="w-full h-full object-cover" />
            ) : (
              <span className="material-symbols-outlined text-6xl text-gray-300">pets</span>
            )}
          </div>
          <h2 className="text-xl font-semibold">{breed.name_cn || breed.name_en}</h2>
          <p className="text-sm text-gray-500 mt-1">
            {SPECIES_MAP[breed.species] || breed.species} · {breed.size_tag} · {breed.life_span}
          </p>
        </div>
        <form onSubmit={handleAdopt} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">给 ta 起个名字</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary"
              placeholder="例如：小橘"
              maxLength={20}
            />
          </div>
          <button
            type="submit"
            disabled={!name.trim() || submitting}
            className="w-full py-3 rounded-xl bg-primary text-white font-medium disabled:opacity-50"
          >
            {submitting ? '领养中...' : '确认领养'}
          </button>
        </form>
      </main>
    </div>
  )
}
