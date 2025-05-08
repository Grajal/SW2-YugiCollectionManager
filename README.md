# YugiCollectionManager

## ⚒️ RUN BACKEND LOCALLY
### 1. Create the .env file in the root of the project:
<pre>
DB_HOST=localhost
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=dbname
DB_PORT=5432
DB_SSLMODE=disable
JWT_SECRET=your_key
JWT_EXPIRES_IN=2h
</pre>

### 2. Set up the services (API + PostgreSQL):
`docker-compose up --build`
