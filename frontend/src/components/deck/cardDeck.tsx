"use client"

import type React from "react"
import { CardDeck } from "@/types/deck"

interface DeckCardProps {
  card: CardDeck
  onClick?: () => void
}

const DeckCard: React.FC<DeckCardProps> = ({ card, onClick }) => {
  return (
    <div
      className="relative bg-gray-700 rounded-lg overflow-hidden shadow-md hover:shadow-lg transition-all hover:scale-105 cursor-pointer group"
      onClick={onClick}
    >
      {/* Solo la imagen de la carta */}
      <div className="relative pb-[140%]">
        <img
          src={card.Card.ImageURL || "/placeholder.svg?height=300&width=200"}
          alt={card.Card.Name}
          className="absolute inset-0 w-full h-full object-cover"
        />
      </div>
      {/* Contador de cantidad */}
        {card.Quantity > 1 && (
          <div className="absolute top-0 right-0 bg-blue-600 text-white rounded-full w-6 h-6 flex items-center justify-center text-xs font-bold">
            {card.Quantity}
          </div>
        )}
    </div>
  )
}

export default DeckCard
