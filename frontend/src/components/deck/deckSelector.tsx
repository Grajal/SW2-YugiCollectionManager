"use client"

import type React from "react"
import { Deck } from "@/types/deck"
import { DeckResult } from "./deckElem"
import {  useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { DialogDescription } from "@radix-ui/react-dialog"

interface DecksGridProps {
  results: Deck[]
  onDeckClick: (card: Deck) => void
}
const API_URL = import.meta.env.VITE_API_URL

export const DeckGrid: React.FC<DecksGridProps> = ({ results, onDeckClick }) => {
  const [deckForm, setDeckForm] = useState<boolean>(false)
  const [newDeckName, setNewDeckName] = useState<string>('')
  const [newDeckDesc, setNewDeckDesc] = useState<string>('')
  const postDeck = async() => {
    try{
      const response = await fetch(`${API_URL}/decks/`,{
        method: 'POST',
        headers: {
        "Content-Type": "application/json",
        },
        body: JSON.stringify({
            name: newDeckName,
            description: newDeckDesc
        }),
        credentials: "include"
      });
      console.log(response)
      if(!response.ok){
        throw new Error("Error al cargar los datos");
      }
    }catch(error){
      console.error("ERROR CARGANDO DECKS: " + error)
    }
  }

  const handleDeckForm = () => {
    console.log(!deckForm)
    setDeckForm(!deckForm)
  }
  const handleNewDeck = () => {
    postDeck()
    handleDeckForm()
    window.location.reload()
  }

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
    <div className="container mx-auto px-4 py-8">
        <h1 className="text-sm font-medium text-white">Selecciona el deck</h1>
        <div className="flex flex-wrap gap-4">
            {results.map((deck) => (
                <DeckResult key={deck.ID} deck={deck} onClick={() => onDeckClick(deck)} />
            ))}
            <div className="w-1/12 h-1/12 overflow-hidden shadow-md hover:shadow-lg transition-all hover:scale-105 cursor-pointer" onClick={handleDeckForm}>
                <div className="relative pb-[140%]">
                    {" "}
                    <img src="https://svgsilh.com/svg/1721865.svg" className="absolute inset-0 w-full h-full object-cover"/>
                </div>
                <div className="p-3 text-center">
                    <h3 className="text-sm font-medium text-white truncate" title="Añadir Deck">"Añadir Deck"</h3>
                </div>
            </div>
        </div>
        <Dialog open={deckForm} onOpenChange={setDeckForm}>
            <DialogContent className="sm:max-w-[425px] bg-gray-800 text-gray-100 border-gray-700">
              <DialogHeader>
                <DialogTitle className="text-xl">Nuevo Deck</DialogTitle>
                <DialogDescription>Introduce los datos</DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid grid-cols-4 items-center gap-4">
                  <label htmlFor="name" className="text-right col-span-1">
                    Nombre
                  </label>
                  <Input
                    id="name"
                    type="string"
                    value={newDeckName}
                    onChange={(e) => setNewDeckName(e.target.value)}
                    className="col-span-3 bg-gray-700 border-gray-600 focus:ring-purple-500"
                  />
                  <label htmlFor="name" className="text-right col-span-1">
                    Descripcion
                  </label>
                  <Input
                    id="description"
                    type="string"
                    value={newDeckDesc}
                    onChange={(e) => setNewDeckDesc(e.target.value)}
                    className="col-span-3 bg-gray-700 border-gray-600 focus:ring-purple-500"
                  />
                </div>
                <div className="flex gap-2">
                  <Button variant="outline" onClick={handleDeckForm} className="bg-transparent border-gray-500 text-gray-200 hover:bg-gray-700 hover:text-gray-100">
                    Cancelar
                  </Button>
                  <Button onClick={handleNewDeck} className="bg-purple-600 hover:bg-purple-700">
                    Guardar Cambios
                  </Button>
                </div>
               </div>
            </DialogContent>
    </Dialog>
    </div>
  )
}
