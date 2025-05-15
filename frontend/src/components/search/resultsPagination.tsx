"use client"

import React from "react"
import { ChevronLeft, ChevronRight } from "lucide-react"

interface PaginationProps {
  currentPage: number
  totalPages: number
  onPageChange: (page: number) => void
}

const Pagination: React.FC<PaginationProps> = ({ currentPage, totalPages, onPageChange }) => {
  // Función para generar el rango de páginas a mostrar
  const getPageRange = () => {
    const range = []
    const maxVisiblePages = 5

    if (totalPages <= maxVisiblePages) {
      // Si hay menos páginas que el máximo visible, mostrar todas
      for (let i = 1; i <= totalPages; i++) {
        range.push(i)
      }
    } else {
      // Siempre mostrar la primera página
      range.push(1)

      let start = Math.max(2, currentPage - 1)
      let end = Math.min(totalPages - 1, currentPage + 1)

      if (currentPage <= 3) {
        end = 4
      } else if (currentPage >= totalPages - 2) {
        start = totalPages - 3
      }

      if (start > 2) {
        range.push("...")
      }

      for (let i = start; i <= end; i++) {
        range.push(i)
      }

      if (end < totalPages - 1) {
        range.push("...")
      }

      range.push(totalPages)
    }

    return range
  }

  return (
    <div className="flex justify-center mt-8 mb-4">
      <div className="flex items-center space-x-1">
        {/* Botón Anterior */}
        <button
          onClick={() => currentPage > 1 && onPageChange(currentPage - 1)}
          disabled={currentPage === 1}
          className={`p-2 rounded-md ${
            currentPage === 1 ? "text-gray-500 cursor-not-allowed" : "text-gray-300 hover:bg-gray-700"
          }`}
        >
          <ChevronLeft size={20} />
        </button>

        {/* Números de página */}
        {getPageRange().map((page, index) => (
          <React.Fragment key={index}>
            {page === "..." ? (
              <span className="px-3 py-2 text-gray-500">...</span>
            ) : (
              <button
                onClick={() => typeof page === "number" && onPageChange(page)}
                className={`px-3 py-1 rounded-md ${
                  currentPage === page ? "bg-purple-600 text-white" : "text-gray-300 hover:bg-gray-700"
                }`}
              >
                {page}
              </button>
            )}
          </React.Fragment>
        ))}

        {/* Botón Siguiente */}
        <button
          onClick={() => currentPage < totalPages && onPageChange(currentPage + 1)}
          disabled={currentPage === totalPages}
          className={`p-2 rounded-md ${
            currentPage === totalPages ? "text-gray-500 cursor-not-allowed" : "text-gray-300 hover:bg-gray-700"
          }`}
        >
          <ChevronRight size={20} />
        </button>
      </div>
    </div>
  )
}

export default Pagination
