# YugiCollectionManager

## ⚒️ RUN BACKEND LOCALLY
### 1. Create the .env file in the root of the project:
<pre>
PGHOST=localhost
PGUSER=your_user
PGPASSWORD=your_password
PGNAME=dbname
PGPORT=5432
PGSSLMODE=disable
JWT_SECRET=your_key
JWT_EXPIRES_IN=2h
</pre>

### 2. Set up the services (API + PostgreSQL):
`docker-compose up --build`
