import { useEffect, useState } from "react"
import { Bar, BarChart, ResponsiveContainer, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from "recharts"
import { Header } from "@/components/landing/header"
import { useUser } from '@/contexts/UserContext'

const API_URL = import.meta.env.VITE_API_URL

interface StatsData {
  monster: number
  spell: number
  trap: number
  attributes: { [key: string]: number }
  average_stats: {
    avg_atk: number
    avg_def: number
  }
  total_cards: number
}

export default function StatisticsPage() {
  const { user } = useUser()
  const [stats, setStats] = useState<StatsData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const response = await fetch(`${API_URL}/stats/collection`, {
          credentials: 'include'
        })
        if (!response.ok) {
          throw new Error(`Error HTTP! estado: ${response.status}`)
        }
        const data = await response.json()
        setStats(data)
      } catch (e) {
        if (e instanceof Error) {
          setError(e.message)
        } else {
          setError("Ocurrió un error desconocido")
        }
      } finally {
        setLoading(false)
      }
    }

    fetchStats()
  }, [])

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 text-gray-100 flex flex-col items-center">
        <Header username={user?.Username || ''} />
        <div className="container mx-auto px-4 py-8 text-center">
          <p className="text-xl">Cargando estadísticas...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-900 text-gray-100 flex flex-col items-center">
        <Header username={user?.Username || ''} />
        <div className="container mx-auto px-4 py-8 text-center">
          <p className="text-xl text-red-500">Error al obtener estadísticas: {error}</p>
        </div>
      </div>
    )
  }

  if (!stats) {
    return (
      <div className="min-h-screen bg-gray-900 text-gray-100 flex flex-col items-center">
        <Header username={user?.Username || ''} />
        <div className="container mx-auto px-4 py-8 text-center">
          <p className="text-xl">No hay datos de estadísticas disponibles.</p>
        </div>
      </div>
    )
  }

  const cardTypeData = [
    { name: 'Monstruos', count: stats.monster },
    { name: 'Magias', count: stats.spell },
    { name: 'Trampas', count: stats.trap },
  ]

  const attributeData = Object.entries(stats.attributes).map(([name, count]) => ({
    name,
    count,
  }))

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100">
      <Header username={user?.Username || ''} />
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-8 text-center">Estadísticas de la colección</h1>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div className="bg-gray-800 p-6 rounded-lg shadow-xl">
            <h2 className="text-xl font-semibold mb-2 text-center">Tipos de Cartas</h2>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={cardTypeData}>
                <CartesianGrid vertical={false} />
                <XAxis
                  dataKey="name"
                  tickLine={false}
                  tickMargin={10}
                  axisLine={false}
                />
                <YAxis />
                <Tooltip
                  cursor={false}
                  contentStyle={{ backgroundColor: "hsl(var(--background))", border: "1px solid hsl(var(--border))" }}
                  labelStyle={{ color: "#FFFFFF" }}
                  itemStyle={{ color: "#FFFFFF" }}
                />
                <Legend />
                <Bar dataKey="count" fill="#7C3AED" radius={4} />
              </BarChart>
            </ResponsiveContainer>
          </div>

          <div className="bg-gray-800 p-6 rounded-lg shadow-xl">
            <h2 className="text-xl font-semibold mb-2 text-center">Atributos de Monstruos</h2>
            {attributeData.length > 0 ? (
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={attributeData}>
                  <CartesianGrid vertical={false} />
                  <XAxis
                    dataKey="name"
                    tickLine={false}
                    tickMargin={10}
                    axisLine={false}
                  />
                  <YAxis />
                  <Tooltip
                    cursor={false}
                    contentStyle={{ backgroundColor: "hsl(var(--background))", border: "1px solid hsl(var(--border))" }}
                    labelStyle={{ color: "#FFFFFF" }}
                    itemStyle={{ color: "#FFFFFF" }}
                  />
                  <Legend />
                  <Bar dataKey="count" fill="#7C3AED" radius={4} />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <p className="text-center text-gray-400">No hay datos de atributos de monstruos disponibles.</p>
            )}
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mt-8">
          <div className="bg-gray-800 p-6 rounded-lg shadow-xl">
            <h2 className="text-xl font-semibold mb-2 text-center">Estadísticas Promedio</h2>
            <p className="text-center">ATK Promedio: {stats.average_stats.avg_atk.toFixed(2)}</p>
            <p className="text-center">DEF Promedio: {stats.average_stats.avg_def.toFixed(2)}</p>
          </div>

          <div className="bg-gray-800 p-6 rounded-lg shadow-xl">
            <h2 className="text-xl font-semibold mb-2 text-center">Total de Cartas</h2>
            <p className="text-center text-2xl font-bold">{stats.total_cards}</p>
          </div>
        </div>
      </div>
    </div>
  )
}