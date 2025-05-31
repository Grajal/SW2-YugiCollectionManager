"use client"

import type React from "react"
import type { SearchResult } from "@/types/search"
import { CardResult } from "./resultsCard"

interface ResultsGridProps {
  results: SearchResult[]
  onCardClick: (card: SearchResult) => void
}

export const ResultsGrid: React.FC<ResultsGridProps> = ({ results, onCardClick }) => {
  if (results.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-400 text-lg">
          No se encontraron resultados. Intenta con otros términos de búsqueda o filtros.
        </p>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
      {results.map((card) => (
        <CardResult key={card.ID} card={card} onClick={() => onCardClick(card)} />
      ))}
    </div>
  )
}
