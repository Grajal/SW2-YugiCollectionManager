"use client"

import { Header } from "@/components/landing/header"
import { useState } from "react"
import SearchBar from "@/components/search/searchBar"
import { ResultsGrid } from "@/components/search/resultsGrid"
import { Sidebar } from "@/components/search/sidebar"
import Pagination from "@/components/search/resultsPagination"
import type { FilterOptions, SearchResult } from "@/types/search"


export default function SearchPage() {
  const [searchQuery, setSearchQuery] = useState<string>("")
  const [archetypeQuery, setArchetypeQuery] = useState<string>("")
  const [atkQuery, setAtkQuery] = useState<string>("")
  const [defQuery, setDefQuery] = useState<string>("")
  const [filters, setFilters] = useState<FilterOptions>({
    tipo: "",
    atributo: "",
    estrellas: "",
  })
  const [currentResults, setCurrentResults] = useState<SearchResult[]>()
  const [currentPage, setCurrentPage] = useState<number>(1)
  const [selectedCard, setSelectedCard] = useState<SearchResult | null>(null)
  const [isSidebarOpen, setIsSidebarOpen] = useState<boolean>(false)

  const resultsPerPage = 50

  const data: SearchResult[] = []

  // Filtrar resultados basados en la búsqueda y filtros
  const filteredResults = data.filter((card) => {
    const nameMatches =
      !searchQuery || card.name.toLowerCase().includes(searchQuery.toLowerCase())

    const archetypeMatches =
      !archetypeQuery || card.arquetipo?.toLowerCase().includes(archetypeQuery.toLowerCase())

    const atkMatches =
      !atkQuery || card.atk?.toString().includes

    const defMatches =
      !defQuery || card.def?.toString() == defQuery

    const tipoMatches = !filters.tipo || card.tipo === filters.tipo
    const atributoMatches = !filters.atributo || card.atributo === filters.atributo
    const estrellasMatches = !filters.estrellas || card.estrellas === filters.estrellas

    return (
      nameMatches &&
      archetypeMatches &&
      tipoMatches &&
      atributoMatches &&
      estrellasMatches &&
      atkMatches &&
      defMatches
    )

  })

  // Calcular resultados para la página actual
  const indexOfLastResult = currentPage * resultsPerPage
  const indexOfFirstResult = indexOfLastResult - resultsPerPage
  //const currentResults : SearchResult[] = filteredResults.slice(indexOfFirstResult, indexOfLastResult)
  const totalPages = Math.ceil(filteredResults.length / resultsPerPage)

  const handleData = () => {

  }

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

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100">
      <Header />
      <div className="container mx-auto px-4 py-8">

        <SearchBar onSearch={handleSearch} onFilterChange={handleFilterChange} filters={filters} />

        <div id="results-section" className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Resultados ({filteredResults.length})</h2>

          <ResultsGrid results={[]} onCardClick={handleCardClick} />

          {totalPages > 1 && (
            <Pagination currentPage={currentPage} totalPages={totalPages} onPageChange={handlePageChange} />
          )}
        </div>
      </div>

      <Sidebar card={selectedCard} isOpen={isSidebarOpen} onClose={closeSidebar} />
    </div>
  )
}

