import { BrowserRouter, Routes, Route } from 'react-router-dom'
import LandingPage from '@/pages/LadingPage';
import SearchPage from './pages/SearchPage';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/cards" element={<SearchPage />} />
      </Routes>
    </BrowserRouter>
  );
}