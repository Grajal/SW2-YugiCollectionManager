import { BrowserRouter, Routes, Route } from 'react-router-dom'
import CatalogPage from '@/pages/CatalogPage'
import LandingPage from '@/pages/LangingPage'
import { Toaster } from '@/components/ui/toaster'
import { ProtectedRoute } from '@/components/auth/ProtectedRoute'
import { UserProvider } from '@/contexts/UserContext'
import MyCollectionPage from '@/pages/MyCollectionPage'
import Decks from './pages/DecksPage'

export default function App() {
  return (
    <UserProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/cards" element={
            <ProtectedRoute>
              <CatalogPage />
            </ProtectedRoute>
          } />
          <Route path="/collection" element={
            <ProtectedRoute>
              <MyCollectionPage />
            </ProtectedRoute>
          } />
          <Route path="/decks" element={
            <ProtectedRoute>
              <Decks />
            </ProtectedRoute>
          } />
        </Routes>
        <Toaster />
      </BrowserRouter>
    </UserProvider>
  )
}