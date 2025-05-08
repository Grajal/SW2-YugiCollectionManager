import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { loginSchema, LoginFormValues } from '@/lib/schemas/authSchemas'
import { toast } from "sonner"

export function LoginForm() {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
  })

  const onLoginSubmit = async (data: LoginFormValues) => {
    try {
      const response = await fetch('http:/localhost:8080/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      })
      const responseData = await response.json()
      if (!response.ok) {
        const errorMessage = responseData.message || 'Error al iniciar sesión. Por favor, inténtalo de nuevo.'
        throw new Error(errorMessage)
      }
      console.log('Inicio de sesión exitoso:', responseData)
      toast.success("¡Bienvenido de nuevo!", {
        description: "Has iniciado sesión correctamente."
      })
    } catch (error) {
      console.error('Error de inicio de sesión:', error)
      toast.error("Error de Inicio de Sesión", {
        description: 'Intentalo de nuevo más tarde'
      })
    }
  }

  return (
    <form onSubmit={handleSubmit(onLoginSubmit)}>
      <div className="grid gap-4 py-4">
        <div className="grid gap-2">
          <Label htmlFor="email-login">
            Correo Electrónico
          </Label>
          <Input id="email-login" type="email" {...register("email")} />
        </div>
        {errors.email && <p className="text-red-500 text-sm">{errors.email.message}</p>}
        <div className="grid gap-2">
          <Label htmlFor="password-login">
            Contraseña
          </Label>
          <Input id="password-login" type="password" {...register("password")} />
        </div>
        {errors.password && <p className="text-red-500 text-sm">{errors.password.message}</p>}
        <Button type="submit" className="w-full mt-2">Iniciar Sesión</Button>
      </div>
    </form>
  )
} 