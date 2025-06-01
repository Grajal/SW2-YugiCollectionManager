"use client"

import type React from "react"
import DeckCard from "./cardDeck"
import { CardDeck } from "@/types/deck"

interface DeckSectionProps {
  title: string
  cards: CardDeck[]
  maxCards: number
  onCardClick?: (card: CardDeck) => void
  isMainDeck: boolean
}

const DeckSection: React.FC<DeckSectionProps> = ({ title, cards, maxCards, onCardClick, isMainDeck }) => {

  const cardsLength = () => {
    var number = 0
    cards.forEach(element => {
    number += element.Quantity
  });
    return number    
}

  return (
    <div className="bg-gray-800 rounded-lg p-6 shadow-lg">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-bold text-white">{title}</h2>
        <span className="text-sm text-gray-400">
          {cardsLength()}/{maxCards}
        </span>
      </div>

      <div
        className={`grid gap-2 ${isMainDeck ? "grid-cols-4 sm:grid-cols-6 md:grid-cols-8 lg:grid-cols-10" : "grid-cols-3 sm:grid-cols-4 md:grid-cols-5"}`}
      >
        {cards.map((card, index) => (
          <DeckCard key={`${card.CardID}-${index}`} card={card} onClick={() => onCardClick?.(card)} />
        ))}
      </div>

      {cards.length === 0 && (
        <div className="text-center py-8">
          <p className="text-gray-500">No hay cartas en este compartimento</p>
        </div>
      )}
    </div>
  )
}

export default DeckSection
