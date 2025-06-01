import { createContext, useEffect, ReactNode, useState, useCallback } from 'react'

interface User {
  ID: number
  Username: string
  Email: string
}

export interface UserContextType {
  user: User | null
  loading: boolean
  error: string | null
  setCurrentUser: (userData: User | null) => void
}

export const UserContext = createContext<UserContextType | undefined>(undefined)

const API_URL = import.meta.env.VITE_API_URL

export function UserProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchUser = useCallback(async () => {
    setLoading(true)
    try {
      const response = await fetch(`${API_URL}/auth/me`, {
        credentials: 'include',
      })

      if (!response.ok) {
        if (user) {
          setError('Error refreshing user session, using cached data.')
          return
        }
        setUser(null)
        setError(response.status === 401 ? 'No autorizado' : 'Error al obtener datos del usuario')
        return
      }

      const userData = await response.json()
      setUser(userData)
      setError(null)
    } catch (error) {
      console.error('Error al cargar usuario:', error)
      setUser(null)
      setError('Error de conexión al cargar datos del usuario')
    } finally {
      setLoading(false)
    }
  }, [user, setUser, setError, setLoading])

  const setCurrentUser = (userData: User | null) => {
    setUser(userData)
    setLoading(false)
    setError(null)
  }

  useEffect(() => {
    if (!user) {
      fetchUser()
    } else {
      setLoading(false)
    }
  }, [user, fetchUser])


  return (
    <UserContext.Provider value={{ user, loading, error, setCurrentUser }}>
      {children}
    </UserContext.Provider>
  )
}

