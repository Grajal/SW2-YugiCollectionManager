import { useState } from 'react';
import type { Collection, CollectionItem } from "@/types/collection";

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
        throw new Error(errorData.error || 'Error al eliminar la carta');
      }

      setCollection(prev => prev.filter(item => item.Card.ID !== selectedCard.Card.ID));
      setIsModalOpen(false);
      setSelectedCard(null);
    } catch (err) {
      console.error("Error deleting card:", err);
      setError(err instanceof Error ? err.message : 'Ocurrió un error desconocido al eliminar la carta');
    }
  };

  const handleUpdateCardQuantity = async () => {
    if (!selectedCard || selectedCardQuantity < 0) {
      setError("La cantidad no puede ser negativa.");
      return;
    }

    const originalQuantity = selectedCard.Quantity;
    const newQuantity = selectedCardQuantity;

    try {
      if (newQuantity === originalQuantity) {
        setIsModalOpen(false);
        return;
      }

      if (newQuantity === 0) {
        const response = await fetch(`${API_URL}/collections/${selectedCard.Card.ID}`, {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ quantity: originalQuantity }),
        });
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Error al eliminar la carta por cantidad cero');
        }
        setCollection(prev => prev.filter(item => item.Card.ID !== selectedCard.Card.ID));
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
          throw new Error(errorData.error || 'Error al actualizar la cantidad de la carta');
        }
        setCollection(prev =>
          prev.map(item =>
            item.Card.ID === selectedCard.Card.ID ? { ...item, Quantity: newQuantity } : item
          )
        );
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
          throw new Error(errorData.error || 'Error al actualizar la cantidad de la carta');
        }
        setCollection(prev =>
          prev.map(item =>
            item.Card.ID === selectedCard.Card.ID ? { ...item, Quantity: newQuantity } : item
          )
        );
      }

      setIsModalOpen(false);
      setSelectedCard(null);
    } catch (err) {
      console.error("Error updating card quantity:", err);
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