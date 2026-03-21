import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { listBreeds } from '../api/breed'
import Header from '../components/Header'
import { useAuth } from '../context/AuthContext'

export default function AdoptionPage() {
  const { user, loading: authLoading } = useAuth()
  const navigate = useNavigate()
  const [species, setSpecies] = useState('CAT')
  const [breeds, setBreeds] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!user) {
      setLoading(false)
      return
    }
    listBreeds(species)
      .then((res) => setBreeds(res?.breeds ?? []))
      .catch(() => setBreeds([]))
      .finally(() => setLoading(false))
  }, [user, species])

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
      <Header title="领养" showBack />
      <main className="p-4">
        <div className="flex gap-2 mb-4">
          <button
            type="button"
            onClick={() => setSpecies('CAT')}
            className={`flex-1 py-3 rounded-xl font-medium ${
              species === 'CAT' ? 'bg-primary text-white' : 'glass'
            }`}
          >
            猫咪
          </button>
          <button
            type="button"
            onClick={() => setSpecies('DOG')}
            className={`flex-1 py-3 rounded-xl font-medium ${
              species === 'DOG' ? 'bg-primary text-white' : 'glass'
            }`}
          >
            狗狗
          </button>
        </div>
        {loading ? (
          <p className="text-center text-gray-500 py-8">加载中...</p>
        ) : breeds.length === 0 ? (
          <p className="text-center text-gray-500 py-8">暂无该类型品种</p>
        ) : (
          <ul className="space-y-2">
            {breeds.map((b) => (
              <li key={b.id}>
                <button
                  type="button"
                  onClick={() => navigate(`/breed-detail/${b.id}`)}
                  className="w-full flex items-center gap-4 py-3 px-4 rounded-xl glass text-left"
                >
                  <div className="w-14 h-14 rounded-xl bg-gray-100 overflow-hidden shrink-0">
                    {b.image_url ? (
                      <img src={b.image_url} alt="" className="w-full h-full object-cover" />
                    ) : (
                      <span className="material-symbols-outlined text-3xl text-gray-400 block text-center leading-14">
                        pets
                      </span>
                    )}
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-medium">{b.name_cn || b.name_en}</p>
                    <p className="text-xs text-gray-500">{b.size_tag} · {b.life_span}</p>
                  </div>
                  <span className="material-symbols-outlined text-gray-400">chevron_right</span>
                </button>
              </li>
            ))}
          </ul>
        )}
      </main>
    </div>
  )
}
