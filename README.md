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
- AWS S3 Bucket

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
AWS_ACCESS_KEY_ID=your_key_aws
AWS_SECRET_ACCESS_KEY=your_secret_key_aws
AWS_REGION=your_region_aws
AWS_BUCKET_NAME=your_bucket_name
COOKIE_DOMAIN=localhost
</pre>

> ⚠️ If you use Railway or Render, you can set these variables in their dashboard. Locally you can set them in an `.env` or in the system environment. Make sure the bucket is created and has public object access enabled if you want to serve images directly from it.

### 3. Launching services

```bash
docker-compose up --build
```
This will lift:

* PostgreSQL (yugi_postgres)
* Backend (yugi_backend)
* Frontend (yugi_frontend)

The backend will be available at http://localhost:8080 and the frontend at http://localhost:3000.