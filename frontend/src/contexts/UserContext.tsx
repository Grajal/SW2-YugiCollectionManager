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
        throw new Error('No autorizado')
      }

      const userData = await response.json()
      setUser(userData)
      setError(null)
    } catch (error) {
      console.error('Error al cargar usuario:', error)
      setUser(null)
      setError('Error al cargar usuario')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchUser()
  }, [])


  return (
    <UserContext.Provider value={{ user, loading, error }}>
      {children}
    </UserContext.Provider>
  )
}

export function useUser() {
  const context = useContext(UserContext)
  if (context === undefined) {
    throw new Error('useUser debe ser usado dentro de un UserProvider')
  }
  return context
}