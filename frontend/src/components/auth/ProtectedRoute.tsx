import { useEffect, useState } from 'react'
import { Navigate, useLocation } from 'react-router-dom'

const API_URL = import.meta.env.VITE_API_URL

interface ProtectedRouteProps {
  children: React.ReactNode
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null)
  const location = useLocation()

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const response = await fetch(`${API_URL}/auth/me`, {
          credentials: 'include',
        })
        setIsAuthenticated(response.ok)
      } catch (error) {
        console.error('Error al verificar la autenticación:', error)
        setIsAuthenticated(false)
      }
    }

    checkAuth()
  }, [])

  if (isAuthenticated === null) {
    return <div>Cargando...</div>
  }

  if (!isAuthenticated) {
    return <Navigate to="/" state={{ from: location }} replace />
  }

  return <>{children}</>
} 