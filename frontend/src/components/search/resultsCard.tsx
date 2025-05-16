"use client"

import type React from "react"
import type { SearchResult } from "../../types/search"

interface CardResultProps {
  card: SearchResult
  onClick: () => void
}

export const CardResult: React.FC<CardResultProps> = ({ card, onClick }) => {
  return (
    <div
      className="bg-gray-800 rounded-lg overflow-hidden shadow-md hover:shadow-lg transition-all hover:scale-105 cursor-pointer"
      onClick={onClick}
    >
      <div className="relative pb-[140%]">
        {" "}
        <img
          src={card.image || "/placeholder.svg?height=300&width=200"}
          className="absolute inset-0 w-full h-full object-cover"
        />
      </div>
      <div className="p-3 text-center">
        <h3 className="text-sm font-medium text-white truncate" title={card.name}>
          {card.name}
        </h3>
        <p className="text-xs text-gray-400 mt-1">{card.tipo}</p>
      </div>
    </div>
  )
}
