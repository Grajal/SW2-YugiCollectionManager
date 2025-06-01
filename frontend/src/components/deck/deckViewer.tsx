"use client"

import type React from "react"
import DeckSection from "./cardsSection"
import type { CardDeck, Deck } from "@/types/deck"

interface DeckViewerProps {
  deck: Deck
  mainDeck: CardDeck[]
  extraDeck: CardDeck[]
  onCardClick?: (card: any) => void
}

const DeckViewer: React.FC<DeckViewerProps> = ({ deck, mainDeck, extraDeck,onCardClick }) => {
  return (
    <div className="bg-gray-900 min-h-screen p-6">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold text-white mb-2">{deck.Name}</h1>
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
    </div>
  )
}

export default DeckViewer
