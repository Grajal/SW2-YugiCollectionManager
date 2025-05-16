import { BrowserRouter, Routes, Route } from 'react-router-dom'
import SearchPage from './pages/SearchPage'
import LandingPage from '@/pages/LangingPage'
import { Toaster } from '@/components/ui/toaster'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/cards" element={<SearchPage />} />
      </Routes>
      <Toaster />
    </BrowserRouter>
  )
}