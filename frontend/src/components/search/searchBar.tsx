"use client"

import type React from "react"
import { useState, useEffect } from "react"
import type { FilterOptions } from "../../types/search"

// Opciones de ejemplo para los filtros - En una aplicación real, estos vendrían de una API
const filterOptions = {
    tipos: ["Monstruo", "Magia", "Trampa", "Normal", "Efecto", "Fusión", "Sincronía", "Xyz", "Péndulo", "Link", "Ritual"],
    atributos: ["Luz", "Oscuridad", "Tierra", "Agua", "Fuego", "Viento", "Divino"],
    estrellas: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"],
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
            estrellas: "",
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

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
                {/* Filtro de Tipo */}
                <div>
                    <label className="block text-sm font-medium text-gray-400 mb-1">Tipo</label>
                    <select
                        value={localFilters.tipo}
                        onChange={(e) => handleFilterChange("tipo", e.target.value)}
                        className="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                    >
                        <option value="">Todos</option>
                        {filterOptions.tipos.map((tipo) => (
                            <option key={tipo} value={tipo}>
                                {tipo}
                            </option>
                        ))}
                    </select>
                </div>

                {/* Filtro de Arquetipo */}
                <div>
                    <label className="block text-sm font-medium text-gray-400 mb-1">Arquetipo</label>
                    <input
                        type="text"
                        value={searchInput}
                        onChange={handleSearchInputChange}
                        placeholder="Buscar cartas..."
                        className="w-full px-4 py-3 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
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
                        {filterOptions.atributos.map((atributo) => (
                            <option key={atributo} value={atributo}>
                                {atributo}
                            </option>
                        ))}
                    </select>
                </div>

                {/* Filtro de Estrellas */}
                <div>
                    <label className="block text-sm font-medium text-gray-400 mb-1">Estrellas</label>
                    <select
                        value={localFilters.estrellas}
                        onChange={(e) => handleFilterChange("estrellas", e.target.value)}
                        className="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                    >
                        <option value="">Todas</option>
                        {filterOptions.estrellas.map((estrellas) => (
                            <option key={estrellas} value={estrellas}>
                                {estrellas}
                            </option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-400 mb-1">Atk</label>
                    <input
                        type="text"
                        value={searchInput}
                        onChange={handleSearchInputChange}
                        placeholder="Buscar cartas..."
                        className="w-full px-4 py-3 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-400 mb-1">Def</label>
                    <input
                        type="text"
                        value={searchInput}
                        onChange={handleSearchInputChange}
                        placeholder="Buscar cartas..."
                        className="w-full px-4 py-3 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
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
