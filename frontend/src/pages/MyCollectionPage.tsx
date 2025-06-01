import { Header } from "@/components/landing/header"
import { useState, useEffect } from "react"
import { useUser } from '@/contexts/UserContext'
import type { Collection } from "@/types/collection"
import type { Deck } from '@/types/deck'
import { ManageCardModal } from "@/components/collection/ManageCardModal"
import { useCollectionManagement } from "@/hooks/useCollectionManagement"
import { CardDeck } from '@/types/deck'

const API_URL = import.meta.env.VITE_API_URL

export default function MyCollectionPage() {
  const { user } = useUser()
  const [collection, setCollection] = useState<Collection>([])
  const [decks, setDecks] = useState<Deck[]>([])
  const [selectedDeck, setSelectedDeck] = useState<Deck | null>(null)
  const [loading, setLoading] = useState<boolean>(true)
  const [error, setError] = useState<string | null>(null)
  const [newDeckName, setNewDeckName] = useState<string>("")

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
          const decksData = await decksResponse.json()
          setDecks(decksData)

        } catch (err) {
          console.error('Error fetching data:', err)
          setError(err instanceof Error ? err.message : 'Ocurrió un error desconocido')
          setCollection([])
          setDecks([])
        } finally {
          setLoading(false)
        }
      }
      fetchData()
    } else if (!user) {
      setLoading(false)
      setError("Por favor, inicia sesión para ver tu colección y mazos.")
      setCollection([])
      setDecks([])
    }
  }, [user, setError, setCollection])

  const handleCreateDeck = async () => {
    console.log('TODO: Implement handleCreateDeck function')
  }

  const removeCardFromDeck = async (cardIndex: number) => {
    console.log(cardIndex)
  }

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

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="md:col-span-2 bg-gray-800 p-6 rounded-lg shadow-xl">
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

          <div className="bg-gray-800 p-6 rounded-lg shadow-xl">
            <h2 className="text-2xl font-semibold mb-6 text-center">Mis mazos</h2>
            <div className="mb-6">
              <input
                type="text"
                value={newDeckName}
                onChange={(e) => setNewDeckName(e.target.value)}
                placeholder="Nombre del mazo"
                className="w-full px-4 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 mb-2"
              />
              <button
                onClick={handleCreateDeck}
                className="w-full bg-purple-600 cursor-pointer purple-600 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded-lg transition duration-150"
              >
                Crear nuevo mazo
              </button>
            </div>

            {decks.length === 0 && !selectedDeck ? (
              <p className="text-center text-gray-400">No tienes mazos. Crea uno para empezar!</p>
            ) : (
              <div className="mb-6 max-h-40 overflow-y-auto pr-1">
                {decks.map((deck) => (
                  <button
                    key={deck.ID}
                    onClick={() => setSelectedDeck(deck)}
                    className={`w - full text - left p - 3 mb - 2 rounded - lg transition - colors ${selectedDeck?.ID === deck.ID ? 'bg-purple-700 text-white' : 'bg-gray-700 hover:bg-gray-600'} `}
                  >
                    {deck.Name}
                  </button>
                ))}
              </div>
            )}

            {selectedDeck && (
              <div>
                <h3 className="text-xl font-semibold mb-4">Editando: {selectedDeck.Name}</h3>
                {selectedDeck.DeckCards.length === 0 ? (
                  <p className="text-center text-gray-400">Este mazo está vacío. Haz click en cartas de tu colección para agregarlas.</p>
                ) : (
                  <div className="space-y-2 max-h-[40vh] overflow-y-auto pr-1">
                    {selectedDeck.DeckCards.map((card: CardDeck, index: number) => (
                      <div key={`${card.CardID} -${index} `} className="flex items-center justify-between bg-gray-700 p-2 rounded-lg">
                        <span className="text-sm truncate" title={card.Card.Name}>{card.Card.Name}</span>
                        <button onClick={() => removeCardFromDeck(index)} className="text-red-400 hover:text-red-300 text-xs">Eliminar</button>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
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