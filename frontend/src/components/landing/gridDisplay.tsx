'use client'
import Window from "./window"

export function Grid() {
    const cards = [
    {
      id: 1,
      title: "",
      description:
        "",
      image: "",
    },
    {
      id: 2,
      title: "",
      description: "",
      image: "",
    },
    {
      id: 3,
      title: "",
      description: "",
      image: "",
    },
  ]

  return (
    <section className="mt-16">
      <div className="mb-12 text-center">
        <img
          src="https://images-wixmp-ed30a86b8c4ca887773594c2.wixmp.com/f/5a6af839-076e-448b-b7e8-47dcfb1f1af3/dez91tq-5368738d-b66f-4964-8a53-8f6916a3c3d2.png/v1/fill/w_869,h_920/dark_magician_render_by_henukim_dez91tq-pre.png?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1cm46YXBwOjdlMGQxODg5ODIyNjQzNzNhNWYwZDQxNWVhMGQyNmUwIiwiaXNzIjoidXJuOmFwcDo3ZTBkMTg4OTgyMjY0MzczYTVmMGQ0MTVlYTBkMjZlMCIsIm9iaiI6W1t7ImhlaWdodCI6Ijw9MjA3MCIsInBhdGgiOiJcL2ZcLzVhNmFmODM5LTA3NmUtNDQ4Yi1iN2U4LTQ3ZGNmYjFmMWFmM1wvZGV6OTF0cS01MzY4NzM4ZC1iNjZmLTQ5NjQtOGE1My04ZjY5MTZhM2MzZDIucG5nIiwid2lkdGgiOiI8PTE5NTQifV1dLCJhdWQiOlsidXJuOnNlcnZpY2U6aW1hZ2Uub3BlcmF0aW9ucyJdfQ.o1rVVMAUDuPpPhwbFqS1YkXjy1d2Lw-RzK4ou2KJIaE"
          alt="Mago Oscuro"
          className="mx-auto"
        />
      </div>

      <div className="grid grid-cols-2 gap-3">
        {/* Ventana 1 */}
        <div className="row-span-2">
          <Window
            title={cards[0].title}
            description={cards[0].description}
            image={cards[0].image}
          />
        </div>

        <div>
          {/* Ventana 2 */}
          <Window
            title={cards[1].title}
            description={cards[1].description}
            image={cards[1].image}
          />
        </div>
        <div>
          {/* Ventana 3 */}
          <Window
            title={cards[2].title}
            description={cards[2].description}
            image={cards[2].image}
          />
        </div>
      </div>
    </section>
  )
}