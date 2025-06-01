import { useContext } from 'react'
import { UserContext, type UserContextType } from '@/contexts/UserContext'

export function useUser(): UserContextType {
  const context: UserContextType | undefined = useContext(UserContext)
  if (context === undefined) {
    throw new Error('useUser debe ser usado dentro de un UserProvider')
  }
  return context
}
