import { Header } from "@/components/landing/header"
import { useState, useEffect } from "react"
import { useUser } from '@/hooks/useUser'
import type { Collection } from "@/types/collection"
import { ManageCardModal } from "@/components/collection/ManageCardModal"
import { useCollectionManagement } from "@/hooks/useCollectionManagement"

const API_URL = import.meta.env.VITE_API_URL

export default function MyCollectionPage() {
  const { user } = useUser()
  const [collection, setCollection] = useState<Collection>([])
  const [loading, setLoading] = useState<boolean>(true)
  const [error, setError] = useState<string | null>(null)

  const {
    isModalOpen,
    selectedCard,
    selectedCardQuantity,
    setIsModalOpen,
    setSelectedCardQuantity,
    handleOpenCardModal,
    handleDeleteCard,
    handleUpdateCardQuantity,
  } = useCollectionManagement({ setCollection, setError })

  useEffect(() => {
    if (user && user.ID) {
      const fetchData = async () => {
        setLoading(true)
        setError(null)
        try {
          const collectionResponse = await fetch(`${API_URL}/collections/`, {
            credentials: 'include',
          })
          if (!collectionResponse.ok) {
            throw new Error('Error al cargar la colección')
          }
          const collectionData = await collectionResponse.json()
          setCollection(collectionData.collection || [])

          const decksResponse = await fetch(`${API_URL}/decks/`, {
            credentials: 'include',
          })
          if (!decksResponse.ok) {
            throw new Error('Error al cargar los mazos')
          }

        } catch (err) {
          console.error('Error fetching data:', err)
          setError(err instanceof Error ? err.message : 'Ocurrió un error desconocido')
          setCollection([])
        } finally {
          setLoading(false)
        }
      }
      fetchData()
    } else if (!user) {
      setLoading(false)
      setError("Por favor, inicia sesión para ver tu colección y mazos.")
      setCollection([])
    }
  }, [user, setError, setCollection])

  if (loading) {
    return null
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-900 text-gray-100 flex flex-col items-center">
        <Header username={user?.Username || ''} />
        <div className="container mx-auto px-4 py-8 text-center">
          <p className="text-xl text-red-500">{error}</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100">
      <Header username={user?.Username || ''} />
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-8 text-center">Mi Colección y Mazos</h1>

        <div className="flex flex-col max-w-7xl mx-auto px-4 py-8 w-full h-full">
          
            <h2 className="text-2xl font-semibold mb-6 text-center">Mis Cartas ({collection.length})</h2>
            {collection.length === 0 ? (
              <p className="text-center text-gray-400">Tu colección está vacía. ¡Añade cartas desde el Catálogo!</p>
            ) : (
              <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 max-h-[70vh] overflow-y-auto pr-2">
                {collection.map((item) => (
                  <div key={item.Card.ID} className="relative bg-gray-700 p-2 rounded-lg shadow-md hover:shadow-lg transition-shadow cursor-pointer hover:ring-2 hover:ring-purple-500"
                    onClick={() => handleOpenCardModal(item)} title={`Gestionar ${item.Card.Name}`}>
                    <img src={item.Card.ImageURL} alt={item.Card.Name} className="w-full h-auto rounded" />
                    <h3 className="text-xs font-semibold mt-1 truncate" title={item.Card.Name}>{item.Card.Name}</h3>
                    {item.Quantity > 1 && (
                      <div className="absolute top-0 right-0 bg-purple-600 text-white text-xs font-bold rounded-full h-5 w-5 flex items-center justify-center">
                        {item.Quantity}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
      </div>

      <ManageCardModal
        isOpen={isModalOpen}
        onOpenChange={setIsModalOpen}
        selectedCard={selectedCard}
        selectedCardQuantity={selectedCardQuantity}
        setSelectedCardQuantity={setSelectedCardQuantity}
        handleDeleteCard={handleDeleteCard}
        handleUpdateCardQuantity={handleUpdateCardQuantity}
      />
    </div>
  )
}