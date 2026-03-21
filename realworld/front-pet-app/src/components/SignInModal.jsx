import { useState, useEffect } from 'react'
import { doCheckin, getCheckinStatus } from '../api/checkin'

const DAYS = ['一', '二', '三', '四', '五', '六', '日']

export default function SignInModal({ open, onClose, onSuccess }) {
  const [status, setStatus] = useState(null)
  const [loading, setLoading] = useState(false)
  const [done, setDone] = useState(false)
  const [result, setResult] = useState(null)

  const loadStatus = async () => {
    try {
      const res = await getCheckinStatus()
      setStatus(res)
    } catch (_) {
      setStatus(null)
    }
  }

  // 当 open 变为 true 时拉取状态
  useEffect(() => {
    if (open) {
      loadStatus()
      setDone(false)
      setResult(null)
    }
  }, [open])

  const handleCheckin = async () => {
    setLoading(true)
    try {
      const res = await doCheckin()
      setResult(res)
      setDone(true)
      loadStatus()
      onSuccess?.(res)
    } catch (e) {
      setResult({ already_checked_today: true, message: e?.data?.message || '签到失败' })
      setDone(true)
    } finally {
      setLoading(false)
    }
  }

  if (!open) return null

  const last7 = status?.last_7_days ?? []
  const checkedToday = status?.checked_today ?? false
  const continuous = status?.continuous_days ?? 0

  return (
    <div className="fixed inset-0 z-[100] flex items-end justify-center bg-black/40" onClick={onClose}>
      <div
        className="w-full max-w-app rounded-t-3xl bg-white shadow-xl p-6 pb-10"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">每日签到</h2>
          <button type="button" onClick={onClose} className="p-1 rounded-lg hover:bg-gray-100">
            <span className="material-symbols-outlined">close</span>
          </button>
        </div>
        <p className="text-gray-600 text-sm mb-2">连续签到 {continuous} 天</p>
        <div className="flex gap-2 mb-6">
          {DAYS.map((d, i) => (
            <div
              key={d}
              className={`flex-1 py-2 rounded-lg text-center text-sm ${
                last7[i] ? 'bg-primary text-white' : 'bg-gray-100 text-gray-400'
              }`}
            >
              {d}
            </div>
          ))}
        </div>
        {!done && (
          <button
            type="button"
            onClick={handleCheckin}
            disabled={loading || checkedToday}
            className="w-full py-3 rounded-xl bg-primary text-white font-medium disabled:opacity-50"
          >
            {loading ? '签到中...' : checkedToday ? '今日已签到' : '立即签到领奖励'}
          </button>
        )}
        {done && result && (
          <div className="rounded-xl bg-gray-50 p-4 text-center">
            {result.already_checked_today ? (
              <p className="text-gray-600">今日已签到过啦，明天再来～</p>
            ) : (
              <p className="text-primary font-medium">
                +{result.coins_reward ?? 0} 金币 +{result.heart_reward ?? 0} 爱心
              </p>
            )}
            <button
              type="button"
              onClick={onClose}
              className="mt-3 text-sm text-gray-500 underline"
            >
              关闭
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
