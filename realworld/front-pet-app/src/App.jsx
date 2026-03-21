import { BrowserRouter, Routes, Route, useLocation } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import BottomNav from './components/BottomNav'

import HomePage from './pages/HomePage'
import InteractionPage from './pages/InteractionPage'
import SocialPage from './pages/SocialPage'
import ProfilePage from './pages/ProfilePage'
import LoginPage from './pages/LoginPage'
import ShopPage from './pages/ShopPage'
import InventoryPage from './pages/InventoryPage'
import AdoptionPage from './pages/AdoptionPage'
import BreedDetailPage from './pages/BreedDetailPage'
import EditProfilePage from './pages/EditProfilePage'
import PurchaseSuccessPage from './pages/PurchaseSuccessPage'
import AchievementsPage from './pages/AchievementsPage'
import DressUpPage from './pages/DressUpPage'
import RankingPage from './pages/RankingPage'
import RequestsPage from './pages/RequestsPage'

const BOTTOM_NAV_PATHS = ['/', '/social', '/interaction', '/profile']

function Layout({ children }) {
  const location = useLocation()
  const showNav = BOTTOM_NAV_PATHS.includes(location.pathname)

  return (
    <>
      {children}
      {showNav && <BottomNav />}
    </>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Layout>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/interaction" element={<InteractionPage />} />
            <Route path="/social" element={<SocialPage />} />
            <Route path="/profile" element={<ProfilePage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/shop" element={<ShopPage />} />
            <Route path="/inventory" element={<InventoryPage />} />
            <Route path="/adoption" element={<AdoptionPage />} />
            <Route path="/breed-detail/:id" element={<BreedDetailPage />} />
            <Route path="/edit-profile" element={<EditProfilePage />} />
            <Route path="/purchase-success" element={<PurchaseSuccessPage />} />
            <Route path="/achievements" element={<AchievementsPage />} />
            <Route path="/dressup" element={<DressUpPage />} />
            <Route path="/ranking" element={<RankingPage />} />
            <Route path="/requests" element={<RequestsPage />} />
          </Routes>
        </Layout>
      </AuthProvider>
    </BrowserRouter>
  )
}
