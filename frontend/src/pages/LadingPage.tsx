import { Header } from "@/components/landing/header"
import { Grid } from "@/components/landing/gridDisplay"
import Hero from "@/components/landing/hero"

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-gray-900 text-gray-100">
      <Header />
      <main className="container mx-auto px-4 py-8">
        <Hero />
        <Grid />
      </main >
    </div>
  )
}
