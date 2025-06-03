"use client"

import type React from "react"
import { useState, useEffect } from "react"
import type { FilterOptions } from "../../types/search"

// Updated filter options
const filterOptionsData = {
  tipos: [
    "Effect Monster",
    "Flip Effect Monster",
    "Flip Tuner Effect Monster",
    "Gemini Monster",
    "Normal Monster",
    "Normal Tuner Monster",
    "Pendulum Effect Monster",
    "Pendulum Effect Ritual Monster",
    "Pendulum Flip Effect Monster",
    "Pendulum Normal Monster",
    "Pendulum Tuner Effect Monster",
    "Ritual Effect Monster",
    "Ritual Monster",
    "Spell Card",
    "Spirit Monster",
    "Toon Monster",
    "Trap Card",
    "Tuner Monster",
    "Union Effect Monster",
    "Fusion Monster",
    "Link Monster",
    "Pendulum Effect Fusion Monster",
    "Synchro Monster",
    "Synchro Pendulum Effect Monster",
    "Synchro Tuner Monster",
    "XYZ Monster",
    "XYZ Pendulum Effect Monster",
    "Skill Card",
    "Token",
  ],
  frameTypes: [
    "normal",
    "effect",
    "ritual",
    "fusion",
    "synchro",
    "xyz",
    "link",
    "normal_pendulum",
    "effect_pendulum",
    "ritual_pendulum",
    "fusion_pendulum",
    "synchro_pendulum",
    "xyz_pendulum",
    "spell",
    "trap",
    "token",
    "skill",
  ],
  atributos: ["LIGHT", "DARK", "EARTH", "WATER", "FIRE", "WIND", "DIVINE"],
  level: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"],
}

interface SearchBarProps {
  onSearch: (query: string) => void
  onFilterChange: (filters: FilterOptions) => void
  filters: FilterOptions
}

const SearchBar: React.FC<SearchBarProps> = ({ onSearch, onFilterChange, filters }) => {
  const [searchInput, setSearchInput] = useState<string>("")
  const [localFilters, setLocalFilters] = useState<FilterOptions>(filters)

  // Actualizar filtros
  useEffect(() => {
    setLocalFilters(filters)
  }, [filters])

  const handleSearchInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setSearchInput(value)
    onSearch(value)
  }

  const handleFilterChange = (filterName: keyof FilterOptions, value: string) => {
    const newFilters = {
      ...localFilters,
      [filterName]: value,
    }
    setLocalFilters(newFilters)
    onFilterChange(newFilters)
  }

  const clearFilters = () => {
    const emptyFilters: FilterOptions = {
      tipo: "",
      atributo: "",
      level: "",
      frameType: "",
    }
    setLocalFilters(emptyFilters)
    onFilterChange(emptyFilters)
  }

  return (
    <div className="bg-gray-800 rounded-lg p-6 shadow-lg">
      <div className="mb-6">
        <input
          type="text"
          value={searchInput}
          onChange={handleSearchInputChange}
          placeholder="Buscar cartas..."
          className="w-full px-4 py-3 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {/* Filtro de Tipo */}
        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1">Tipo</label>
          <select
            value={localFilters.tipo}
            onChange={(e) => handleFilterChange("tipo", e.target.value)}
            className="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
          >
            <option value="">Todos</option>
            {filterOptionsData.tipos.map((tipo) => (
              <option key={tipo} value={tipo}>
                {tipo}
              </option>
            ))}
          </select>
        </div>

        {/* Filtro de FrameType */}
        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1">Frame Type</label>
          <select
            value={localFilters.frameType}
            onChange={(e) => handleFilterChange("frameType" as keyof FilterOptions, e.target.value)}
            className="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
          >
            <option value="">Todos</option>
            {filterOptionsData.frameTypes.map((frame) => (
              <option key={frame} value={frame}>
                {frame}
              </option>
            ))}
          </select>
        </div>

        {/* Filtro de Atributo */}
        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1">Atributo</label>
          <select
            value={localFilters.atributo}
            onChange={(e) => handleFilterChange("atributo", e.target.value)}
            className="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
          >
            <option value="">Todos</option>
            {filterOptionsData.atributos.map((atributo) => (
              <option key={atributo} value={atributo}>
                {atributo}
              </option>
            ))}
          </select>
        </div>

        {/* Filtro de Nivel (antes Estrellas) */}
        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1">Nivel</label>
          <select
            value={localFilters.level}
            onChange={(e) => handleFilterChange("level" as keyof FilterOptions, e.target.value)}
            className="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
          >
            <option value="">Todos</option>
            {filterOptionsData.level.map((lvl) => (
              <option key={lvl} value={lvl}>
                {lvl}
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className="mt-4 text-right">
        <button onClick={clearFilters} className="px-4 py-2 text-sm text-gray-300 hover:text-white">
          Limpiar filtros
        </button>
      </div>
    </div>
  )
}

export default SearchBar
