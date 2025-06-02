import { Header } from "@/components/landing/header"
import { useUser } from '@/hooks/useUser'
import { useEffect, useState } from "react"
import { CardDeck, type Deck } from "@/types/deck"
import { DeckGrid } from "@/components/deck/deckSelector"
import DeckViewer from "@/components/deck/deckViewer"
import { Sidebar } from "@/components/search/sidebar"

const API_URL = import.meta.env.VITE_API_URL

export default function Decks() {
  const { user } = useUser()
  const [decks, setDecks] = useState<Deck[]>([])
  const [deckEditor, setDeckEditor] = useState<boolean>(true)
  const [selectedDeck, setSelectedDeck] = useState<Deck | null>(null)
  const [mainDeck, setMainDeck] = useState<CardDeck[]>([])
  const [extraDeck, setExtraDeck] = useState<CardDeck[]>([])
  const [sideDeck, setSideDeck] = useState<CardDeck[]>([])
  const [selectedCard, setSelectedCard] = useState<CardDeck | null>(null)
  const [isSidebarOpen, setIsSidebarOpen] = useState<boolean>(false)

  const fetchDecks = async () => {
    try {
      const response = await fetch(`${API_URL}/decks/`, {
        method: 'GET',
        credentials: 'include',
      })
      if (!response.ok) {
        throw new Error("Error al cargar los datos")
      }
      const data = await response.json()
      setDecks(data)
      console.log(data)
    } catch (error) {
      console.error("ERROR CARGANDO DECKS: " + error)
    }
  }

  const handleDeckCards = () => {
    try {
      const main: CardDeck[] = []
      const extra: CardDeck[] = []
      const side: CardDeck[] = []
      selectedDeck?.DeckCards.forEach(element => {
        switch (element.Zone) {
          case "main":
            main.push(element)
            break

          case "extra":
            extra.push(element)
            break

          default:
            side.push(element)
            break
        }
      })
      setMainDeck(main)
      setExtraDeck(extra)
      setSideDeck(side)
    } catch (error) {
      console.error(error)
    }
  }

  const handleDeckClick = (deck: Deck) => {
    setSelectedDeck(deck)
    if (!deckEditor) {
      setDeckEditor(true)
    }
  }

  const handleSidebarClose = () => {
    setIsSidebarOpen(false)
    setSelectedCard(null)
  }

  const removeCard = async (numE?: number) => {
    if (selectedCard && selectedCard.Quantity > 0) {
      try {
        const response = await fetch(`${API_URL}/decks/${selectedDeck?.ID}/cards/${selectedCard.CardID}`, {
          method: 'DELETE',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            quantity: numE === undefined ? 1 : numE,
          }),
        })

        const data = await response.json()

        if (!response.ok) {
          throw new Error(data.error || 'Error removing card')
        }
        fetchDecks()
        handleDeckCards()
      } catch (error) {
        console.error('Failed to remove card:', error)
      }
    }
  }

  const addCard = async (numE?: number) => {
    if (selectedCard && selectedCard.Quantity > 0) {
      try {
        const response = await fetch(`${API_URL}/decks/${selectedDeck?.ID}/cards`, {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            card_id: selectedCard.CardID,
            quantity: numE === undefined ? 1 : numE,
          }),
        })

        const data = await response.json()

        if (!response.ok) {
          throw new Error(data.error || 'Error adding card')
        }
        fetchDecks()
        handleDeckCards()
      } catch (error) {
        console.error('Failed to add card:', error)
      }
    }
  }

  const handleCardClick = (card: CardDeck) => {
    setSelectedCard(card)
    setIsSidebarOpen(true)
  }


  useEffect(() => {
    if (selectedDeck || selectedCard?.Quantity) {
      handleDeckCards()
    }
  }, [selectedDeck])



  return (
    <div className="min-h-screen bg-gray-900 text-gray-100" onLoad={fetchDecks}>
      <Header username={user?.Username || ''} />
      <DeckGrid results={decks} onDeckClick={handleDeckClick}></DeckGrid>
      {deckEditor && selectedDeck != null && (<DeckViewer deck={selectedDeck} mainDeck={mainDeck} extraDeck={extraDeck} sideDeck={sideDeck} onCardClick={(card) => { handleCardClick(card) }}></DeckViewer>)}
      {selectedCard != null && (<Sidebar type="deck" card={selectedCard.Card} isOpen={isSidebarOpen} onClose={handleSidebarClose} onAction={(quantity) => { removeCard(quantity) }} onAdd={(quantity) => { addCard(quantity) }} onAddToCollection={() => { }} onQuantityChange={() => { }}></Sidebar>)}
    </div>
  )
}