import { SignInButton } from "@clerk/clerk-react"
import { WalletCardsIcon as Cards } from "lucide-react"

export default function LandingPage() {
  return (
    <div className="flex min-h-screen flex-col">
      <header className="w-full border-b bg-background/95 backdrop-blur">
        <div className="flex h-16 items-center space-x-4 sm:justify-between sm:space-x-0 px-4">
          <div className="flex gap-2 items-center text-xl font-bold">
            <Cards className="h-6 w-6 text-purple-600" />
            <span>YuGiDeck</span>
          </div>
          <div className="flex flex-1 items-center justify-end space-x-4">
            <nav className="flex items-center space-x-1">
              <SignInButton
                mode="modal"
              />
            </nav>
          </div>
        </div>
      </header>
      <main className="flex-1">
        <section className="w-full py-12 md:py-24 lg:py-32 xl:py-48 bg-gradient-to-b from-background to-purple-950/20">
          <div className="container px-4 md:px-6">
            <div className="grid gap-6 lg:grid-cols-[1fr_400px] lg:gap-12 xl:grid-cols-[1fr_600px]">
              <div className="flex flex-col justify-center space-y-4">
                <div className="space-y-2">
                  <h1 className="text-3xl font-bold tracking-tighter sm:text-5xl xl:text-6xl/none">
                    Gestiona tus mazos de YuGi-Oh como un profesional
                  </h1>
                  <p className="max-w-[600px] text-muted-foreground md:text-xl">
                    Organiza, optimiza y comparte tus estrategias con la aplicación definitiva para duelistas.
                  </p>
                </div>
              </div>
              <div className="mx-auto flex w-full items-center justify-center p-4 sm:p-8">
                <div className="relative aspect-video overflow-hidden rounded-xl border bg-background shadow-xl">
                  <div className="w-full h-full bg-gray-200 dark:bg-gray-800 flex items-center justify-center">
                    <span className="text-gray-500 dark:text-gray-400">App Screenshot</span>
                  </div>
                  <div className="absolute inset-0 bg-gradient-to-t from-background/80 to-background/20" />
                </div>
              </div>
            </div>
          </div>
        </section>
      </main>
    </div>
  )
}
