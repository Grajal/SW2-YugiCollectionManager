import type React from "react"

interface CardProps {
    title: string
    description: string
    image: string
}
const Window: React.FC<CardProps> = ({ title, description, image }) => {
  return (
    <div
      className={`
        bg-gray-800 rounded-xl overflow-hidden shadow-lg transition-transform hover:scale-[1.02] cursor-pointer
        h-full
      `}
    >
      <div className="relative">
        <img src={image || "/placeholder.svg"} alt={title} className="w-full h-48 object-cover" />
        <div className="absolute inset-0 bg-gradient-to-t from-gray-900 to-transparent opacity-70"></div>
      </div>
      <div className="p-6">
        <h3 className="text-xl font-bold mb-2 text-white">{title}</h3>
        <p className="text-gray-400">{description}</p>
      </div>
    </div>
  )
}

export default Window
