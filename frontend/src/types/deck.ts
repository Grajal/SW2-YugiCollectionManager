import type { SearchResult } from "./search";

export interface Deck {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  UserID: number;
  Name: string;
  Description: string;
  DeckCards: CardDeck[];
}

export interface CardDeck {
    Card: SearchResult;
    CardID: number;
    DeckID: number;
    Quantity: number;
    Zone: string;
}