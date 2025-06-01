"use client"
import { Header } from "@/components/landing/header"
import { useUser } from '@/contexts/UserContext'
import { useEffect, useState } from "react"
import { CardDeck, type Deck } from "@/types/deck"
import { DeckGrid } from "@/components/deck/deckSelector"
import DeckViewer from "@/components/deck/deckViewer"
import { Sidebar } from "@/components/search/sidebar"

const API_URL = import.meta.env.VITE_API_URL

export default function Decks() {
    const {user} = useUser()
    const[decks,setDecks] = useState<Deck[]>([])
    const [deckEditor, setDeckEditor] = useState<Boolean>(true)
    const [selectedDeck, setSelectedDeck] = useState<Deck | null>(null)
    const [mainDeck, setMainDeck] = useState<CardDeck[]>([])
    const [extraDeck, setExtraDeck] = useState<CardDeck[]>([])
    const [selectedCard, setSelectedCard] = useState<CardDeck | null>(null)
    const [isSidebarOpen, setIsSidebarOpen] = useState<boolean>(false)

    const fetchDecks = async() => {
    try{
      const response = await fetch(`${API_URL}/decks/`,{
        method: 'GET',
        credentials: 'include',
      })
      if(!response.ok){
        throw new Error("Error al cargar los datos");
      }
      const data = await response.json();
      setDecks(data)
      console.log(data)
    }catch(error){
      console.error("ERROR CARGANDO DECKS: " + error)
    }
  }

  const handleDeckCards = () => {
    try{
        var main:CardDeck[] = []
        var extra:CardDeck[] = []
        selectedDeck?.DeckCards.forEach(element => {
            element.Zone === 'main' ? main.push(element) : extra.push(element)
        });
        setMainDeck(main)
        setExtraDeck(extra)
    }catch(error){
        console.error(error)
    }
  }

  const handleDeckClick = (deck:Deck) => {
    setSelectedDeck(deck)
    if(!deckEditor){
        setDeckEditor(true)
    }
  }

  const handleSidebarClose = () => {
    setIsSidebarOpen(false)
    setSelectedCard(null)
  }

 const removeCard = async() => {
  if (selectedCard && selectedCard.Quantity > 0) {
    try {
    const response = await fetch(`${API_URL}/decks/${selectedDeck?.ID}/cards/${selectedCard.CardID}`, {
      method: 'DELETE',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        quantity: 1,
      }),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.error || 'Error removing card')
    }

    console.log('Card removed:', data.message)
  } catch (error) {
    console.error('Failed to remove card:', error)
  }
  }
}

 const handleCardClick = (card:CardDeck) => {
    setSelectedCard(card)
    setIsSidebarOpen(true)
 }


  useEffect(() => {
    if (selectedDeck) {
        handleDeckCards()
    }
    }, [selectedDeck])


    return (
        <div className="min-h-screen bg-gray-900 text-gray-100" onLoad={fetchDecks}>
            <Header username={user?.Username || ''}/>
            <DeckGrid results={decks} onDeckClick={handleDeckClick}></DeckGrid>
            {deckEditor && selectedDeck != null && (<DeckViewer deck={selectedDeck} mainDeck={mainDeck} extraDeck={extraDeck} onCardClick={(card) => {handleCardClick(card)}}></DeckViewer>)}
            {selectedCard != null && (<Sidebar card={selectedCard.Card} isOpen={isSidebarOpen} onClose={handleSidebarClose} onAddToCollection={removeCard}></Sidebar>)}
        </div>
    )
}