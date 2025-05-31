import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { loginSchema, LoginFormValues } from '@/lib/schemas/authSchemas'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

const API_URL = import.meta.env.VITE_API_URL

export function LoginForm() {
  const [formError, setFormError] = useState<string>('')
  const navigate = useNavigate()
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
  })

  const onLoginSubmit = async (data: LoginFormValues) => {
    try {
      const response = await fetch(`${API_URL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(data),
      })
      if (!response.ok) {
        if (response.status === 401) {
          throw new Error('Correo o contraseña incorrectos')
        } else {
          throw new Error('Error al iniciar sesión. Por favor, inténtalo de nuevo.')
        }
      }
      navigate('/collection')
      window.location.reload()
    } catch (error) {
      if (error instanceof Error) {
        setFormError(error.message)
      } else {
        setFormError('Hubo un error al iniciar sesión. Intentalo de nuevo más tarde')
      }
    }
  }

  return (
    <form onSubmit={handleSubmit(onLoginSubmit)}>
      <div className="grid gap-4 py-4">
        <div className="grid gap-2">
          <Label htmlFor="username-login">
            Nombre de usuario
          </Label>
          <Input id="username-login" type="text" {...register("username")} onFocus={() => setFormError('')} />
        </div>
        {errors.username && <p className="text-red-500 text-sm">{errors.username.message}</p>}
        <div className="grid gap-2">
          <Label htmlFor="password-login">
            Contraseña
          </Label>
          <Input id="password-login" type="password" {...register("password")} onFocus={() => setFormError('')} />
        </div>
        {errors.password && <p className="text-red-500 text-sm">{errors.password.message}</p>}
        {formError && <p className="text-red-500 text-sm mt-2">{formError}</p>}
        <Button type="submit" className="w-full mt-2">Iniciar Sesión</Button>
      </div>
    </form>
  )
} 