import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog"
import { toast } from "sonner"

interface UploadYDKDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onFileSelected: (file: File) => void
}

export default function UploadYDKDialog({
  open,
  onOpenChange,
  onFileSelected
}: UploadYDKDialogProps) {
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0]
    if (!selectedFile) return

    if (!selectedFile.name.endsWith(".ydk")) {
      toast("Solo se permiten archivos con extensión .ydk")
      return
    }

    onFileSelected(selectedFile)
    onOpenChange(false) // cerrar el diálogo tras seleccionar
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Sube un archivo .ydk</DialogTitle>
        </DialogHeader>
        <input
          type="file"
          accept=".ydk"
          onChange={handleFileChange}
          className="mt-4"
        />
      </DialogContent>
    </Dialog>
  )
}
