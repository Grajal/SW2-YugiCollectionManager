'use client'
import {
  PopoverGroup,
} from '@headlessui/react'
import { AuthModal } from '../auth/AuthModal'

export function Header({ username }: { username: string }) {


  return (
    <header className="bg-black">
      <nav aria-label="Global" className="mx-auto flex max-w-7xl items-center justify-between p-6 lg:px-8">
        <div className="flex lg:flex-1">
          <a href="#" className="-m-1.5 p-1.5">
            <span className="sr-only">Your Company</span>
            <img
              alt=""
              src="https://cardcosmos.de/cdn/shop/articles/Yu-Gi-Oh.svg?v=1714914619&width=1400"
              className="h-16 w-auto"
            />
          </a>
        </div>
        <PopoverGroup className="hidden lg:flex lg:gap-x-12">
          <a href="/collection" className="text-sm/6 font-semibold text-white">
            Colección
          </a>
          <a href="#" className="text-sm/6 font-semibold text-white">
            Estadísticas
          </a>
          <a href="/cards" className="text-sm/6 font-semibold text-white">
            Cartas
          </a>
        {username === '' ? 
        (<AuthModal></AuthModal>)
        :
        (<div className="hidden lg:flex lg:flex-1 lg:justify-end">
          <span>Bienvenido, {username}</span>
        </div>)}
        </PopoverGroup>
      </nav>
    </header>
  )
}