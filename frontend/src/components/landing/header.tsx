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
          <a href="/decks" className="text-sm/6 font-semibold text-white">
            Decks
          </a>
          <a href="/collection" className="text-sm/6 font-semibold text-white">
            Colección
          </a>

          <a href="/statistics" className="text-sm/6 font-semibold text-white">

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
      <Dialog open={mobileMenuOpen} onClose={setMobileMenuOpen} className="lg:hidden">
        <div className="fixed inset-0 z-10" />
        <DialogPanel className="fixed inset-y-0 right-0 z-10 w-full overflow-y-auto bg-white px-6 py-6 sm:max-w-sm sm:ring-1 sm:ring-gray-200/10">
          <div className="flex items-center justify-between">
            <a href="#" className="-m-1.5 p-1.5">
              <span className="sr-only">Your Company</span>
              <img
                alt=""
                src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600"
                className="h-8 w-auto"
              />
            </a>
            <button
              type="button"
              onClick={() => setMobileMenuOpen(false)}
              className="-m-2.5 rounded-md p-2.5 text-white"
            >
              <span className="sr-only">Close menu</span>
              <XMarkIcon aria-hidden="true" className="size-6" />
            </button>
          </div>
          <div className="mt-6 flow-root">
            <div className="-my-6 divide-y divide-gray-200/10">
              <div className="space-y-2 py-6">
                <Disclosure as="div" className="-mx-3">
                  <DisclosureButton className="group flex w-full items-center justify-between rounded-lg py-2 pr-3.5 pl-3 text-base/7 font-semibold text-white-900 hover:bg-gray-400">
                    Colección
                  </DisclosureButton>
                  <DisclosurePanel className="mt-2 space-y-2">
                  </DisclosurePanel>
                </Disclosure>
                <a
                  href="#"
                  className="-mx-3 block rounded-lg px-3 py-2 text-base/7 font-semibold text-white hover:bg-gray-400"
                >
                  Cartas
                </a>
                <a
                  href="#"
                  className="-mx-3 block rounded-lg px-3 py-2 text-base/7 font-semibold text-white hover:bg-gray-400"
                >
                  Estadísticas
                </a>
                <a
                  href="/statistics"
                  className="-mx-3 block rounded-lg px-3 py-2 text-base/7 font-semibold text-white hover:bg-gray-400"
                >
                  Decks
                </a>
              </div>
              <div className="py-6">
                <a
                  href="#"
                  className="-mx-3 block rounded-lg px-3 py-2.5 text-base/7 font-semibold text-white hover:bg-gray-400"
                >
                  Iniciar sesión
                </a>
              </div>
            </div>
          </div>
        </DialogPanel>
      </Dialog>


    </header>
  )
}