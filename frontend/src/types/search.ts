// For the nested Card object within SpellTrapCardDetails, based on the provided example
interface BaseCardPlaceholder {
  ID: number;
  CardYGOID: number;
  Name: string;
  Desc: string;
  FrameType: string;
  Type: string;
  ImageURL: string;
  MonsterCard: null;
  SpellTrapCard: null;
  LinkMonsterCard: null;
  PendulumMonsterCard: null;
}

// Details for Spell/Trap cards, based on the provided example
interface SpellTrapCardDetails {
  CardID: number;
  Type: string;
  Card: BaseCardPlaceholder;
}

// Details for Monster cards (non-Link, non-Pendulum specific parts)
// This includes Normal, Effect, Ritual, Fusion, Synchro, Xyz monsters
interface MonsterCardDetails {
  Atk?: number | null;
  Def?: number | null;
  Level?: number | null;    // For non-Xyz monsters (includes Ritual, Fusion, Synchro, Effect, Normal)
  Rank?: number | null;     // For Xyz monsters
  Attribute?: string | null;
  Race?: string | null;     // e.g., Warrior, Dragon, Fiend (monster type category)
  Archetype?: string | null;
}

// Details specific to Link Monsters
interface LinkMonsterCardDetails {
  Atk?: number | null;
  Attribute?: string | null;
  Race?: string | null;
  Archetype?: string | null;
  LinkValue?: number | null;
  LinkMarkers?: string[] | null;
}

// Details specific to Pendulum Monsters
interface PendulumMonsterCardDetails {
  Atk?: number | null;
  Def?: number | null;
  Level?: number | null;
  Attribute?: string | null;
  Race?: string | null;
  Archetype?: string | null;
  PendulumScale?: number | null;
  PendulumEffect?: string | null; // Text of the Pendulum effect
  // Monster effect is in the top-level Desc
}

// The main SearchResult interface, reflecting the structure of the provided example
export interface SearchResult {
  ID: number;
  CardYGOID: number;
  Name: string;
  Desc: string;
  FrameType: string;
  Type: string;
  ImageURL: string;

  MonsterCard?: MonsterCardDetails | null;
  SpellTrapCard?: SpellTrapCardDetails | null;
  LinkMonsterCard?: LinkMonsterCardDetails | null;
  PendulumMonsterCard?: PendulumMonsterCardDetails | null;
}

export interface FilterOptions {
  tipo: string
  atributo: string
  estrellas: string
}
