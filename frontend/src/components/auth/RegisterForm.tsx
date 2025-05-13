import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { registerSchema, RegisterFormValues } from '@/lib/schemas/authSchemas'
import { toast } from "sonner"
import { useState } from 'react'

export function RegisterForm() {
  const [formError, setFormError] = useState<string>('')
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
  })

  const onRegisterSubmit = async (data: RegisterFormValues) => {
    try {
      const payloadToSend = {
        email: data.email,
        username: data.username,
        password: data.password,
      }
      console.log('Enviando al backend:', payloadToSend)

      const response = await fetch('http://localhost:8080/api/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payloadToSend),
      })
      const responseData = await response.json()
      if (!response.ok) {
        if (response.status === 409) {
          throw new Error('Ya hay un usuario registrado con esos datos')
        } else {
          throw new Error('Error en el registro. Por favor, inténtalo de nuevo.')
        }
      }

      setFormError('')
      console.log('Registro exitoso:', responseData)
      toast.success("Registro Exitoso", {
        description: "Tu cuenta ha sido creada correctamente."
      })
      // TODO: Manejar registro exitoso (e.g., auto-login, redirigir, close modal)
    } catch (error) {
      console.error('Error en el registro:', error)
      if (error instanceof Error) {
        setFormError(error.message)
      } else {
        setFormError('Hubo un error al registrar tu cuenta. Intentalo de nuevo más tarde')
      }
    }
  }

  return (
    <form onSubmit={handleSubmit(onRegisterSubmit)}>
      <div className="grid gap-4 py-4">
        <div className="grid gap-2">
          <Label htmlFor="email-register">
            Correo Electrónico
          </Label>
          <Input id="email-register" type="email" {...register("email")} onFocus={() => setFormError('')} />
        </div>
        {errors.email && <p className="text-red-500 text-sm">{errors.email.message}</p>}
        <div className="grid gap-2">
          <Label htmlFor="username-register">
            Nombre de usuario
          </Label>
          <Input id="username-register" type="text" {...register("username")} onFocus={() => setFormError('')} />
        </div>
        {errors.username && <p className="text-red-500 text-sm">{errors.username.message}</p>}
        <div className="grid gap-2">
          <Label htmlFor="password-register">
            Contraseña
          </Label>
          <Input id="password-register" type="password" {...register("password")} onFocus={() => setFormError('')} />
        </div>
        {errors.password && <p className="text-red-500 text-sm">{errors.password.message}</p>}
        <div className="grid gap-2">
          <Label htmlFor="confirm-password-register">
            Confirmar Contraseña
          </Label>
          <Input id="confirm-password-register" type="password" {...register("confirmPassword")} onFocus={() => setFormError('')} />
        </div>
        {errors.confirmPassword && <p className="text-red-500 text-sm">{errors.confirmPassword.message}</p>}
        {formError && <p className="text-red-500 text-sm mt-2">{formError}</p>}
        <Button type="submit" className="w-full mt-2">Registrarse</Button>
      </div>
    </form>
  )
} 