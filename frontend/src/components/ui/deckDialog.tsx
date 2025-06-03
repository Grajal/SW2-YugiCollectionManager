import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Deck } from "@/types/deck"

interface SelectDeckDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  decks: Deck[]
  onDeckSelected: (deckId: number) => void
}

export default function SelectDeckDialog({open, onOpenChange, decks, onDeckSelected}: SelectDeckDialogProps) {
  const handleSelect = (deckId: number) => {
    onDeckSelected(deckId)
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Selecciona un deck</DialogTitle>
        </DialogHeader>

        {decks.length === 0 ? (
          <p className="text-gray-500 text-sm mt-2">No hay decks disponibles.</p>
        ) : (
          <div className="flex flex-col gap-2 mt-4 max-h-64 overflow-y-auto pr-1">
            {decks.map((deck) => (
              <Button
                key={deck.ID}
                variant="outline"
                className="justify-start"
                onClick={() => handleSelect(deck.ID)}
              >
                {deck.Name}
              </Button>
            ))}
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}
