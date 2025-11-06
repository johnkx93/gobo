# to access docker bash
docker exec -it bo_postgres /bin/sh

# to access db bash
psql -U {{username}} -d {{dbname}}

# check extension installed
psql -U postgres -d appdb -c "SELECT extname, extversion FROM pg_extension;"

# to check current db migration version
migrate -path db/schema -database "postgres://postgres:postgres@localhost:5431/mydb?sslmode=disable" version


# to create new postgres server
docker run -d \
  --name prod_bo \ 
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=mydb \
  -p 5431:5432 \
  postgres:16.10-trixie

  "postgres:16.10-trixie" is the image name

# PRODUCTION - check migration status
PRODUCTION_DATABASE_URL='postgres://postgres:postgres@localhost:5431/mydb?sslmode=disable' make migrate-status-prod
migrate -path db/schema -database "$PRODUCTION_DATABASE_URL" version

# set version
migrate -path db/schema -database "$PRODUCTION_DATABASE_URL" force 1
