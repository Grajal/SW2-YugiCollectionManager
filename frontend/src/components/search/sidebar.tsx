"use client"

import type React from "react"
import { useEffect, useState } from "react"
import type { SearchResult } from "../../types/search"
import { X } from "lucide-react"
import { CardDeck } from "@/types/deck"

interface DetailsSidebarProps {
  type?: string
  card: SearchResult | null
  isOpen: boolean
  onClose: () => void
  onAction: (quantity?: number) => void
  onAdd?: (quantity: number) => void
  onAddToCollection: (card: SearchResult) => void
  onQuantityChange: (quantity: number) => void
}

export const Sidebar: React.FC<DetailsSidebarProps> = ({ type = "search", card, isOpen, onClose, onAction, onAdd, onAddToCollection, onQuantityChange }) => {

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

  const [quantity, setQuantity] = useState(1)
  const isCardDeck = (card: SearchResult | CardDeck): card is CardDeck => {
    return (card as CardDeck).Quantity !== undefined
  }


  if (!card) return null

  return (
    <>
      {/* Overlay */}
      <div
        className={`fixed inset-0 bg-black bg-opacity-50 z-40 transition-opacity duration-300 ${isOpen ? "opacity-100" : "opacity-0 pointer-events-none"
          }`}
        onClick={onClose}
      />

      <div
        className={`fixed top-0 right-0 h-full w-full md:w-96 bg-gray-800 shadow-xl z-50 overflow-y-auto transition-transform duration-300 transform ${isOpen ? "translate-x-0" : "translate-x-full"
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
              src={card.ImageURL || "/placeholder.svg?height=400&width=280"}
              className="w-48 h-auto mb-4 rounded-lg shadow-md"
            />
            <h3 className="text-xl font-bold text-white">{card.Name}</h3>
          </div>

          <div className="space-y-4">
            <div>
              <h4 className="text-sm font-medium text-gray-400">Tipo</h4>
              <p className="text-white">{card.Type}</p>
            </div>

            <div>
              <h4 className="text-sm font-medium text-gray-400">Arquetipo</h4>
              <p className="text-white">{card.MonsterCard?.Archetype || "N/A"}</p>
            </div>

            <div>
              <h4 className="text-sm font-medium text-gray-400">Atributo</h4>
              <p className="text-white">{card.MonsterCard?.Attribute || "N/A"}</p>
            </div>

            {card.MonsterCard?.Level && (
              <div>
                <h4 className="text-sm font-medium text-gray-400">Estrellas</h4>
                <p className="text-white">{card.MonsterCard?.Level || "N/A"}</p>
              </div>
            )}


            {(card.MonsterCard?.Atk !== undefined || card.MonsterCard?.Def !== undefined) && (
              <div>
                <h4 className="text-sm font-medium text-gray-400">ATK / DEF</h4>
                <p className="text-white">
                  {card.MonsterCard?.Atk || "?"} / {card.MonsterCard?.Def || "?"}
                </p>
              </div>
            )}

            <div>
              <h4 className="text-sm font-medium text-gray-400">Descripción</h4>
              <p className="text-white text-sm mt-1">{card.Desc}</p>
            </div>

            {/*  Button */}
            <div className="mt-6">
              {type === "search" && (
                <>
                  <div className="mb-4">
                    <label htmlFor="quantity" className="block text-sm font-medium text-gray-400 mb-1">
                      Cantidad
                    </label>
                    <input
                      type="number"
                      id="quantity"
                      name="quantity"
                      min="1"
                      value={quantity}
                      onChange={(e) => onQuantityChange(parseInt(e.target.value, 10))}
                      className="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
                  </div>
                  <button
                    onClick={() => onAddToCollection(card)}
                    className="w-full cursor-pointer bg-purple-600 hover:bg-purple-700 text-white font-bold py-3 px-4 rounded-lg transition duration-150 ease-in-out"
                  >
                    Añadir a la Colección
                  </button>
                </>
              )}
              {type === "deck" && onAdd != undefined && <div>
                <div>
                  <input
                    type="number"
                    min={1}
                    max={isCardDeck(card) ? card.Quantity : 1}
                    value={quantity}
                    onChange={(e) => setQuantity(Number(e.target.value))}
                    className="w-full bg-gray-700 text-white rounded px-3 py-2 border border-gray-600 focus:outline-none focus:ring-2 focus:ring-purple-500"
                  />
                  <button
                    onClick={() => onAction(quantity)}
                    className="w-full cursor-pointer bg-red-600 hover:bg-red-700 text-white font-bold py-3 px-4 rounded-lg transition duration-150 ease-in-out"
                  >
                    Eliminar
                  </button>
                </div>
                <div>
                  <input
                    type="number"
                    min={1}
                    max={isCardDeck(card) ? 3 - card.Quantity : 1}
                    value={quantity}
                    onChange={(e) => setQuantity(Number(e.target.value))}
                    className="w-full bg-gray-700 text-white rounded px-3 py-2 border border-gray-600 focus:outline-none focus:ring-2 focus:ring-purple-500"
                  />
                  <button
                    onClick={() => onAdd(quantity)}
                    className="w-full cursor-pointer bg-green-600 hover:bg-green-700 text-white font-bold py-3 px-4 rounded-lg transition duration-150 ease-in-out"
                  >
                    Añadir
                  </button>
                </div>
              </div>}
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
