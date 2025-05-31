import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import type { CollectionItem } from "@/types/collection"

interface ManageCardModalProps {
  isOpen: boolean
  onOpenChange: (isOpen: boolean) => void
  selectedCard: CollectionItem | null
  selectedCardQuantity: number
  setSelectedCardQuantity: (quantity: number) => void
  handleDeleteCard: () => Promise<void>
  handleUpdateCardQuantity: () => Promise<void>
}

export function ManageCardModal({
  isOpen,
  onOpenChange,
  selectedCard,
  selectedCardQuantity,
  setSelectedCardQuantity,
  handleDeleteCard,
  handleUpdateCardQuantity,
}: ManageCardModalProps) {
  if (!selectedCard) return null

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px] bg-gray-800 text-gray-100 border-gray-700">
        <DialogHeader>
          <DialogTitle className="text-xl">Gestionar {selectedCard.Card.Name}</DialogTitle>
          <DialogDescription>
            Actualiza la cantidad o elimina esta carta de tu colección.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
            <label htmlFor="quantity" className="text-right col-span-1">
              Cantidad
            </label>
            <Input
              id="quantity"
              type="number"
              min="1"
              value={selectedCardQuantity}
              onChange={(e) => setSelectedCardQuantity(parseInt(e.target.value, 10))}
              className="col-span-3 bg-gray-700 border-gray-600 focus:ring-purple-500"
            />
          </div>
          <img src={selectedCard.Card.ImageURL} alt={selectedCard.Card.Name} className="w-full h-auto rounded my-2" />
        </div>
        <DialogFooter className="sm:justify-between">
          <Button variant="destructive" onClick={handleDeleteCard}>
            Eliminar Carta
          </Button>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => onOpenChange(false)} className="bg-transparent border-gray-500 text-gray-200 hover:bg-gray-700 hover:text-gray-100">
              Cancelar
            </Button>
            <Button onClick={handleUpdateCardQuantity} className="bg-purple-600 hover:bg-purple-700">
              Guardar Cambios
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
} 