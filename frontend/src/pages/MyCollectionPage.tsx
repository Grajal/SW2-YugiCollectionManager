import { Header } from "@/components/landing/header"
import { useState, useEffect } from "react"
import { useUser } from '@/contexts/UserContext'
import type { SearchResult } from "@/types/search"
import type { Deck, Collection } from "@/types/collection"

const API_URL = import.meta.env.VITE_API_URL

export default function MyCollectionPage() {
  const { user } = useUser()
  const [collection, setCollection] = useState<Collection>([])
  const [decks, setDecks] = useState<Deck[]>([])
  const [selectedDeck, setSelectedDeck] = useState<Deck | null>(null)
  const [loading, setLoading] = useState<boolean>(true)
  const [error, setError] = useState<string | null>(null)
  const [newDeckName, setNewDeckName] = useState<string>("")

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
            throw new Error('Failed to fetch collection')
          }
          const collectionData = await collectionResponse.json()
          setCollection(collectionData.collection || [])

          const decksResponse = await fetch(`${API_URL}/decks/`, {
            credentials: 'include',
          })
          if (!decksResponse.ok) {
            throw new Error('Failed to fetch decks')
          }
          const decksData = await decksResponse.json()
          setDecks(decksData)

        } catch (err) {
          console.error('Error fetching data:', err)
          setError(err instanceof Error ? err.message : 'An unknown error occurred')
          setCollection([])
          setDecks([])
        } finally {
          setLoading(false)
        }
      }
      fetchData()
    } else if (!user) {
      setLoading(false)
      setError("Please log in to view your collection and decks.")
    }
  }, [user])

  const handleCreateDeck = async () => {
    // TODO: Implement handleCreateDeck function
    console.log('TODO: Implement handleCreateDeck function')
  }

  const addCardToDeck = async (card: SearchResult) => {
    // TODO: Implement addCardToDeck function
    console.log(card)
  }

  const removeCardFromDeck = async (cardIndex: number) => {
    // TODO: Implement removeCardFromDeck function
    console.log(cardIndex)
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 text-gray-100 flex flex-col items-center justify-center">
        <Header username={user?.Username || ''} />
        <p className="text-xl">Loading your collection and decks...</p>
      </div>
    )
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
        <h1 className="text-3xl font-bold mb-8 text-center">My Collection & Decks</h1>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* Columna 1: Colección de cartas del usuario */}
          <div className="md:col-span-2 bg-gray-800 p-6 rounded-lg shadow-xl">
            <h2 className="text-2xl font-semibold mb-6 text-center">My Cards ({collection.length})</h2>
            {collection.length === 0 ? (
              <p className="text-center text-gray-400">Your collection is empty. Add cards from the Catalog!</p>
            ) : (
              <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 max-h-[70vh] overflow-y-auto pr-2">
                {collection.map((card) => (
                  <div key={card.ID} className="bg-gray-700 p-2 rounded-lg shadow-md hover:shadow-lg transition-shadow cursor-pointer hover:ring-2 hover:ring-purple-500"
                    onClick={() => addCardToDeck(card.Card)} title={`Add ${card.Card.Name} to deck`}>
                    <img src={card.Card.ImageURL} alt={card.Card.Name} className="w-full h-auto rounded" />
                    <h3 className="text-xs font-semibold mt-1 truncate" title={card.Card.Name}>{card.Card.Name}</h3>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Columna 2: Gestión de mazos */}
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
                    {selectedDeck.DeckCards.map((card, index) => (
                      <div key={`${card.ID} -${index} `} className="flex items-center justify-between bg-gray-700 p-2 rounded-lg">
                        <span className="text-sm truncate" title={card.Name}>{card.Name}</span>
                        <button onClick={() => removeCardFromDeck(index)} className="text-red-400 hover:text-red-300 text-xs">Eliminar</button>
                      </div>
                    ))}
                  </div>
                )}
                {/* Mostrar el total de cartas en el mazo, mazo principal, mazo extra, mazo lateral podría ir aquí */}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
} 