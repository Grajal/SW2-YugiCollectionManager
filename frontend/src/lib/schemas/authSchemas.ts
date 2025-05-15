import { z } from 'zod';


export const loginSchema = z.object({
  username: z.string().min(1, { message: "Nombre de usuario requerido" }),
  password: z.string().min(1, { message: "Contraseña requerida" }),
});
export type LoginFormValues = z.infer<typeof loginSchema>;


export const registerSchema = z.object({
  email: z.string().email({ message: "Correo electrónico inválido" }).min(1, { message: "Correo electrónico requerido" }),
  username: z.string().min(1, { message: "Nombre de usuario requerido" }),
  password: z.string().min(8, { message: "La contraseña debe tener al menos 8 caracteres" }),
  confirmPassword: z.string().min(1, { message: "Confirmar contraseña es requerido" }),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Las contraseñas no coinciden",
  path: ["confirmPassword"],
});
export type RegisterFormValues = z.infer<typeof registerSchema>; 