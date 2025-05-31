import type { SearchResult } from "./search";

export interface Deck {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  UserID: number;
  Name: string;
  Description: string;
  DeckCards: SearchResult[];
}

export interface CollectionItem {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  UserID: number;
  CardID: number;
  Quantity: number;
  Card: SearchResult;
}

export type Collection = CollectionItem[]; 