import { BrowserRouter, Routes, Route } from 'react-router-dom'
import SearchPage from './pages/SearchPage'
import LandingPage from '@/pages/LangingPage'
import { Toaster } from '@/components/ui/toaster'
import { ProtectedRoute } from '@/components/auth/ProtectedRoute'
import { UserProvider } from '@/contexts/UserContext'

export default function App() {
  return (
    <UserProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/cards" element={
            <ProtectedRoute>
              <SearchPage />
            </ProtectedRoute>
          } />
        </Routes>
        <Toaster />
      </BrowserRouter>
    </UserProvider>
  )
}