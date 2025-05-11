# Yugi Collection Manager

Application for the digital management of Yu-Gi-Oh! card collections, developed as a Software Engineering II project.

## Technologies used

- **Frontend**: React + Vite + TypeScript
- **Backend**: Go
- **Database**: PostgreSQL
- **Containers**: Docker y Docker Compose

## 🚀 Launch the whole application locally (frontend + backend + database)

### 1. Requirements

- Docker
- Docker Compose

### 2. Environment variables

Make sure you have the following environment variables defined on your system or in an `.env` file:

<pre>
PGHOST=localhost
PGUSER=your_user
PGPASSWORD=your_password
PGDATABASE=dbname
PGPORT=5432
PGSSLMODE=disable
JWT_SECRET=your_key
JWT_EXPIRES_IN=2h
</pre>

> ⚠️ If you use Railway or Render, you can set these variables in their dashboard. Locally you can set them in an `.env` or in the system environment.

### 3. Launching services

```bash
docker-compose up --build
```
Esto levantará:

* PostgreSQL (yugi_postgres)
* Backend (yugi_backend)
* Frontend (yugi_frontend)

El backend estará disponible en http://localhost:8080 y el frontend en http://localhost:3000