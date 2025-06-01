import { Header } from "@/components/landing/header"
import Hero from "@/components/landing/hero"

const API_URL = import.meta.env.VITE_API_URL

export default function LandingPage() {

  const fetchCollection = async() => {
    try{
      const response = await fetch(`${API_URL}/cards/`,{
        method: 'GET',
        credentials: 'include',
      })
      if(!response.ok){
        throw new Error("Error al cargar los datos");
      }
      const data = await response.json();
      console.log(data)
    }catch(error){
      console.error(error)
    }
  }

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100">
      <Header username={''} />
      <div className="container mx-auto px-4 py-8">
        <Hero></Hero>
      </div>
    </div>
  )
}