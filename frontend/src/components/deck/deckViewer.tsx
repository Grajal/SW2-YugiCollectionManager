"use client"

import type React from "react"
import DeckSection from "./cardsSection"
import type { CardDeck, Deck } from "@/types/deck"
import { useState } from "react"
import UploadYDKDialog from "../ui/import"

interface DeckViewerProps {
  deck: Deck
  mainDeck: CardDeck[]
  extraDeck: CardDeck[]
  onCardClick?: (card: CardDeck) => void
}

const API_URL = import.meta.env.VITE_API_URL

const DeckViewer: React.FC<DeckViewerProps> = ({ deck, mainDeck, extraDeck, onCardClick }) => {
  const [isImportOpen, setIsImportOpen] = useState<boolean>(false)
  const removeDeck = async () => {
      try {
        const response = await fetch(`${API_URL}/decks/${deck.ID}`, {
          method: 'DELETE',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
          },
        })

        const data = await response.json()

        if (!response.ok) {
          throw new Error(data.error || 'Error removing card')
        }
      } catch (error) {
        console.error('Failed to remove card:', error)
      }
  }

  const importDeck = async (selectedFile: File) => {
    const formData = new FormData()
    formData.append("file", selectedFile)
    try {
      const response = await fetch(`${API_URL}/decks/import/${deck.ID}`, {
        method: 'POST',
        body: formData,
        credentials: 'include',
      })
      if (!response.ok) {
        throw new Error("Error al cargar los datos")
      }
    } catch (error) {
      console.error(error)
    }
  }

  const exportDeck = async () => {
    if (!deck) return

    try {
      const response = await fetch(`${API_URL}/decks/export/${deck.ID}`, {
        method: 'POST',
        credentials: 'include',
      })

      if (!response.ok) {
        throw new Error('Error exportando el deck')
      }

      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)

      const a = document.createElement('a')
      a.href = url
      a.download = `${deck.Name}.ydk`
      a.click()
      window.URL.revokeObjectURL(url)
    } catch (error) {
      console.error('Error exportando el deck:', error)
    }
  }

  const handleImportClick = () => {
    setIsImportOpen(!isImportOpen)
    console.log(isImportOpen)
  }

  return (
    <div className="bg-gray-900 min-h-screen p-6">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold text-white mb-2">{deck.Name}</h1>
        </div>
        <div className="justify-start gap-2">
            <button
                    onClick={removeDeck}
                    className="ml-6 cursor-pointer bg-red-600 hover:bg-red-700 text-white font-bold py-3 px-4 rounded-lg transition duration-150 ease-in-out"
                  >
                    Eliminar
            </button>
            {deck.DeckCards.length === 0 ? 
            <button 
              onClick={handleImportClick} 
              className="ww-full cursor-pointer bg-white-600 hover:bg-white-700 text-white font-bold py-3 px-4 rounded-md transition duration-150 ease-in-out">
              Importar
            </button>
            : 
            <button 
              onClick={exportDeck} 
              className="ww-full cursor-pointer bg-white-600 hover:bg-white-700 text-white font-bold py-3 px-4 rounded-md transition duration-150 ease-in-out">
              Exportar
            </button>
}
          </div>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Deck Principal - Ocupa 2/3 del espacio */}
          <div className="lg:col-span-2">
            <DeckSection
              title="Deck Principal"
              cards={mainDeck}
              maxCards={60}
              onCardClick={onCardClick}
              isMainDeck={true}
            />
          </div>

          {/* Extra Deck - Ocupa 1/3 del espacio */}
          <div className="lg:col-span-1">
            <DeckSection
              title="Extra Deck"
              cards={extraDeck}
              maxCards={15}
              onCardClick={onCardClick}
              isMainDeck={false}
            />
          </div>
        </div>
      </div>
      {isImportOpen && (<UploadYDKDialog open={isImportOpen}
        onOpenChange={setIsImportOpen}
        onFileSelected={(file) => { importDeck(file) }}>
      </UploadYDKDialog>)}
    </div>
  )
}

export default DeckViewer
