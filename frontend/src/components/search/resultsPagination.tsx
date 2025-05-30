"use client"

import React from "react"
import { ChevronLeft, ChevronRight } from "lucide-react"

interface PaginationProps {
  page: String
  onPageChange: (page: string) => void
}

const Pagination: React.FC<PaginationProps> = ({page, onPageChange }) => {

  return (
    <div className="flex justify-center mt-8 mb-4">
      <div className="flex items-center space-x-1">
        {/* Botón Anterior */}
        <button
          onClick={() => page !== "First" && onPageChange("")}
          disabled={page === "First"}
          className={`p-2 rounded-md ${
            page === "First" ? "text-gray-500 cursor-not-allowed" : "text-gray-300 hover:bg-gray-700"
          }`}
        >
          <ChevronLeft size={20} />
        </button>

        {/* Botón Siguiente */}
        <button
          onClick={() => page !== "Last" && onPageChange("")}
          disabled={page === "Last"}
          className={`p-2 rounded-md ${
            page === "Last" ? "text-gray-500 cursor-not-allowed" : "text-gray-300 hover:bg-gray-700"
          }`}
        >
          <ChevronRight size={20} />
        </button>
      </div>
    </div>
  )
}

export default Pagination
