import { Navigate, useLocation } from 'react-router-dom'
import { useUser } from '@/contexts/UserContext'

// const API_URL = import.meta.env.VITE_API_URL // No longer needed here

interface ProtectedRouteProps {
  children: React.ReactNode
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { user, loading } = useUser()
  const location = useLocation()

  // useEffect removed

  if (loading) {
    // Or a global spinner, or null to render nothing while loading
    return null
  }

  if (!user) {
    return <Navigate to="/" state={{ from: location }} replace />
  }

  return <>{children}</>
} 