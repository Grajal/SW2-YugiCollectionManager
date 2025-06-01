import { createContext, useContext, useEffect, ReactNode, useState } from 'react'

interface User {
  ID: number
  Username: string
  Email: string
}

interface UserContextType {
  user: User | null
  loading: boolean
  error: string | null
  setCurrentUser: (userData: User | null) => void
}

const UserContext = createContext<UserContextType | undefined>(undefined)

const API_URL = import.meta.env.VITE_API_URL

export function UserProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchUser = async () => {
    try {
      const response = await fetch(`${API_URL}/auth/me`, {
        credentials: 'include',
      })

      if (!response.ok) {
        if (user) {
          setError('Error refreshing user session, using cached data.')
          return
        }
        throw new Error('No autorizado')
      }

      const userData = await response.json()
      setUser(userData)
      setError(null)
      setLoading(false)
    } catch (error) {
      console.error('Error al cargar usuario:', error)
      if (!user) {
        setUser(null)
      }
      setError('Error al cargar usuario')
      setLoading(false)
    }
  }

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
  }, [])


  return (
    <UserContext.Provider value={{ user, loading, error, setCurrentUser }}>
      {children}
    </UserContext.Provider>
  )
} export function useUser() {
  const context = useContext(UserContext)
  if (context === undefined) {
    throw new Error('useUser debe ser usado dentro de un UserProvider')
  }
  return context
}

