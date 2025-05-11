import { BrowserRouter, Routes, Route } from 'react-router-dom'
import LandingPage from '@/pages/LangingPage'
import { Toaster } from '@/components/ui/toaster'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
      </Routes>
      <Toaster />
    </BrowserRouter>
  )
}