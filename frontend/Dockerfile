# Etapa 1: Construcción
FROM node:20-alpine AS builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos del proyecto
COPY . .

# Instalar pnpm globalmente
RUN npm install -g pnpm

# Instalar dependencias y construir el proyecto
RUN pnpm install
RUN pnpm run build

# Etapa 2: Servidor
FROM nginx:alpine

# Copiar los archivos construidos al directorio público de NGINX
COPY --from=builder /app/dist /usr/share/nginx/html

# Copiar configuración personalizada de nginx si es necesario
# COPY nginx.conf /etc/nginx/nginx.conf

# Exponer el puerto
EXPOSE 80

# Comando por defecto
CMD ["nginx", "-g", "daemon off;"]
