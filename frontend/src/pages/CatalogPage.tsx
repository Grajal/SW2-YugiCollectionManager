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
  })
  // const [currentResults, setCurrentResults] = useState<SearchResult[]>()
  const [currentPage, setCurrentPage] = useState<number>(1)
  const [selectedCard, setSelectedCard] = useState<SearchResult | null>(null)
  const [isSidebarOpen, setIsSidebarOpen] = useState<boolean>(false)

  const resultsPerPage = 50

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(`${API_URL}/cards/`, {
          credentials: 'include',
        })

        if (!response.ok) {
          throw new Error('Failed to fetch cards')
        }

        const fetchedData = await response.json()
        console.log('Fetched data from API:', fetchedData)

        setCards(fetchedData["Some cards"])

      } catch (error) {
        console.error('Error fetching cards:', error)
        setCards([])
      }
    }
    fetchData()
  }, [])

  // Filtrar resultados basados en la búsqueda y filtros
  const filteredResults = cards.filter((card) => {
    const nameMatches =
      !searchQuery || card.Name.toLowerCase().includes(searchQuery.toLowerCase())

    const tipoMatches = !filters.tipo || card.Type === filters.tipo

    return (
      nameMatches &&
      tipoMatches
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
  }

  const handleAddToCollection = async (card: SearchResult) => {
    if (!user || !user.ID) {
      console.error("User not logged in or user ID is missing.")
      return
    }
    if (!card || !card.ID) {
      console.error("Card data is missing or card ID is missing.")
      return
    }

    try {
      const response = await fetch(`${API_URL}/collections/`, { // Assuming this is the endpoint
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ card_id: card.ID, quantity: 1 })
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
          <h2 className="text-xl font-semibold mb-4">Resultados ({filteredResults.length})</h2>

          <ResultsGrid results={currentResults} onCardClick={handleCardClick} />

          {totalPages > 1 && (
            <Pagination currentPage={currentPage} totalPages={totalPages} onPageChange={handlePageChange} />
          )}
        </div>
      </div>

      <Sidebar card={selectedCard} isOpen={isSidebarOpen} onClose={closeSidebar} onAddToCollection={handleAddToCollection} />
    </div>
  )
}