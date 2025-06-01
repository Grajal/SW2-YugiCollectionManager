"use client"

import type React from "react"
import type { Deck } from "../../types/deck"

interface DeckResultProps {
  deck: Deck
  onClick: () => void
}

export const DeckResult: React.FC<DeckResultProps> = ({ deck, onClick }) => {
  return (
    <div
      className="w-8/100 h-4/100 overflow-hidden shadow-md hover:shadow-lg transition-all hover:scale-105 cursor-pointer"
      onClick={onClick}
    >
      <div className="relative pb-[140%]">
        {" "}
        <img
          src={"https://carduniverse.cl/cdn/shop/products/DeckBoxDiosespng-min_1024x1024@2x.png?v=1652570195"}
          className="absolute inset-0 w-full h-full object-cover"
        />
      </div>
      <div className="p-3 text-center">
        <h3 className="text-sm font-medium text-white truncate" title={deck.Name}>
          {deck.Name}
        </h3>
      </div>
    </div>
  )
}
