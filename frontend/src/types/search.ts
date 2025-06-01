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

interface SpellTrapCardDetails {
  CardID: number;
  Type: string;
  Card: BaseCardPlaceholder;
}

interface MonsterCardDetails {
  Atk?: number | null;
  Def?: number | null;
  Level?: number | null;
  Rank?: number | null;
  Attribute?: string | null;
  Race?: string | null;
  Archetype?: string | null;
}

interface LinkMonsterCardDetails {
  Atk?: number | null;
  Attribute?: string | null;
  Race?: string | null;
  Archetype?: string | null;
  LinkValue?: number | null;
  LinkMarkers?: string[] | null;
}

interface PendulumMonsterCardDetails {
  Atk?: number | null;
  Def?: number | null;
  Level?: number | null;
  Attribute?: string | null;
  Race?: string | null;
  Archetype?: string | null;
  PendulumScale?: number | null;
  PendulumEffect?: string | null;
}

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
  level: string
  frameType: string
}
