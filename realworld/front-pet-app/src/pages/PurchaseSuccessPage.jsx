import { useNavigate, useLocation } from 'react-router-dom'
import Header from '../components/Header'

export default function PurchaseSuccessPage() {
  const navigate = useNavigate()
  const { state } = useLocation()
  const item = state?.item
  const quantity = state?.quantity ?? 1

  return (
    <div className="min-h-screen pb-20 flex flex-col">
      <Header title="购买成功" showBack />
      <main className="flex-1 flex flex-col items-center justify-center p-6">
        <span className="material-symbols-outlined text-6xl text-primary mb-4">check_circle</span>
        <h2 className="text-xl font-semibold mb-2">购买成功</h2>
        {item && (
          <p className="text-gray-600">
            获得 {item.name} x{quantity}
          </p>
        )}
        <button
          type="button"
          onClick={() => navigate('/inventory')}
          className="mt-6 px-6 py-3 rounded-xl bg-primary text-white font-medium"
        >
          去背包查看
        </button>
        <button
          type="button"
          onClick={() => navigate('/shop')}
          className="mt-3 text-gray-500 text-sm"
        >
          继续逛逛
        </button>
      </main>
    </div>
  )
}
