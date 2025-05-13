import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { registerSchema, RegisterFormValues } from '@/lib/schemas/authSchemas'
import { toast } from "sonner"

const API_URL = import.meta.env.VITE_API_URL

export function RegisterForm() {
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

      const response = await fetch(`${API_URL}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payloadToSend),
      })
      const responseData = await response.json()
      if (!response.ok) {
        const errorMessage = responseData.message || 'Error en el registro. Por favor, inténtalo de nuevo.'
        throw new Error(errorMessage)
      }
      console.log('Registro exitoso:', responseData)
      toast.success("Registro Exitoso", {
        description: "Tu cuenta ha sido creada correctamente."
      })
      // TODO: Manejar registro exitoso (e.g., auto-login, redirigir, close modal)
    } catch (error) {
      console.error('Error en el registro:', error)
      toast.error("Error de Registro", {
        description: 'Hubo un error al registrar tu cuenta. Intentalo de nuevo más tarde',
      })
    }
  }

  return (
    <form onSubmit={handleSubmit(onRegisterSubmit)}>
      <div className="grid gap-4 py-4">
        <div className="grid gap-2">
          <Label htmlFor="email-register">
            Correo Electrónico
          </Label>
          <Input id="email-register" type="email" {...register("email")} />
        </div>
        {errors.email && <p className="text-red-500 text-sm">{errors.email.message}</p>}
        <div className="grid gap-2">
          <Label htmlFor="username-register">
            Nombre de usuario
          </Label>
          <Input id="username-register" type="text" {...register("username")} />
        </div>
        {errors.username && <p className="text-red-500 text-sm">{errors.username.message}</p>}
        <div className="grid gap-2">
          <Label htmlFor="password-register">
            Contraseña
          </Label>
          <Input id="password-register" type="password" {...register("password")} />
        </div>
        {errors.password && <p className="text-red-500 text-sm">{errors.password.message}</p>}
        <div className="grid gap-2">
          <Label htmlFor="confirm-password-register">
            Confirmar Contraseña
          </Label>
          <Input id="confirm-password-register" type="password" {...register("confirmPassword")} />
        </div>
        {errors.confirmPassword && <p className="text-red-500 text-sm">{errors.confirmPassword.message}</p>}
        <Button type="submit" className="w-full mt-2">Registrarse</Button>
      </div>
    </form>
  )
}