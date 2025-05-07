import { BrowserRouter, Routes, Route } from 'react-router-dom'
import LandingPage from '@/pages/LadingPage';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
      </Routes>
    </BrowserRouter>
  );
}