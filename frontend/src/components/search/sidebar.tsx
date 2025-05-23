"use client"

import type React from "react"
import { useEffect } from "react"
import type { SearchResult } from "../../types/search"
import { X } from "lucide-react"

interface DetailsSidebarProps {
  card: SearchResult | null
  isOpen: boolean
  onClose: () => void
}

export const Sidebar: React.FC<DetailsSidebarProps> = ({ card, isOpen, onClose }) => {
  // Bloquear el scroll del body cuando el sidebar está abierto
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = "hidden"
    } else {
      document.body.style.overflow = "auto"
    }

    return () => {
      document.body.style.overflow = "auto"
    }
  }, [isOpen])

  if (!card) return null

  return (
    <>
      {/* Overlay */}
      <div
        className={`fixed inset-0 bg-black bg-opacity-50 z-40 transition-opacity duration-300 ${
          isOpen ? "opacity-100" : "opacity-0 pointer-events-none"
        }`}
        onClick={onClose}
      />

      <div
        className={`fixed top-0 right-0 h-full w-full md:w-96 bg-gray-800 shadow-xl z-50 overflow-y-auto transition-transform duration-300 transform ${
          isOpen ? "translate-x-0" : "translate-x-full"
        }`}
      >
        <div className="p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-bold text-white">Detalles de la Carta</h2>
            <button onClick={onClose} className="text-gray-400 hover:text-white">
              <X size={24} />
            </button>
          </div>

          <div className="flex flex-col items-center mb-6">
            <img
              src={card.image || "/placeholder.svg?height=400&width=280"}
              className="w-48 h-auto mb-4 rounded-lg shadow-md"
            />
            <h3 className="text-xl font-bold text-white">{card.name}</h3>
          </div>

          <div className="space-y-4">
            <div>
              <h4 className="text-sm font-medium text-gray-400">Tipo</h4>
              <p className="text-white">{card.tipo}</p>
            </div>

            <div>
              <h4 className="text-sm font-medium text-gray-400">Arquetipo</h4>
              <p className="text-white">{card.arquetipo || "N/A"}</p>
            </div>

            <div>
              <h4 className="text-sm font-medium text-gray-400">Atributo</h4>
              <p className="text-white">{card.atributo || "N/A"}</p>
            </div>

            {card.estrellas && (
              <div>
                <h4 className="text-sm font-medium text-gray-400">Estrellas</h4>
                <p className="text-white">{card.estrellas}</p>
              </div>
            )}


            {(card.atk !== undefined || card.def !== undefined) && (
              <div>
                <h4 className="text-sm font-medium text-gray-400">ATK / DEF</h4>
                <p className="text-white">
                  {card.atk || "?"} / {card.def || "?"}
                </p>
              </div>
            )}

            <div>
              <h4 className="text-sm font-medium text-gray-400">Descripción</h4>
              <p className="text-white text-sm mt-1">{card.descripcion}</p>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
