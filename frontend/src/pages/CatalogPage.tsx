"use client"

import { Header } from "@/components/landing/header"
import { useState, useEffect } from "react"
import SearchBar from "@/components/search/searchBar"
import { ResultsGrid } from "@/components/search/resultsGrid"
import { Sidebar } from "@/components/search/sidebar"
import Pagination from "@/components/search/resultsPagination"
import type { FilterOptions, SearchResult } from "@/types/search"
import { useUser } from '@/contexts/UserContext'
import { toast } from 'sonner'

const API_URL = import.meta.env.VITE_API_URL

export default function CatalogPage() {
  const { user } = useUser()

  const [cards, setCards] = useState<SearchResult[]>([])
  const [searchQuery, setSearchQuery] = useState<string>("")
  const [filters, setFilters] = useState<FilterOptions>({
    tipo: "",
    atributo: "",
    estrellas: "",
    frameType: "",
  })
  // const [currentResults, setCurrentResults] = useState<SearchResult[]>()
  const [currentPage, setCurrentPage] = useState<number>(1)
  const [selectedCard, setSelectedCard] = useState<SearchResult | null>(null)
  const [isSidebarOpen, setIsSidebarOpen] = useState<boolean>(false)
  const [quantity, setQuantity] = useState<number>(1)
  const [isLoading, setIsLoading] = useState<boolean>(false)

  const resultsPerPage = 50

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true)
      try {
        let effectiveApiUrl = `${API_URL}/cards/` // Default to fetching all cards

        const queryParams = new URLSearchParams()
        if (searchQuery.trim() !== "") {
          queryParams.append('name', searchQuery)
        }
        if (filters.tipo.trim() !== "") {
          queryParams.append('type', filters.tipo)
        }
        if (filters.frameType.trim() !== "") {
          queryParams.append('frameType', filters.frameType)
        }
        // Note: filters.atributo and filters.estrellas are not currently sent to the backend.
        // Add them to queryParams here if backend support is added.
        // e.g.:
        // if (filters.atributo.trim() !== "") {
        //   queryParams.append('attribute', filters.atributo);
        // }
        // if (filters.estrellas.trim() !== "") {
        //   queryParams.append('level', filters.estrellas);
        // }

        const queryString = queryParams.toString()

        if (queryString) { // If there are any parameters, use the search endpoint
          effectiveApiUrl = `${API_URL}/cards/search?${queryString}`
        }

        console.log('Fetching from URL:', effectiveApiUrl) // For debugging

        const response = await fetch(effectiveApiUrl, {
          credentials: 'include',
        })

        if (!response.ok) {
          throw new Error(`Failed to fetch cards from ${effectiveApiUrl}`)
        }

        const fetchedData = await response.json()
        console.log('Fetched data from API:', fetchedData)

        setCards(fetchedData["Some cards"])

      } catch (error) {
        console.error('Error fetching cards:', error)
        // Check if error is an instance of Error to safely access message property
        const errorMessage = error instanceof Error ? error.message : 'An unknown error occurred'
        toast.error(`Error fetching cards: ${errorMessage}`)
        setCards([]) // Set to empty array on error
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [searchQuery, filters]) // Dependencies: re-fetch when query or filters change

  // Filtrar resultados basados en la búsqueda y filtros
  const filteredResults = cards.filter((card) => {
    const nameMatches =
      !searchQuery || card.Name.toLowerCase().includes(searchQuery.toLowerCase())

    const tipoMatches = !filters.tipo || card.Type === filters.tipo
    const frameTypeMatches = !filters.frameType || card.FrameType === filters.frameType

    return (
      nameMatches &&
      tipoMatches &&
      frameTypeMatches
    )
  })

  // Calcular resultados para la página actual
  const indexOfLastResult = currentPage * resultsPerPage
  const indexOfFirstResult = indexOfLastResult - resultsPerPage
  const currentResults: SearchResult[] = filteredResults.slice(indexOfFirstResult, indexOfLastResult)
  const totalPages = Math.ceil(filteredResults.length / resultsPerPage)

  // const handleData = () => {

  // }

  const handleSearch = (query: string) => {
    setSearchQuery(query)
    setCurrentPage(1) // Resetear a la primera página cuando se busca
  }

  const handleFilterChange = (newFilters: FilterOptions) => {
    setFilters(newFilters)
    setCurrentPage(1) // Resetear a la primera página cuando se cambian los filtros
  }

  const handleCardClick = (card: SearchResult) => {
    setSelectedCard(card)
    setIsSidebarOpen(true)
  }

  const handlePageChange = (page: number) => {
    setCurrentPage(page)
    // Scroll al inicio de los resultados
    window.scrollTo({
      top: document.getElementById("results-section")?.offsetTop || 0,
      behavior: "smooth",
    })
  }

  const closeSidebar = () => {
    setIsSidebarOpen(false)
    setQuantity(1)
  }

  const handleAddToCollection = async (card: SearchResult, count: number) => {
    if (!user || !user.ID) {
      console.error("User not logged in or user ID is missing.")
      return
    }
    if (!card || !card.ID) {
      console.error("Card data is missing or card ID is missing.")
      return
    }

    try {
      const response = await fetch(`${API_URL}/collections/`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ card_id: card.ID, quantity: count })
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.error || 'Failed to add card to collection')
      }

      console.log("Card added to collection:", card.Name)
      toast.success(`Carta ${card.Name} añadida a la colección`)
      closeSidebar()

    } catch (error) {
      console.error('Error adding card to collection:', error)
      toast.error('Error añadiendo carta a la colección')
    }
  }

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100">
      <Header username={user?.Username || ''} />
      <div className="container mx-auto px-4 py-8">

        <SearchBar onSearch={handleSearch} onFilterChange={handleFilterChange} filters={filters} />

        <div id="results-section" className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Resultados ({isLoading ? '...' : filteredResults.length})</h2>

          {isLoading ? (
            <div className="text-center py-10">
              <p className="text-lg text-gray-400">Cargando cartas...</p>
              {/* You can add a spinner component here if you have one */}
            </div>
          ) : (
            <>
              <ResultsGrid results={currentResults} onCardClick={handleCardClick} />

              {totalPages > 1 && (
                <Pagination currentPage={currentPage} totalPages={totalPages} onPageChange={handlePageChange} />
              )}
            </>
          )}
        </div>
      </div>

      <Sidebar
        card={selectedCard}
        isOpen={isSidebarOpen}
        onClose={closeSidebar}
        onAddToCollection={(card) => handleAddToCollection(card, quantity)}
        quantity={quantity}
        onQuantityChange={setQuantity}
      />
    </div>
  )
}