import { useState } from 'react';
import type { Collection, CollectionItem } from "@/types/collection";
import { toast } from 'sonner';

const API_URL = import.meta.env.VITE_API_URL;

interface UseCollectionManagementParams {
  setCollection: React.Dispatch<React.SetStateAction<Collection>>;
  setError: React.Dispatch<React.SetStateAction<string | null>>;
}

export function useCollectionManagement({
  setCollection,
  setError,
}: UseCollectionManagementParams) {
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [selectedCard, setSelectedCard] = useState<CollectionItem | null>(null);
  const [selectedCardQuantity, setSelectedCardQuantity] = useState<number>(1);

  const handleOpenCardModal = (item: CollectionItem) => {
    setSelectedCard(item);
    setSelectedCardQuantity(item.Quantity);
    setIsModalOpen(true);
  };

  const handleDeleteCard = async () => {
    if (!selectedCard) return;
    try {
      const response = await fetch(`${API_URL}/collections/${selectedCard.Card.ID}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ quantity: selectedCard.Quantity }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        const errorMessage = errorData.error || 'Error al eliminar la carta';
        toast.error(errorMessage);
        throw new Error(errorMessage);
      }

      toast.success(`Carta '${selectedCard.Card.Name}' eliminada de la colección`);
      setCollection(prev => prev.filter(item => item.Card.ID !== selectedCard.Card.ID));
      setIsModalOpen(false);
      setSelectedCard(null);
    } catch (err) {
      console.error("Error deleting card:", err);
      if (!(err instanceof Error && err.message.startsWith('Error'))) {
        toast.error('Ocurrió un error desconocido al eliminar la carta');
      }
      setError(err instanceof Error ? err.message : 'Ocurrió un error desconocido al eliminar la carta');
    }
  };

  const handleUpdateCardQuantity = async () => {
    if (!selectedCard) return;
    if (selectedCardQuantity < 0) {
      const errorMsg = "La cantidad no puede ser negativa.";
      toast.error(errorMsg);
      setError(errorMsg);
      return;
    }

    const originalQuantity = selectedCard.Quantity;
    const newQuantity = selectedCardQuantity;
    const cardName = selectedCard.Card.Name;

    try {
      if (newQuantity === originalQuantity) {
        setIsModalOpen(false);
        return;
      }

      let successMessage = "";

      if (newQuantity === 0) {
        const response = await fetch(`${API_URL}/collections/${selectedCard.Card.ID}`, {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ quantity: originalQuantity }),
        });
        if (!response.ok) {
            const errorData = await response.json();
            const errorMessage = errorData.error || 'Error al eliminar la carta por cantidad cero';
            toast.error(errorMessage);
            throw new Error(errorMessage);
        }
        setCollection(prev => prev.filter(item => item.Card.ID !== selectedCard.Card.ID));
        successMessage = `Carta '${cardName}' eliminada de la colección`;
      } else if (newQuantity < originalQuantity) {
        const quantityToRemove = originalQuantity - newQuantity;
        const response = await fetch(`${API_URL}/collections/${selectedCard.Card.ID}`, {
          method: 'DELETE',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({ quantity: quantityToRemove }),
        });
        if (!response.ok) {
          const errorData = await response.json();
          const errorMessage = errorData.error || 'Error al actualizar la cantidad de la carta';
          toast.error(errorMessage);
          throw new Error(errorMessage);
        }
        setCollection(prev =>
          prev.map(item =>
            item.Card.ID === selectedCard.Card.ID ? { ...item, Quantity: newQuantity } : item
          )
        );
        successMessage = `Cantidad de '${cardName}' actualizada a ${newQuantity}`;
      } else {
        const quantityToAdd = newQuantity - originalQuantity;
        const response = await fetch(`${API_URL}/collections/`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({ card_id: selectedCard.Card.ID, quantity: quantityToAdd }),
        });
        if (!response.ok) {
          const errorData = await response.json();
          const errorMessage = errorData.error || 'Error al actualizar la cantidad de la carta';
          toast.error(errorMessage);
          throw new Error(errorMessage);
        }
        setCollection(prev =>
          prev.map(item =>
            item.Card.ID === selectedCard.Card.ID ? { ...item, Quantity: newQuantity } : item
          )
        );
        successMessage = `Cantidad de '${cardName}' actualizada a ${newQuantity}`;
      }
      
      toast.success(successMessage);
      setIsModalOpen(false);
      setSelectedCard(null);
    } catch (err) {
      console.error("Error updating card quantity:", err);
      if (!(err instanceof Error && err.message.startsWith('Error'))) {
        toast.error('Ocurrió un error desconocido al actualizar la cantidad');
      }
      setError(err instanceof Error ? err.message : 'Ocurrió un error desconocido al actualizar la cantidad');
    }
  };

  return {
    isModalOpen,
    selectedCard,
    selectedCardQuantity,
    setIsModalOpen,
    setSelectedCardQuantity,
    handleOpenCardModal,
    handleDeleteCard,
    handleUpdateCardQuantity,
  };
}