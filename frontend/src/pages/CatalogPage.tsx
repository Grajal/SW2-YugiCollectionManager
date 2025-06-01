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
import { useDebounce } from 'use-debounce'

const API_URL = import.meta.env.VITE_API_URL
const DEBOUNCE_DELAY = 500

export default function CatalogPage() {
  const { user } = useUser()

  const [cards, setCards] = useState<SearchResult[]>([])
  const [searchQuery, setSearchQuery] = useState<string>("")
  const [filters, setFilters] = useState<FilterOptions>({
    tipo: "",
    atributo: "",
    level: "",
    frameType: "",
  })
  const [currentPage, setCurrentPage] = useState<number>(1)
  const [selectedCard, setSelectedCard] = useState<SearchResult | null>(null)
  const [isSidebarOpen, setIsSidebarOpen] = useState<boolean>(false)
  const [quantity, setQuantity] = useState<number>(1)
  const [isLoading, setIsLoading] = useState<boolean>(false)

  const [debouncedSearchQuery] = useDebounce(searchQuery, DEBOUNCE_DELAY)

  const resultsPerPage = 50

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true)
      try {
        let effectiveApiUrl = `${API_URL}/cards/`

        const queryParams = new URLSearchParams()
        if (debouncedSearchQuery.trim() !== "") {
          queryParams.append('name', debouncedSearchQuery)
        }
        if (filters.tipo.trim() !== "") {
          queryParams.append('type', filters.tipo)
        }
        if (filters.frameType.trim() !== "") {
          queryParams.append('frameType', filters.frameType)
        }

        const queryString = queryParams.toString()

        if (queryString) {
          effectiveApiUrl = `${API_URL}/cards/search?${queryString}`
        }

        console.log('Fetching from URL:', effectiveApiUrl)

        const response = await fetch(effectiveApiUrl, {
          credentials: 'include',
        })

        if (!response.ok) {
          if (response.status === 400) {
            try {
              const errorData = await response.json()
              if (errorData && errorData.error === "Invalid search term") {
                setCards([])
                return
              } else {
                throw new Error(errorData.error || `Error 400: Solicitud incorrecta al buscar cartas.`)
              }
            } catch {
              throw new Error(`Solicitud incorrecta al buscar cartas, respuesta no es JSON válido.`)
            }
          }
          throw new Error(`Error al buscar cartas. Estado: ${response.status}`)
        }

        const fetchedData = await response.json()
        console.log('Fetched data from API:', fetchedData)

        setCards(fetchedData["cards"] || [])

      } catch (error) {
        console.error('Error fetching cards:', error)
        const errorMessage = error instanceof Error ? error.message : 'An unknown error occurred'
        toast.error(`Error fetching cards: ${errorMessage}`)
        setCards([])
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [debouncedSearchQuery, filters])

  const filteredResults = cards.filter((card) => {
    const nameMatches =
      !searchQuery || card.Name.toLowerCase().includes(searchQuery.toLowerCase())

    const tipoMatches = !filters.tipo || card.Type === filters.tipo
    const frameTypeMatches = !filters.frameType || card.FrameType === filters.frameType

    const levelMatches = (() => {
      if (!filters.level) return true
      const targetLevel = parseInt(filters.level, 10)
      if (isNaN(targetLevel)) return true

      const cardLevel = card.MonsterCard?.Level ?? card.PendulumMonsterCard?.Level
      return cardLevel === targetLevel
    })()

    return (
      nameMatches &&
      tipoMatches &&
      frameTypeMatches &&
      levelMatches
    )
  })

  const indexOfLastResult = currentPage * resultsPerPage
  const indexOfFirstResult = indexOfLastResult - resultsPerPage
  const currentResults: SearchResult[] = filteredResults.slice(indexOfFirstResult, indexOfLastResult)
  const totalPages = Math.ceil(filteredResults.length / resultsPerPage)

  const handleSearch = (query: string) => {
    setSearchQuery(query)
    setCurrentPage(1)
    setIsLoading(true)
  }

  const handleFilterChange = (newFilters: FilterOptions) => {
    setFilters(newFilters)
    setCurrentPage(1)
    setIsLoading(true)
  }

  const handleCardClick = (card: SearchResult) => {
    setSelectedCard(card)
    setIsSidebarOpen(true)
  }

  const handlePageChange = (page: number) => {
    setCurrentPage(page)
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