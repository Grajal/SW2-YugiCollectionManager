export interface SearchResult {
  id: string
  name: string
  image: string
  tipo: string
  arquetipo?: string
  atributo?: string
  estrellas?: string
  carta: string
  atk?: number
  def?: number
  descripcion: string
}

export interface FilterOptions {
  tipo: string
  atributo: string
  estrellas: string
}
